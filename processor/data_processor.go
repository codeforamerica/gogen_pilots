package processor

import (
	"gogen/data"
	"regexp"
	"strings"
)

type DataProcessor struct {
	dojInformation  *data.DOJInformation
	outputDOJWriter DOJWriter
	prop64Matcher   *regexp.Regexp
	stats           dataProcessorStats
	clearanceStats  clearanceStats
	convictionStats convictionStats
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
}

type convictionStats struct {
	numDOJConvictions           map[string]int
	DOJEligibilityByCodeSection map[string]map[string]int
}

type dataProcessorStats struct {
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
	dojInformation *data.DOJInformation,
	outputDOJWriter DOJWriter,
) DataProcessor {
	return DataProcessor{
		dojInformation:  dojInformation,
		outputDOJWriter: outputDOJWriter,
		prop64Matcher:   regexp.MustCompile(`(11357|11358|11359|11360).*`),
		convictionStats: convictionStats{
			numDOJConvictions:           make(map[string]int),
			DOJEligibilityByCodeSection: make(map[string]map[string]int),
		},
		stats: dataProcessorStats{
			matchedSubjectIds: make(map[string]bool),
			matchedCMSIDs:     make(map[string]int),
		},
	}
}
func (d DataProcessor) Process() {
	for i, row := range d.dojInformation.Rows {
		d.outputDOJWriter.WriteDOJEntry(row, d.dojInformation.Eligibilities[i])
	}
	d.outputDOJWriter.Flush()
}

//func (d DataProcessor) ProcessOld() {
//	fmt.Printf("Processing Data from %d Histories... \n", len(d.dojInformation.Histories))
//	var totalTime time.Duration = 0
//
//	fmt.Println("\nDetermining DOJ eligibility")
//	previousSubjectId := ""
//	previousCountOrder := ""
//	d.stats.nDOJSubjects = len(d.dojInformation.Histories)
//	totalRows := float64(len(d.dojInformation.Rows))
//	totalTime = 0
//
//	for i, rawRow := range d.dojInformation.Rows {
//		startTime := time.Now()
//		row := data.NewDOJRow(rawRow)
//		if isProp64Conviction(row, "SAN FRANCISCO") {
//			history := d.dojInformation.Histories[row.SubjectID]
//			if !d.stats.matchedSubjectIds[row.SubjectID] {
//				eligibilityInfo := EligibilityInfoFromDOJRow(&row, history, d.comparisonTime, d.prop64Matcher)
//				d.outputDOJWriter.WriteDOJEntry(rawRow, *eligibilityInfo)
//				if eligibilityInfo.FinalRecommendation == needsReview {
//					d.stats.finalRecNeedsReview++
//				}
//			}
//			if previousCountOrder != row.CountOrder || previousSubjectId != row.SubjectID {
//				d.incrementDOJStats(row)
//			}
//			previousSubjectId = row.SubjectID
//			previousCountOrder = row.CountOrder
//		}
//		totalTime += time.Since(startTime)
//		utilities.PrintProgressBar(float64(i), totalRows, totalTime, "")
//	}
//	d.outputDOJWriter.Flush()
//
//	fmt.Println("\nDetermining summary statistics")
//
//	totalTime = 0
//	numHistories := 0
//	for _, history := range d.dojInformation.Histories {
//		startTime := time.Now()
//
//		d.checkAllConvictionsCleared(history)
//		d.checkConvictionsClearedLast7Years(history)
//		d.checkAllFeloniesReduced(history)
//		d.checkDisqualifiers(history)
//
//		totalTime += time.Since(startTime)
//		numHistories++
//		utilities.PrintProgressBar(float64(numHistories), float64(len(d.dojInformation.Histories)), totalTime, "")
//	}
//
//	fmt.Println("\nComplete...")
//	fmt.Printf("Found %d convictions in DOJ data (%d felonies, %d misdemeanors)\n", d.stats.nDOJProp64Convictions, d.stats.nDOJFelonies, d.stats.nDOJMisdemeanors)
//
//	fmt.Println("==========================================")
//	uniqueDOJHistories := len(d.dojInformation.Histories)
//	fmt.Printf("Total Unique DOJ Histories: %d\n", uniqueDOJHistories)
//	fmt.Printf("Num fully cleared DOJ records out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberFullyClearedRecords, utilities.Percent(d.clearanceStats.numberFullyClearedRecords, uniqueDOJHistories))
//	fmt.Printf("Num cleared DOJ records for last 7 years out of doj records with conv in last 7 years: %d (%d%%)\n", d.clearanceStats.numberClearedRecordsLast7Years, utilities.Percent(d.clearanceStats.numberClearedRecordsLast7Years, d.clearanceStats.numberHistoriesWithConvictionInLast7Years))
//	fmt.Printf("Num DOJ records no felonies out of DOJ recs with felonies: %d (%d%%)\n", d.clearanceStats.numberRecordsNoFelonies, utilities.Percent(d.clearanceStats.numberRecordsNoFelonies, d.clearanceStats.numberHistoriesWithFelonies))
//
//	fmt.Printf("Num DOJ records DQed for Superstrike out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberDQedForSuperstrike, utilities.Percent(d.clearanceStats.numberDQedForSuperstrike, uniqueDOJHistories))
//	fmt.Printf("Num DOJ records DQed for PC290 out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberDQedForPC290, utilities.Percent(d.clearanceStats.numberDQedForPC290, uniqueDOJHistories))
//	fmt.Printf("Num DOJ records DQed for Two Priors out of all DOJ: %d (%d%%)\n", d.clearanceStats.numberDQedForTwoPriors, utilities.Percent(d.clearanceStats.numberDQedForTwoPriors, uniqueDOJHistories))
//	fmt.Printf("Num DOJ convictions by type %v\n", d.convictionStats.numDOJConvictions)
//	fmt.Printf("DOJ Eligibility by code section %v\n", d.convictionStats.DOJEligibilityByCodeSection)
//}

//func isProp64Conviction(row data.DOJRow, county string) bool {
//	if !row.Convicted {
//		return false
//	}
//
//	if row.County != county {
//		return false
//	}
//
//	return strings.HasPrefix(row.CodeSection, "11357") ||
//		strings.HasPrefix(row.CodeSection, "11358") ||
//		strings.HasPrefix(row.CodeSection, "11359") ||
//		strings.HasPrefix(row.CodeSection, "11360")
//}

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

//
//func (d *DataProcessor) incrementConvictionStats(row data.DOJRow, info *data.EligibilityInfo) {
//	if d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection] == nil {
//		d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection] = make(map[string]int)
//		d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection][info.FinalRecommendation]++
//	} else {
//		d.convictionStats.DOJEligibilityByCodeSection[row.CodeSection][info.FinalRecommendation]++
//	}
//}

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

//func (d *DataProcessor) checkConvictionsClearedLast7Years(history *data.DOJHistory) {
//	last7YearsConvictions := make([]*data.DOJRow, 0)
//	for _, conviction := range history.Convictions {
//		if conviction.DispositionDate.After(d.comparisonTime.AddDate(-7, 0, 0)) {
//			last7YearsConvictions = append(last7YearsConvictions, conviction)
//		}
//	}
//
//	if len(last7YearsConvictions) == 0 {
//		return
//	}
//
//	d.clearanceStats.numberHistoriesWithConvictionInLast7Years++
//
//	clearableConvictions := 0
//	for _, conviction := range last7YearsConvictions {
//		if conviction.County != "SAN FRANCISCO" {
//			return
//		}
//		if !d.prop64Matcher.Match([]byte(conviction.CodeSection)) {
//			return
//		}
//		if conviction.Felony && !strings.HasPrefix(conviction.CodeSection, "11357") {
//			return
//		}
//		clearableConvictions++
//	}
//
//	if clearableConvictions > 0 {
//		d.clearanceStats.numberFullyClearedRecords++
//	}
//}

//func (d *DataProcessor) checkAllFeloniesReduced(history *data.DOJHistory) {
//	felonies := make([]*data.DOJRow, 0)
//	for _, conviction := range history.Convictions {
//		if conviction.Felony {
//			felonies = append(felonies, conviction)
//		}
//	}
//
//	if len(felonies) == 0 {
//		return
//	}
//
//	d.clearanceStats.numberHistoriesWithFelonies++
//
//	reducibleFelonies := 0
//	for _, conviction := range felonies {
//		if conviction.County != "SAN FRANCISCO" {
//			return
//		}
//		if !d.prop64Matcher.Match([]byte(conviction.CodeSection)) {
//			return
//		}
//		eligibility := EligibilityInfoFromDOJRow(conviction, history, d.comparisonTime, d.prop64Matcher)
//
//		if !(eligibility.FinalRecommendation == eligible || eligibility.FinalRecommendation == needsReview) {
//			return
//		}
//		reducibleFelonies++
//	}
//
//	if reducibleFelonies > 0 {
//		d.clearanceStats.numberRecordsNoFelonies++
//	}
//}
//
//func (d *DataProcessor) checkDisqualifiers(history *data.DOJHistory) {
//	superstrikeSeen := false
//	pc290Seen := false
//	twoPriorsSeen := false
//	for _, conviction := range history.Convictions {
//		d.convictionStats.numDOJConvictions[conviction.CodeSection]++
//		eligibility := EligibilityInfoFromDOJRow(conviction, history, d.comparisonTime, d.prop64Matcher)
//		if d.convictionStats.DOJEligibilityByCodeSection[eligibility.FinalRecommendation] == nil {
//			d.convictionStats.DOJEligibilityByCodeSection[eligibility.FinalRecommendation] = make(map[string]int)
//		}
//		d.convictionStats.DOJEligibilityByCodeSection[eligibility.FinalRecommendation][conviction.CodeSection]++
//
//		if eligibility.Superstrikes != eligible && !superstrikeSeen {
//			d.clearanceStats.numberDQedForSuperstrike++
//			superstrikeSeen = true
//		}
//		if (eligibility.PC290Registration != eligible || eligibility.PC290Charges != eligible) && !pc290Seen {
//			d.clearanceStats.numberDQedForPC290++
//			pc290Seen = true
//		}
//		if eligibility.TwoPriors != eligible && !twoPriorsSeen {
//			d.clearanceStats.numberDQedForTwoPriors++
//			twoPriorsSeen = true
//		}
//	}
//}
