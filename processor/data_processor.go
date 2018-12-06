package processor

import (
	"encoding/csv"
	"fmt"
	"gogen/data"
	"gogen/utilities"
	"io"
	"time"
)

type DataProcessor struct {
	cmsCSV             *csv.Reader
	weightsInformation *data.WeightsInformation
	dojInformation     *data.DOJInformation
	outputCMSWriter    CMSWriter
	stats              dataProcessorStats
	comparisonTime     time.Time
}

type dataProcessorStats struct {
	nCMSRows         int
	nCMSFelonies     int
	nCMSMisdemeanors int
	unmatchedCMSRows int
}

func NewDataProcessor(
	cmsCSV *csv.Reader,
	weightsInformation *data.WeightsInformation,
	dojInformation *data.DOJInformation,
	outputCMSWriter CMSWriter,
	comparisonTime time.Time,
) DataProcessor {
	return DataProcessor{
		cmsCSV:             cmsCSV,
		weightsInformation: weightsInformation,
		dojInformation:     dojInformation,
		outputCMSWriter:    outputCMSWriter,
		comparisonTime:     comparisonTime,
	}
}

/*
Some Notes:
Using a pure csv.Reader means we don't get line count - how to progress bar?
*/

func (d DataProcessor) Process() {
	d.readHeaders()

	currentRowIndex := 0.0
	totalRows := 9102.0

	fmt.Println("Processing Data...")
	var totalTime time.Duration = 0
	var totalWeightSearchTime time.Duration = 0
	var totalDOJSearchTime time.Duration = 0
	var totalEligibilityTime time.Duration = 0

	for {
		startTime := time.Now()
		rawRow, err := d.cmsCSV.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		row := data.NewCMSEntry(rawRow)

		weightStartTime := time.Now()
		weightsEntry := d.weightsInformation.GetWeight(row.CourtNumber)
		weightEndTime := time.Now()
		totalWeightSearchTime += weightEndTime.Sub(weightStartTime)

		dojStartTime := time.Now()
		dojHistory := d.dojInformation.FindDOJHistory(row)
		dojEndTime := time.Now()
		totalDOJSearchTime += dojEndTime.Sub(dojStartTime)

		eligibilityStartTime := time.Now()
		eligibilityInfo := NewEligibilityInfo(row, weightsEntry, dojHistory, d.comparisonTime)
		eligibilityEndTime := time.Now()
		totalEligibilityTime += eligibilityEndTime.Sub(eligibilityStartTime)

		d.incrementStats(row, dojHistory)
		d.outputCMSWriter.WriteEntry(row, dojHistory, *eligibilityInfo)

		currentRowIndex++
		avgWeightSearchTime := utilities.AverageTime(totalWeightSearchTime, currentRowIndex)
		avgDOJSearchTime := utilities.AverageTime(totalDOJSearchTime, currentRowIndex)
		avgEligibilityTime := utilities.AverageTime(totalEligibilityTime, currentRowIndex)

		tail := fmt.Sprintf("weight: %s, doj: %s, eligibility: %s", avgWeightSearchTime, avgDOJSearchTime, avgEligibilityTime)

		totalTime += time.Since(startTime)
		utilities.PrintProgressBar(currentRowIndex, totalRows, totalTime, tail)
	}
	d.outputCMSWriter.Flush()
	fmt.Println("\nComplete...")
	fmt.Printf("Found %d charges in CMS data (%d felonies, %d misdemeanors)\n", d.stats.nCMSRows, d.stats.nCMSFelonies, d.stats.nCMSMisdemeanors)
	fmt.Printf("Failed to match %d out of %d charges in CMS data (%d%%)\n", d.stats.unmatchedCMSRows, d.stats.nCMSRows, ((d.stats.unmatchedCMSRows)*100)/d.stats.nCMSRows)
	fmt.Printf("Summary Match Data: %+v", d.dojInformation.SummaryMatchData)
}

func (d *DataProcessor) incrementStats(row data.CMSEntry, history *data.DOJHistory) {
	d.stats.nCMSRows++
	if row.Level == "F" {
		d.stats.nCMSFelonies++
	} else {
		d.stats.nCMSMisdemeanors++
	}
	if history == nil {
		d.stats.unmatchedCMSRows++
	}
}

func (d DataProcessor) readHeaders() {
	_, err := d.cmsCSV.Read()
	if err != nil {
		panic(err)
	}
}
