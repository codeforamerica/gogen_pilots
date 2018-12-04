package processor

import (
	"encoding/csv"
	"fmt"
	"gogen/data"
	"io"
)

type DataProcessor struct {
	cmsCSV             *csv.Reader
	weightsInformation *data.WeightsInformation
	dojInformation     *data.DOJInformation
	outputCMSWriter    CMSWriter
	stats              dataProcessorStats
}

type dataProcessorStats struct {
	nCMSRows         int
	nCMSFelonies     int
	nCMSMisdemeanors int
}

func NewDataProcessor(
	cmsCSV *csv.Reader,
	weightsInformation *data.WeightsInformation,
	dojInformation *data.DOJInformation,
	outputCMSWriter CMSWriter,
) DataProcessor {
	return DataProcessor{
		cmsCSV:             cmsCSV,
		weightsInformation: weightsInformation,
		dojInformation:     dojInformation,
		outputCMSWriter:    outputCMSWriter,
	}
}

/*
Some Notes:
Using a pure csv.Reader means we don't get line count - how to progress bar?

*/

func (d DataProcessor) Process() {
	d.readHeaders()

	for {
		rawRow, err := d.cmsCSV.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		row := data.NewCMSEntry(rawRow)
		d.incrementStats(row)

		weightsEntry := d.weightsInformation.GetWeight(row.CourtNumber)
		dojHistory := d.dojInformation.FindDOJHistory(row)
		eligibilityInfo := ComputeEligibility(row, weightsEntry, dojHistory)

		d.outputCMSWriter.WriteEntry(row, *eligibilityInfo)
	}
	d.outputCMSWriter.Flush()
	fmt.Printf("Found %d charges in CMS data (%d felonies, %d misdemeanors)", d.stats.nCMSRows, d.stats.nCMSFelonies, d.stats.nCMSMisdemeanors)
}

func (d *DataProcessor) incrementStats(row data.CMSEntry) {
	d.stats.nCMSRows++
	if row.Level == "F" {
		d.stats.nCMSFelonies++
	} else {
		d.stats.nCMSMisdemeanors++
	}
}

func (d DataProcessor) readHeaders() {
	_, err := d.cmsCSV.Read()
	if err != nil {
		panic(err)
	}
}
