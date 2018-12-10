package processor

import (
	"encoding/csv"
	"fmt"
	"gogen/data"
	"gogen/utilities"
	"io"
	"strings"
	"time"
)

type DataProcessor struct {
	cmsCSV             *csv.Reader
	weightsInformation *data.WeightsInformation
	dojInformation     *data.DOJInformation
	outputCMSWriter    CMSWriter
	outputDOJWriter    CMSWriter
	stats              dataProcessorStats
	comparisonTime     time.Time
}

type dataProcessorStats struct {
	nCMSRows                 int
	nCMSFelonies             int
	nCMSMisdemeanors         int
	unmatchedCMSRows         int
	unmatchedCMSFelonies     int
	unmatchedCMSMisdemeanors int
	nDOJProp64Convictions    int
	nDOJSubjects             int
	nDOJFelonies             int
	nDOJMisdemeanors         int
	unmatchedDOJConvictions  int
	unmatchedDOJFelonies     int
	unmatchedDOJMisdemeanors int
	matchedSubjectIds        map[string]bool
}

func NewDataProcessor(
	cmsCSV *csv.Reader,
	weightsInformation *data.WeightsInformation,
	dojInformation *data.DOJInformation,
	outputCMSWriter CMSWriter,
	outputDOJWriter CMSWriter,
	comparisonTime time.Time,
) DataProcessor {
	return DataProcessor{
		cmsCSV:             cmsCSV,
		weightsInformation: weightsInformation,
		dojInformation:     dojInformation,
		outputCMSWriter:    outputCMSWriter,
		outputDOJWriter:    outputDOJWriter,
		comparisonTime:     comparisonTime,
		stats:              dataProcessorStats{matchedSubjectIds: make(map[string]bool)},
	}
}

func (d DataProcessor) Process() {
	d.readHeaders()

	currentRowIndex := 0.0
	totalRows := 9102.0

	fmt.Printf("Processing Data from %d Histories... \n", len(d.dojInformation.Histories))
	var totalTime time.Duration = 0
	var totalWeightSearchTime time.Duration = 0
	var totalDOJSearchTime time.Duration = 0
	var totalEligibilityTime time.Duration = 0
	var totalMatchingTime time.Duration = 0

	for {
		startTime := time.Now()
		rawRow, err := d.cmsCSV.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		row := data.NewCMSEntry(rawRow)
		if !row.MJCharge() {
			continue
		}

		weightStartTime := time.Now()
		weightsEntry := d.weightsInformation.GetWeight(row.FormattedCourtNumber)
		weightEndTime := time.Now()
		totalWeightSearchTime += weightEndTime.Sub(weightStartTime)

		dojStartTime := time.Now()
		dojHistory, avgMatchTime := d.dojInformation.FindDOJHistory(row)
		dojEndTime := time.Now()
		totalDOJSearchTime += dojEndTime.Sub(dojStartTime)
		totalMatchingTime += avgMatchTime

		eligibilityStartTime := time.Now()
		eligibilityInfo := NewEligibilityInfo(row, weightsEntry, dojHistory, d.comparisonTime)
		eligibilityEndTime := time.Now()
		totalEligibilityTime += eligibilityEndTime.Sub(eligibilityStartTime)

		d.incrementCMSStats(row, dojHistory)
		d.outputCMSWriter.WriteEntry(row, dojHistory, *eligibilityInfo)

		currentRowIndex++
		avgWeightSearchTime := utilities.AverageTime(totalWeightSearchTime, currentRowIndex)
		avgDOJSearchTime := utilities.AverageTime(totalDOJSearchTime, currentRowIndex)
		avgEligibilityTime := utilities.AverageTime(totalEligibilityTime, currentRowIndex)
		avgMatchingTime := utilities.AverageTime(totalMatchingTime, currentRowIndex)

		tail := fmt.Sprintf("weight: %s, doj: %s (match: %s), eligibility: %s", avgWeightSearchTime, avgDOJSearchTime, avgMatchingTime, avgEligibilityTime)

		totalTime += time.Since(startTime)
		utilities.PrintProgressBar(currentRowIndex, totalRows, totalTime, tail)
	}
	d.outputCMSWriter.Flush()

	previousSubjectId := ""
	previousCountOrder := ""
	d.stats.nDOJSubjects = len(d.dojInformation.Histories)
	for _, rawRow := range (d.dojInformation.Rows) {
		row := data.NewDOJRow(rawRow)
		if isProp64Conviction(row, "SAN FRANCISCO") {
			if !d.stats.matchedSubjectIds[row.SubjectID] {
				history := d.dojInformation.Histories[row.SubjectID]
				d.outputDOJWriter.WriteDOJEntry(rawRow, *EligibilityInfoFromDOJRow(&row, history, d.comparisonTime))
			}
			if previousCountOrder != row.CountOrder || previousSubjectId != row.SubjectID {
				d.incrementDOJStats(row)
			}
			previousSubjectId = row.SubjectID
			previousCountOrder = row.CountOrder
		}

	}
	d.outputDOJWriter.Flush()

	fmt.Println("\nComplete...")
	fmt.Printf("Found %d charges in CMS data (%d felonies, %d misdemeanors)\n", d.stats.nCMSRows, d.stats.nCMSFelonies, d.stats.nCMSMisdemeanors)
	fmt.Printf("Found %d charges in DOJ data (%d felonies, %d misdemeanors)\n", d.stats.nDOJProp64Convictions, d.stats.nDOJFelonies, d.stats.nDOJMisdemeanors)

	fmt.Printf("Failed to match %d out of %d charges in CMS data (%d%%)\n", d.stats.unmatchedCMSRows, d.stats.nCMSRows, ((d.stats.unmatchedCMSRows)*100)/d.stats.nCMSRows)
	fmt.Printf("Failed to match %d out of %d felonies in CMS data (%d%%)\n", d.stats.unmatchedCMSFelonies, d.stats.nCMSFelonies, ((d.stats.unmatchedCMSFelonies)*100)/d.stats.nCMSFelonies)
	fmt.Printf("Failed to match %d out of %d misdemeanors in CMS data (%d%%)\n", d.stats.unmatchedCMSMisdemeanors, d.stats.nCMSMisdemeanors, ((d.stats.unmatchedCMSMisdemeanors)*100)/d.stats.nCMSMisdemeanors)
	fmt.Printf("Summary Match Data: %+v", d.dojInformation.SummaryMatchData)
}

func isProp64Conviction(row data.DOJRow, county string) bool {
	if !row.Convicted {
		return false
	}

	if row.County != county {
		return false
	}

	return strings.HasPrefix(row.CodeSection, "11357") ||
		strings.HasPrefix(row.CodeSection, "11358") ||
		strings.HasPrefix(row.CodeSection, "11359") ||
		strings.HasPrefix(row.CodeSection, "11360")
}

func (d *DataProcessor) incrementCMSStats(row data.CMSEntry, history *data.DOJHistory) {
	d.stats.nCMSRows++
	if row.Level == "F" {
		d.stats.nCMSFelonies++
	} else {
		d.stats.nCMSMisdemeanors++
	}
	if history == nil {
		d.stats.unmatchedCMSRows++
		if row.Level == "F" {
			d.stats.unmatchedCMSFelonies++
		} else {
			d.stats.unmatchedCMSMisdemeanors++
		}
	} else {
		d.stats.matchedSubjectIds[history.SubjectID] = true
	}
}

func (d *DataProcessor) incrementDOJStats(row data.DOJRow) {

	d.stats.nDOJProp64Convictions++
	if row.Felony {
		d.stats.nDOJFelonies++
	} else {
		d.stats.nDOJMisdemeanors++
	}

	if !d.stats.matchedSubjectIds[row.SubjectID] {
		d.stats.unmatchedDOJConvictions++
		if row.Felony {
			d.stats.unmatchedDOJFelonies++
		} else {
			d.stats.unmatchedDOJMisdemeanors++
		}
	}
}

func (d DataProcessor) readHeaders() {
	_, err := d.cmsCSV.Read()
	if err != nil {
		panic(err)
	}
}
