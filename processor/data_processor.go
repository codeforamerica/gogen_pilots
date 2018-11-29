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
		raw_row, err := d.cmsCSV.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		row := data.NewCMSEntry(raw_row)

		//weightsEntry, err := d.weightsInformation.GetWeight(row.CourtNumber)
		//if err != nil {
		//	panic(err)
		//}
		//dojHistory := findDOJHistory(row)

		eligibilityInfo := d.computeEligibility(row)

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

func (d DataProcessor) computeEligibility(entry data.CMSEntry) data.EligibilityInfo {
	var eligibleString string
	weight, err := d.weightsInformation.GetWeight(entry.CourtNumber)
	if err != nil {
		eligibleString = "no match"
	}
	eligible, _ := d.weightsInformation.Under1LB(entry.CourtNumber)

	if eligible {
		eligibleString = "eligible"
	} else {
		eligibleString = "ineligible"
	}

	return data.EligibilityInfo{
		QFinalSum: weight,
		Over1Lb:   eligibleString,
	}
}
