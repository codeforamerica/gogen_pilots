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
	clearanceStats     clearanceStats
	convictionStats    convictionStats
	comparisonTime     time.Time
}

type clearanceStats struct {
	numberFullyClearedRecords      int
	numberClearedRecordsLast7Years int
	numberRecordsNoFelonies        int
	numberDQedForSuperstrike       int
	numberDQedForPC290             int
	numberDQedForTwoPriors         int
	numberDQedForOver1LB           int
}

type convictionStats struct {
	numCMSConvictions           map[string]int
	numDOJConvictions           map[string]int
	CMSEligibilityByCodeSection map[string]map[string]int
	DOJEligibilityByCodeSection map[string]map[string]int
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
	finalRecNeedsReview      int
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
		convictionStats: convictionStats{
			numCMSConvictions:           make(map[string]int),
			numDOJConvictions:           make(map[string]int),
			CMSEligibilityByCodeSection: make(map[string]map[string]int),
			DOJEligibilityByCodeSection: make(map[string]map[string]int),
		},
		stats: dataProcessorStats{matchedSubjectIds: make(map[string]bool)},
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

		d.incrementCMSStats(row, dojHistory, eligibilityInfo)
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

	fmt.Println("\nDetermining Unmatched DOJ eligibility")
	previousSubjectId := ""
	previousCountOrder := ""
	d.stats.nDOJSubjects = len(d.dojInformation.Histories)
	totalRows = float64(len(d.dojInformation.Rows))
	totalTime = 0

	for i, rawRow := range d.dojInformation.Rows {
		startTime := time.Now()
		row := data.NewDOJRow(rawRow)
		if isProp64Conviction(row, "SAN FRANCISCO") {
			history := d.dojInformation.Histories[row.SubjectID]
			if !d.stats.matchedSubjectIds[row.SubjectID] {
				d.outputDOJWriter.WriteDOJEntry(rawRow, *EligibilityInfoFromDOJRow(&row, history, d.comparisonTime))
			}
			d.incrementConvictionStats(row, EligibilityInfoFromDOJRow(&row, history, d.comparisonTime))
			if previousCountOrder != row.CountOrder || previousSubjectId != row.SubjectID {
				d.incrementDOJStats(row)
				d.incrementClearanceStats(row, history, EligibilityInfoFromDOJRow(&row, history, d.comparisonTime))
			}
			previousSubjectId = row.SubjectID
			previousCountOrder = row.CountOrder
		}
		totalTime += time.Since(startTime)
		utilities.PrintProgressBar(float64(i), totalRows, totalTime, "")
	}
	d.outputDOJWriter.Flush()

	fmt.Println("\nComplete...")
	fmt.Printf("Found %d convictions in CMS data (%d felonies, %d misdemeanors)\n", d.stats.nCMSRows, d.stats.nCMSFelonies, d.stats.nCMSMisdemeanors)
	fmt.Printf("Found %d convictions in DOJ data (%d felonies, %d misdemeanors)\n", d.stats.nDOJProp64Convictions, d.stats.nDOJFelonies, d.stats.nDOJMisdemeanors)

	fmt.Printf("Failed to match %d out of %d convictions in CMS data (%d%%)\n", d.stats.unmatchedCMSRows, d.stats.nCMSRows, ((d.stats.unmatchedCMSRows)*100)/d.stats.nCMSRows)
	fmt.Printf("Failed to match %d out of %d felonies in CMS data (%d%%)\n", d.stats.unmatchedCMSFelonies, d.stats.nCMSFelonies, ((d.stats.unmatchedCMSFelonies)*100)/d.stats.nCMSFelonies)
	fmt.Printf("Failed to match %d out of %d misdemeanors in CMS data (%d%%)\n", d.stats.unmatchedCMSMisdemeanors, d.stats.nCMSMisdemeanors, ((d.stats.unmatchedCMSMisdemeanors)*100)/d.stats.nCMSMisdemeanors)

	fmt.Printf("Failed to match %d out of %d convictions in DOJ data (%d%%)\n", d.stats.unmatchedDOJConvictions, d.stats.nDOJProp64Convictions, ((d.stats.unmatchedDOJConvictions)*100)/d.stats.nDOJProp64Convictions)
	fmt.Printf("Failed to match %d out of %d unique subjects in DOJ data (%d%%)\n", len(d.dojInformation.Histories)-len(d.stats.matchedSubjectIds), len(d.dojInformation.Histories), ((len(d.dojInformation.Histories)-len(d.stats.matchedSubjectIds))*100)/len(d.dojInformation.Histories))

	fmt.Printf("Summary Match Data: %+v\n", d.dojInformation.SummaryMatchData)

	fmt.Println("==========================================")
	uniqueDOJHistories := len(d.dojInformation.Histories)
	fmt.Printf("Total Unique DOJ Histories: %d\n", uniqueDOJHistories)
	fmt.Printf("Num Convictions Needs Review: %d (%d%%)\n", d.stats.finalRecNeedsReview, utilities.Percent(d.stats.finalRecNeedsReview, d.stats.nCMSRows))
	fmt.Printf("Num fully cleared DOJ records: %d (%d%%)\n", d.clearanceStats.numberFullyClearedRecords, utilities.Percent(d.clearanceStats.numberFullyClearedRecords, uniqueDOJHistories))
	fmt.Printf("Num cleared DOJ records for last 7 years: %d (%d%%)\n", d.clearanceStats.numberClearedRecordsLast7Years, utilities.Percent(d.clearanceStats.numberClearedRecordsLast7Years, uniqueDOJHistories))
	fmt.Printf("Num DOJ records no felonies: %d (%d%%)\n", d.clearanceStats.numberRecordsNoFelonies, utilities.Percent(d.clearanceStats.numberRecordsNoFelonies, uniqueDOJHistories))
	fmt.Printf("Num DOJ records DQed for Superstrike: %d (%d%%)\n", d.clearanceStats.numberDQedForSuperstrike, utilities.Percent(d.clearanceStats.numberDQedForSuperstrike, uniqueDOJHistories))
	fmt.Printf("Num DOJ records DQed for PC290: %d (%d%%)\n", d.clearanceStats.numberDQedForPC290, utilities.Percent(d.clearanceStats.numberDQedForPC290, uniqueDOJHistories))
	fmt.Printf("Num DOJ records DQed for Two Priors: %d (%d%%)\n", d.clearanceStats.numberDQedForTwoPriors, utilities.Percent(d.clearanceStats.numberDQedForTwoPriors, uniqueDOJHistories))
	fmt.Printf("Num CMS records DQed for Over 1lb: %d (%d%%)\n", d.clearanceStats.numberDQedForOver1LB, utilities.Percent(d.clearanceStats.numberDQedForOver1LB, d.stats.nCMSRows))
	fmt.Printf("Num CMS convictions by type %v\n", d.convictionStats.numCMSConvictions)
	fmt.Printf("Num DOJ convictions by type %v\n", d.convictionStats.numDOJConvictions)
	fmt.Printf("CMS Eligibility by code section %v\n", d.convictionStats.CMSEligibilityByCodeSection)
	fmt.Printf("DOJ Eligibility by code section %v\n", d.convictionStats.DOJEligibilityByCodeSection)
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

func (d *DataProcessor) incrementCMSStats(row data.CMSEntry, history *data.DOJHistory, info *EligibilityInfo) {
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
	if info.FinalRecommendation == needsReview {
		d.stats.finalRecNeedsReview++
	}
	if info.Over1Lb == ineligible {
		d.clearanceStats.numberDQedForOver1LB++
	}
	d.convictionStats.numCMSConvictions[row.Charge]++
	if d.convictionStats.CMSEligibilityByCodeSection[row.Charge] == nil {
		d.convictionStats.CMSEligibilityByCodeSection[row.Charge] = make(map[string]int)
		d.convictionStats.CMSEligibilityByCodeSection[row.Charge][info.FinalRecommendation]++
	} else {
		d.convictionStats.CMSEligibilityByCodeSection[row.Charge][info.FinalRecommendation]++
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

func (d *DataProcessor) incrementClearanceStats(row data.DOJRow, history *data.DOJHistory, eligibilityInfo *EligibilityInfo) {
	if history.OnlyProp64MisdemeanorsSince(time.Time{}) && eligibilityInfo.FinalRecommendation == eligible {
		d.clearanceStats.numberFullyClearedRecords++
	}

	if history.OnlyProp64MisdemeanorsSince(d.comparisonTime.AddDate(-7, 0, 0)) && eligibilityInfo.FinalRecommendation == eligible {
		d.clearanceStats.numberClearedRecordsLast7Years++
	}

	if history.OnlyProp64FeloniesSince(time.Time{}) && eligibilityInfo.FinalRecommendation == eligible {
		d.clearanceStats.numberRecordsNoFelonies++
	}

	if eligibilityInfo.Superstrikes == ineligible {
		d.clearanceStats.numberDQedForSuperstrike++
	}

	pc290Charges := eligibilityInfo.PC290Charges == ineligible
	pc290CodeSections := eligibilityInfo.PC290CodeSections == ineligible
	pc290Registration := eligibilityInfo.PC290Registration == ineligible
	if pc290Charges || pc290CodeSections || pc290Registration {
		d.clearanceStats.numberDQedForPC290++
	}

	if eligibilityInfo.TwoPriors == ineligible {
		d.clearanceStats.numberDQedForTwoPriors++
	}

	for _, conviction := range history.Convictions {
		d.convictionStats.numDOJConvictions[conviction.CodeSection]++
	}
}

func (d *DataProcessor) incrementConvictionStats(row data.DOJRow, info *EligibilityInfo) {
	if d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection] == nil {
		d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection] = make(map[string]int)
		d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection][info.FinalRecommendation]++
	} else {
		d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection][info.FinalRecommendation]++
	}
}

func (d DataProcessor) readHeaders() {
	_, err := d.cmsCSV.Read()
	if err != nil {
		panic(err)
	}
}
