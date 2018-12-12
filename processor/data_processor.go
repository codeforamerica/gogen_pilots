package processor

import (
	"encoding/csv"
	"fmt"
	"gogen/data"
	"gogen/utilities"
	"io"
	"regexp"
	"strings"
	"time"
)

type DataProcessor struct {
	cmsCSV             *csv.Reader
	weightsInformation *data.WeightsInformation
	dojInformation     *data.DOJInformation
	outputCMSWriter    CMSWriter
	outputDOJWriter    CMSWriter
	prop64Matcher      *regexp.Regexp
	stats              dataProcessorStats
	clearanceStats     clearanceStats
	convictionStats    convictionStats
	comparisonTime     time.Time
}

type clearanceStats struct {
	numberFullyClearedRecords                 int
	numberClearedRecordsLast7Years            int
	numberHistoriesWithConvictionInLast7Years int
	numberRecordsNoFelonies                   int
	numberHistoriesWithFelonies               int
	numberDQedForSuperstrike                  int
	numberDQedForPC290                        int
	numberDQedForTwoPriors                    int
	numberDQedForOver1LB                      map[string]bool
	numberOnlyDQIsWeight                      int
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
	matchedCMSIDs            map[string]int
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
		prop64Matcher:      regexp.MustCompile(`(11357|11358|11359|11360).*`),
		convictionStats: convictionStats{
			numCMSConvictions:           make(map[string]int),
			numDOJConvictions:           make(map[string]int),
			CMSEligibilityByCodeSection: make(map[string]map[string]int),
			DOJEligibilityByCodeSection: make(map[string]map[string]int),
		},
		stats: dataProcessorStats{
			matchedSubjectIds: make(map[string]bool),
			matchedCMSIDs:     make(map[string]int),
		},
		clearanceStats: clearanceStats{
			numberDQedForOver1LB: make(map[string]bool),
		},
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
		d.stats.matchedCMSIDs[row.Name+row.DateOfBirth.Format("01/02/2006")]++
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
		eligibilityInfo := NewEligibilityInfo(row, weightsEntry, dojHistory, d.comparisonTime, d.prop64Matcher)
		d.checkDisqualifiedFor1LB(row, eligibilityInfo)
		if eligibilityInfo.FinalRecommendation == needsReview {
			d.stats.finalRecNeedsReview++
		}
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
				eligibilityInfo := EligibilityInfoFromDOJRow(&row, history, d.comparisonTime, d.prop64Matcher)
				d.outputDOJWriter.WriteDOJEntry(rawRow, *eligibilityInfo)
				if eligibilityInfo.FinalRecommendation == needsReview {
					d.stats.finalRecNeedsReview++
				}
			}
			if previousCountOrder != row.CountOrder || previousSubjectId != row.SubjectID {
				d.incrementDOJStats(row)
			}
			previousSubjectId = row.SubjectID
			previousCountOrder = row.CountOrder
		}
		totalTime += time.Since(startTime)
		utilities.PrintProgressBar(float64(i), totalRows, totalTime, "")
	}
	d.outputDOJWriter.Flush()

	fmt.Println("\nDetermining summary statistics")

	totalTime = 0
	numHistories := 0
	for _, history := range d.dojInformation.Histories {
		startTime := time.Now()

		d.checkAllConvictionsCleared(history)
		d.checkConvictionsClearedLast7Years(history)
		d.checkAllFeloniesReduced(history)
		d.checkDisqualifiers(history)

		totalTime += time.Since(startTime)
		numHistories++
		utilities.PrintProgressBar(float64(numHistories), float64(len(d.dojInformation.Histories)), totalTime, "")
	}

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
	uniqueCMSHistories := len(d.stats.matchedCMSIDs)
	fmt.Printf("Total Unique DOJ Histories: %d\n", uniqueDOJHistories)
	fmt.Printf("Total Unique CMS Histories: %d\n", uniqueCMSHistories)
	fmt.Printf("Num fully cleared DOJ records out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberFullyClearedRecords, utilities.Percent(d.clearanceStats.numberFullyClearedRecords, uniqueDOJHistories))
	fmt.Printf("Num cleared DOJ records for last 7 years out of doj records with conv in last 7 years: %d (%d%%)\n", d.clearanceStats.numberClearedRecordsLast7Years, utilities.Percent(d.clearanceStats.numberClearedRecordsLast7Years, d.clearanceStats.numberHistoriesWithConvictionInLast7Years))
	fmt.Printf("Num DOJ records no felonies out of DOJ recs with felonies: %d (%d%%)\n", d.clearanceStats.numberRecordsNoFelonies, utilities.Percent(d.clearanceStats.numberRecordsNoFelonies, d.clearanceStats.numberHistoriesWithFelonies))

	fmt.Printf("Num Convictions Needs Review: %d (%d%%)\n", d.stats.finalRecNeedsReview, utilities.Percent(d.stats.finalRecNeedsReview, d.stats.nCMSRows))
	fmt.Printf("Num DOJ records DQed for Superstrike out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberDQedForSuperstrike, utilities.Percent(d.clearanceStats.numberDQedForSuperstrike, uniqueDOJHistories))
	fmt.Printf("Num DOJ records DQed for PC290 out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberDQedForPC290, utilities.Percent(d.clearanceStats.numberDQedForPC290, uniqueDOJHistories))
	fmt.Printf("Num DOJ records DQed for Two Priors out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberDQedForTwoPriors, utilities.Percent(d.clearanceStats.numberDQedForTwoPriors, uniqueDOJHistories))
	fmt.Printf("Num CMS records DQed for Over 1lb out of all CMS: %d (%d%%)\n", len(d.clearanceStats.numberDQedForOver1LB), utilities.Percent(len(d.clearanceStats.numberDQedForOver1LB), d.stats.nCMSRows))
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
	d.convictionStats.numCMSConvictions[row.Charge]++
	if d.convictionStats.CMSEligibilityByCodeSection[info.FinalRecommendation] == nil {
		d.convictionStats.CMSEligibilityByCodeSection[info.FinalRecommendation] = make(map[string]int)
	}
	d.convictionStats.CMSEligibilityByCodeSection[info.FinalRecommendation][row.Charge]++
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

func (d *DataProcessor) checkAllConvictionsCleared(history *data.DOJHistory) {
	clearableConvictions := 0
	for _, conviction := range history.Convictions {
		if conviction.County != "SAN FRANCISCO" {
			return
		}
		if !d.prop64Matcher.Match([]byte(conviction.CodeSection)) {
			return
		}
		if conviction.Felony && !strings.HasPrefix(conviction.CodeSection, "11357") {
			return
		}
		clearableConvictions++
	}
	if clearableConvictions > 0 {
		d.clearanceStats.numberFullyClearedRecords++
	}
}

func (d *DataProcessor) checkConvictionsClearedLast7Years(history *data.DOJHistory) {
	last7YearsConvictions := make([]*data.DOJRow, 0)
	for _, conviction := range history.Convictions {
		if conviction.DispositionDate.After(d.comparisonTime.AddDate(-7, 0, 0)) {
			last7YearsConvictions = append(last7YearsConvictions, conviction)
		}
	}

	if len(last7YearsConvictions) == 0 {
		return
	}

	d.clearanceStats.numberHistoriesWithConvictionInLast7Years++

	clearableConvictions := 0
	for _, conviction := range last7YearsConvictions {
		if conviction.County != "SAN FRANCISCO" {
			return
		}
		if !d.prop64Matcher.Match([]byte(conviction.CodeSection)) {
			return
		}
		if conviction.Felony && !strings.HasPrefix(conviction.CodeSection, "11357") {
			return
		}
		clearableConvictions++
	}

	if clearableConvictions > 0 {
		d.clearanceStats.numberFullyClearedRecords++
	}
}

func (d *DataProcessor) checkAllFeloniesReduced(history *data.DOJHistory) {
	felonies := make([]*data.DOJRow, 0)
	for _, conviction := range history.Convictions {
		if conviction.Felony {
			felonies = append(felonies, conviction)
		}
	}

	if len(felonies) == 0 {
		return
	}

	d.clearanceStats.numberHistoriesWithFelonies++

	reducibleFelonies := 0
	for _, conviction := range felonies {
		if conviction.County != "SAN FRANCISCO" {
			return
		}
		if !d.prop64Matcher.Match([]byte(conviction.CodeSection)) {
			return
		}
		eligibility := EligibilityInfoFromDOJRow(conviction, history, d.comparisonTime, d.prop64Matcher)

		if !(eligibility.FinalRecommendation == eligible || eligibility.FinalRecommendation == needsReview) {
			return
		}
		reducibleFelonies++
	}

	if reducibleFelonies > 0 {
		d.clearanceStats.numberRecordsNoFelonies++
	}
}

func (d *DataProcessor) checkDisqualifiers(history *data.DOJHistory) {
	superstrikeSeen := false
	pc290Seen := false
	twoPriorsSeen := false
	for _, conviction := range history.Convictions {
		d.convictionStats.numDOJConvictions[conviction.CodeSection]++
		eligibility := EligibilityInfoFromDOJRow(conviction, history, d.comparisonTime, d.prop64Matcher)
		if d.convictionStats.DOJEligibilityByCodeSection[eligibility.FinalRecommendation] == nil {
			d.convictionStats.DOJEligibilityByCodeSection[eligibility.FinalRecommendation] = make(map[string]int)
		}
		d.convictionStats.DOJEligibilityByCodeSection[eligibility.FinalRecommendation][conviction.CodeSection]++

		if eligibility.Superstrikes != eligible && !superstrikeSeen {
			d.clearanceStats.numberDQedForSuperstrike++
			superstrikeSeen = true
		}
		if (eligibility.PC290Registration != eligible || eligibility.PC290Charges != eligible) && !pc290Seen {
			d.clearanceStats.numberDQedForPC290++
			pc290Seen = true
		}
		if eligibility.TwoPriors != eligible && !twoPriorsSeen {
			d.clearanceStats.numberDQedForTwoPriors++
			twoPriorsSeen = true
		}
	}
}

func (d *DataProcessor) checkDisqualifiedFor1LB(entry data.CMSEntry, info *EligibilityInfo) {
	if info.Over1Lb == ineligible {
		d.clearanceStats.numberDQedForOver1LB[entry.Name+entry.DateOfBirth.Format("01/02/2006")] = true
	}
}
