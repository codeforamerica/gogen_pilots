package processor

import (
	"encoding/csv"
	"gogen/data"
	"io"
)

type DataProcessor struct {
	cmsCSV             *csv.Reader
	weightsInformation *data.WeightsInformation
	dojInformation     *data.DOJInformation
	outputCMSWriter    data.CMSWriter
}

func NewDataProcessor(
	cmsCSV *csv.Reader,
	weightsInformation *data.WeightsInformation,
	dojInformation *data.DOJInformation,
	outputCMSWriter data.CMSWriter,
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

		weightsEntry := d.weightsInformation.GetWeight(row.CourtNumber)
		//dojHistory := d.dojInformation.findDOJHistory(row)

		eligibilityInfo := ComputeEligibility(row, weightsEntry)

		d.outputCMSWriter.WriteEntry(row, eligibilityInfo)
	}
	d.outputCMSWriter.Flush()
}

func (d DataProcessor) readHeaders() {
	_, err := d.cmsCSV.Read()
	if err != nil {
		panic(err)
	}
}
