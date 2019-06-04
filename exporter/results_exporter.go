package exporter

import (
	"fmt"
	"gogen/data"
	"regexp"
	"sort"
)

type Exporter struct {
	dojInformation                      *data.DOJInformation
	outputDOJWriter                     DOJWriter
	outputCondensedDOJWriter            DOJWriter
	outputProp64ConvictionsDOJWriter    DOJWriter
	prop64Matcher                       *regexp.Regexp
	stats                               dataProcessorStats
	totalClearanceResults               totalClearanceResults
	clearanceByCodeSection              clearanceByCodeSection
	relatedChargeClearanceByCodeSection clearanceByCodeSection
	totalExistingConvictions            totalExistingConvictions
	convictionsByCodeSection            convictionsByCodeSection
	relatedConvictionsByCodeSection     convictionsByCodeSection
}

type totalClearanceResults struct {
	numberHistoriesWithFelonies                               int
	numberHistoriesWithConvictionInLast7Years                 int
	numberNoLongerHaveFelony                                  int
	numberNoLongerHaveFelonyIfAllSealed                       int
	numberNoLongerHaveFelonyIfAllSealedIncludingRelated       int
	numberClearedRecordsLast7Years                            int
	numberClearedRecordsLast7YearsIfAllSealed                 int
	numberClearedRecordsLast7YearsIfAllSealedIncludingRelated int
	numberNoMoreConvictions                                   int
	numberNoMoreConvictionsIfAllSealed                        int
	numberNoMoreConvictionsIfAllSealedIncludingRelated        int
}

type clearanceByCodeSection struct {
	numberEligibilityByReason                     map[string]int
	numberDismissedByCodeSection                  map[string]int
	numberReducedByCodeSection                    map[string]int
	numberIneligibleByCodeSection                 map[string]int
	numberMaybeEligibleFlagForReviewByCodeSection map[string]int
}

type totalExistingConvictions struct {
	totalCountyConvictions       int
	totalHasFelony               int
	totalHasConvictionLast7Years int
	totalHasConvictions          int
}

type convictionsByCodeSection struct {
	totalConvictionsByCodeSection  map[string]int
	countyConvictionsByCodeSection map[string]int
	DOJEligibilityByCodeSection    map[string]map[string]int
}

type historySummaryStats struct {
	feloniesDismissed               int
	feloniesReduced                 int
	maybeEligibleFelonies           int
	notEligibleFelonies             int
	misdemeanorsDismissed           int
	feloniesDismissedLast7Years     int
	feloniesReducedLast7Years       int
	maybeEligibleFeloniesLast7Years int
	notEligibleFeloniesLast7Years   int
	misdemeanorsDismissedLast7Years int
}

type dataProcessorStats struct {
	nDOJProp64Convictions int
	nDOJSubjects          int
	nDOJFelonies          int
	nDOJMisdemeanors      int
}

func NewExporter(
	dojInformation *data.DOJInformation,
	outputDOJWriter DOJWriter,
	outputCondensedDOJWriter DOJWriter,
	outputProp64ConvictionsDOJWriter DOJWriter,
) Exporter {
	return Exporter{
		dojInformation:           				dojInformation,
		outputDOJWriter:          				outputDOJWriter,
		outputCondensedDOJWriter: 				outputCondensedDOJWriter,
		outputProp64ConvictionsDOJWriter: 		outputProp64ConvictionsDOJWriter,
		totalClearanceResults:    				totalClearanceResults{},
		clearanceByCodeSection: 				clearanceByCodeSection{
			numberEligibilityByReason:                     make(map[string]int),
			numberDismissedByCodeSection:                  make(map[string]int),
			numberReducedByCodeSection:                    make(map[string]int),
			numberIneligibleByCodeSection:                 make(map[string]int),
			numberMaybeEligibleFlagForReviewByCodeSection: make(map[string]int),
		},
		relatedChargeClearanceByCodeSection: 	clearanceByCodeSection{
			numberEligibilityByReason:                     make(map[string]int),
			numberDismissedByCodeSection:                  make(map[string]int),
			numberReducedByCodeSection:                    make(map[string]int),
			numberIneligibleByCodeSection:                 make(map[string]int),
			numberMaybeEligibleFlagForReviewByCodeSection: make(map[string]int),
		},
		totalExistingConvictions: 				totalExistingConvictions{},
		convictionsByCodeSection: 				convictionsByCodeSection{
			totalConvictionsByCodeSection:  make(map[string]int),
			countyConvictionsByCodeSection: make(map[string]int),
		},
		relatedConvictionsByCodeSection: 		convictionsByCodeSection{
			totalConvictionsByCodeSection:  make(map[string]int),
			countyConvictionsByCodeSection: make(map[string]int),
		},
	}
}

func (e *Exporter) incrementConvictionsByCodeSection(
	conviction *data.DOJRow,
	county string,
	matchedCodeSection string,
	convictionsByCodeSection *convictionsByCodeSection,
) {
	if matchedCodeSection != "" {
		convictionsByCodeSection.totalConvictionsByCodeSection[matchedCodeSection]++
		if conviction.County == county {
			convictionsByCodeSection.countyConvictionsByCodeSection[matchedCodeSection]++
		}
	}

}

func (e *Exporter) incrementEligibilities(
	eligibility *data.EligibilityInfo,
	conviction *data.DOJRow,
	county string,
	matchedCodeSection string,
	historyStats *historySummaryStats,
	clearanceByCodeSection *clearanceByCodeSection,
) {
	switch eligibility.EligibilityDetermination {
	case "Eligible for Dismissal":
		clearanceByCodeSection.numberDismissedByCodeSection[matchedCodeSection]++
		if conviction.IsFelony {
			historyStats.feloniesDismissed++
			if conviction.OccurredInLast7Years() {
				historyStats.feloniesDismissedLast7Years++
			}
		} else {
			historyStats.misdemeanorsDismissed++
			if conviction.OccurredInLast7Years() {
				historyStats.misdemeanorsDismissedLast7Years++
			}
		}

	case "Eligible for Reduction":
		clearanceByCodeSection.numberReducedByCodeSection[matchedCodeSection]++
		historyStats.feloniesReduced++
		if conviction.OccurredInLast7Years() {
			historyStats.feloniesReducedLast7Years++
		}

	case "Not eligible":
		clearanceByCodeSection.numberIneligibleByCodeSection[matchedCodeSection]++
		historyStats.notEligibleFelonies++
		if conviction.OccurredInLast7Years() {
			historyStats.notEligibleFeloniesLast7Years++
		}

	case "Maybe Eligible - Flag for Review":
		clearanceByCodeSection.numberMaybeEligibleFlagForReviewByCodeSection[matchedCodeSection]++
		historyStats.maybeEligibleFelonies++
		if conviction.OccurredInLast7Years() {
			historyStats.maybeEligibleFeloniesLast7Years++
		}
	}

	clearanceByCodeSection.numberEligibilityByReason[eligibility.EligibilityReason]++
}

func (e *Exporter) incrementConvictionAndClearanceStats(
	county string,
	subject *data.Subject,
	prop64HistoryStats *historySummaryStats,
	relatedChargeHistoryStats *historySummaryStats,
	convictionStats *totalExistingConvictions,
	clearanceStats *totalClearanceResults,
) {

	convictionStats.totalCountyConvictions += subject.NumberOfConvictionsInCounty(county)
	if subject.NumberOfFelonies() > 0 {
		convictionStats.totalHasFelony++

		if subject.NumberOfFelonies() == (prop64HistoryStats.feloniesDismissed + prop64HistoryStats.feloniesReduced) {
			clearanceStats.numberNoLongerHaveFelony++
		}

		if subject.NumberOfFelonies() == (prop64HistoryStats.feloniesDismissed + prop64HistoryStats.feloniesReduced + prop64HistoryStats.maybeEligibleFelonies + prop64HistoryStats.notEligibleFelonies) {
			clearanceStats.numberNoLongerHaveFelonyIfAllSealed++
		}

		if subject.NumberOfFelonies() == (prop64HistoryStats.feloniesDismissed + prop64HistoryStats.feloniesReduced +
			prop64HistoryStats.maybeEligibleFelonies + prop64HistoryStats.notEligibleFelonies +
			relatedChargeHistoryStats.feloniesDismissed + relatedChargeHistoryStats.notEligibleFelonies) {
			clearanceStats.numberNoLongerHaveFelonyIfAllSealedIncludingRelated++
		}
	}

	if len(subject.Convictions) > 0 {
		convictionStats.totalHasConvictions++

		if len(subject.Convictions) == (prop64HistoryStats.feloniesDismissed + prop64HistoryStats.misdemeanorsDismissed) {
			clearanceStats.numberNoMoreConvictions++
		}

		if len(subject.Convictions) == (prop64HistoryStats.feloniesDismissed + prop64HistoryStats.feloniesReduced + prop64HistoryStats.maybeEligibleFelonies + prop64HistoryStats.notEligibleFelonies + prop64HistoryStats.misdemeanorsDismissed) {
			clearanceStats.numberNoMoreConvictionsIfAllSealed++
		}

		if len(subject.Convictions) == (prop64HistoryStats.feloniesDismissed + prop64HistoryStats.feloniesReduced +
			prop64HistoryStats.maybeEligibleFelonies + prop64HistoryStats.notEligibleFelonies + prop64HistoryStats.misdemeanorsDismissed +
			relatedChargeHistoryStats.feloniesDismissed + relatedChargeHistoryStats.misdemeanorsDismissed) {
			clearanceStats.numberNoMoreConvictionsIfAllSealedIncludingRelated++
		}
	}

	if subject.NumberOfConvictionsInLast7Years() > 0 {
		convictionStats.totalHasConvictionLast7Years++

		if subject.NumberOfConvictionsInLast7Years() == (prop64HistoryStats.feloniesDismissedLast7Years + prop64HistoryStats.misdemeanorsDismissedLast7Years) {
			clearanceStats.numberClearedRecordsLast7Years++
		}

		if subject.NumberOfConvictionsInLast7Years() == (prop64HistoryStats.feloniesDismissedLast7Years + prop64HistoryStats.feloniesReducedLast7Years + prop64HistoryStats.maybeEligibleFeloniesLast7Years + prop64HistoryStats.notEligibleFeloniesLast7Years + prop64HistoryStats.misdemeanorsDismissedLast7Years) {
			clearanceStats.numberClearedRecordsLast7YearsIfAllSealed++
		}

		if subject.NumberOfConvictionsInLast7Years() == (prop64HistoryStats.feloniesDismissedLast7Years + prop64HistoryStats.feloniesReducedLast7Years +
			prop64HistoryStats.maybeEligibleFeloniesLast7Years + prop64HistoryStats.notEligibleFeloniesLast7Years + prop64HistoryStats.misdemeanorsDismissedLast7Years +
			relatedChargeHistoryStats.feloniesDismissedLast7Years + relatedChargeHistoryStats.notEligibleFeloniesLast7Years + relatedChargeHistoryStats.misdemeanorsDismissedLast7Years) {
			clearanceStats.numberClearedRecordsLast7YearsIfAllSealedIncludingRelated++
		}
	}
}

func (e *Exporter) SummarizeAndExport(county string) {
	fmt.Printf("Processing Subjects\n")
	for _, subject := range e.dojInformation.Subjects {

		historyStats := historySummaryStats{}
		relatedChargeHistoryStats := historySummaryStats{}

		for _, conviction := range subject.Convictions {
			matchedCodeSection := data.EligibilityFlows[county].MatchedCodeSection(conviction.CodeSection)
			matchedRelatedCodeSection := data.EligibilityFlows[county].MatchedRelatedCodeSection(conviction.CodeSection)

			e.incrementConvictionsByCodeSection(conviction, county, matchedCodeSection, &e.convictionsByCodeSection)
			e.incrementConvictionsByCodeSection(conviction, county, matchedRelatedCodeSection, &e.relatedConvictionsByCodeSection)

			eligibility, ok := e.dojInformation.Eligibilities[conviction.Index]

			if ok && matchedCodeSection != "" {
				e.incrementEligibilities(eligibility, conviction, county, matchedCodeSection, &historyStats, &e.clearanceByCodeSection)
			}
			if ok && matchedRelatedCodeSection != "" {
				e.incrementEligibilities(eligibility, conviction, county, matchedRelatedCodeSection, &relatedChargeHistoryStats, &e.relatedChargeClearanceByCodeSection)
			}
		}

		e.incrementConvictionAndClearanceStats(county, subject, &historyStats, &relatedChargeHistoryStats, &e.totalExistingConvictions, &e.totalClearanceResults)
	}

	for i, row := range e.dojInformation.Rows {
		e.outputDOJWriter.WriteEntryWithEligibilityInfo(row, e.dojInformation.Eligibilities[i])
		e.outputCondensedDOJWriter.WriteCondensedEntryWithEligibilityInfo(row, e.dojInformation.Eligibilities[i])
		if e.dojInformation.Eligibilities[i] != nil {
			e.outputProp64ConvictionsDOJWriter.WriteEntryWithEligibilityInfo(row, e.dojInformation.Eligibilities[i])
		}
	}

	e.outputDOJWriter.Flush()
	e.outputCondensedDOJWriter.Flush()
	e.outputProp64ConvictionsDOJWriter.Flush()

	fmt.Println()
	fmt.Println("----------- Overall summary of DOJ file --------------------")
	fmt.Printf("Found %d Total rows in DOJ file\n", len(e.dojInformation.Rows))
	fmt.Printf("Found %d Total individuals in DOJ file\n", len(e.dojInformation.Subjects))
	fmt.Printf("Found %d Total convictions in DOJ file\n", e.dojInformation.TotalConvictions)
	fmt.Printf("Found %d convictions in this county\n", e.totalExistingConvictions.totalCountyConvictions)

	fmt.Println()
	fmt.Printf("----------- Prop64 Convictions Overall--------------------")
	printSummaryByCodeSection("total", e.convictionsByCodeSection.totalConvictionsByCodeSection)
	fmt.Println()
	fmt.Printf("----------- Prop64 Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", e.convictionsByCodeSection.countyConvictionsByCodeSection)
	printSummaryByCodeSection("that are eligible for dismissal", e.clearanceByCodeSection.numberDismissedByCodeSection)
	printSummaryByCodeSection("that are eligible for reduction", e.clearanceByCodeSection.numberReducedByCodeSection)
	printSummaryByCodeSection("that are flagged for review", e.clearanceByCodeSection.numberMaybeEligibleFlagForReviewByCodeSection)
	printSummaryByCodeSection("that are not eligible", e.clearanceByCodeSection.numberIneligibleByCodeSection)

	fmt.Println()
	fmt.Println("----------- Eligibility Reasons --------------------")
	printConvictionsCountByReason(e.clearanceByCodeSection.numberEligibilityByReason)
	fmt.Println()
	fmt.Printf("----------- Prop64 Related Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", e.relatedConvictionsByCodeSection.countyConvictionsByCodeSection)
	printSummaryByCodeSection("that are eligible for dismissal", e.relatedChargeClearanceByCodeSection.numberDismissedByCodeSection)
	printSummaryByCodeSection("that are flagged for review", e.relatedChargeClearanceByCodeSection.numberMaybeEligibleFlagForReviewByCodeSection)
	printSummaryByCodeSection("that are not eligible", e.relatedChargeClearanceByCodeSection.numberIneligibleByCodeSection)

	fmt.Println()
	fmt.Println("----------- Impact to individuals --------------------")
	fmt.Printf("%d individuals currently have a felony on their record\n", e.totalExistingConvictions.totalHasFelony)
	fmt.Printf("%d individuals currently have convictions on their record\n", e.totalExistingConvictions.totalHasConvictions)
	fmt.Printf("%d individuals currently have convictions on their record in the last 7 years\n", e.totalExistingConvictions.totalHasConvictionLast7Years)
	fmt.Println()
	fmt.Println("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", e.totalClearanceResults.numberNoLongerHaveFelony)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", e.totalClearanceResults.numberNoMoreConvictions)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", e.totalClearanceResults.numberClearedRecordsLast7Years)
	fmt.Println()
	fmt.Println("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", e.totalClearanceResults.numberNoLongerHaveFelonyIfAllSealed)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", e.totalClearanceResults.numberNoMoreConvictionsIfAllSealed)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", e.totalClearanceResults.numberClearedRecordsLast7YearsIfAllSealed)
	fmt.Println()
	fmt.Println("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", e.totalClearanceResults.numberNoLongerHaveFelonyIfAllSealedIncludingRelated)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", e.totalClearanceResults.numberNoMoreConvictionsIfAllSealedIncludingRelated)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", e.totalClearanceResults.numberClearedRecordsLast7YearsIfAllSealedIncludingRelated)

}

func printConvictionsCountByReason(numberEligibilityByReason map[string]int) {
	printMap("Found %d convictions in this county with eligibility reason: %s\n", numberEligibilityByReason)
}

func printSummaryByCodeSection(description string, resultsByCodeSection map[string]int) {
	fmt.Printf("\nFound %d convictions %s\n", sumValues(resultsByCodeSection), description)
	formatString := fmt.Sprintf("Found %%d %%s convictions %s\n", description)
	printMap(formatString, resultsByCodeSection)
}

func printMap(formatString string, values map[string]int) {
	keys := getSortedKeys(values)

	for _, key := range keys {
		fmt.Printf(formatString, values[key], key)
	}
}

func getSortedKeys(mapWithStringKeys map[string]int) []string {
	keys := make([]string, 0, len(mapWithStringKeys))
	for key := range mapWithStringKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sumValues(mapOfInts map[string]int) int {
	total := 0
	for _, value := range mapOfInts {
		total += value
	}
	return total
}
