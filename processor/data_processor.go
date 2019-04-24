package processor

import (
	"fmt"
	"gogen/data"
	"regexp"
	"time"
)

type DataProcessor struct {
	dojInformation           				*data.DOJInformation
	outputDOJWriter          				DOJWriter
	outputCondensedDOJWriter 				DOJWriter
	prop64Matcher            				*regexp.Regexp
	stats                    				dataProcessorStats
	totalClearanceResults    				totalClearanceResults
	clearanceByCodeSection   				clearanceByCodeSection
	relatedChargeClearanceByCodeSection   	clearanceByCodeSection
	totalExistingConvictions 				totalExistingConvictions
	convictionsByCodeSection 				convictionsByCodeSection
	relatedConvictionsByCodeSection 		convictionsByCodeSection
}

type totalClearanceResults struct {
	numberFullyClearedRecords                     int
	numberClearedRecordsLast7Years                int
	numberClearedRecordsLast7YearsIfAllSealed     int
	numberHistoriesWithConvictionInLast7Years     int
	numberRecordsNoFelonies                       int
	numberHistoriesWithFelonies                   int
	numberNoLongerHaveFelony                      int
	numberNoLongerHaveFelonyIfAllSealed           int
	numberNoMoreConvictions                       int
	numberNoMoreConvictionsIfAllSealed            int
}

type clearanceByCodeSection struct {
	numberEligibilityByReason                     map[string]int
	numberDismissedByCodeSection                  map[string]int
	numberReducedByCodeSection                    map[string]int
	numberIneligibleByCodeSection                 map[string]int
	numberMaybeEligibleFlagForReviewByCodeSection map[string]int
}

type totalExistingConvictions struct {
	totalConvictions               int
	totalCountyConvictions         int
	totalHasFelony                 int
	totalHasConvictionLast7Years   int
	totalHasConvictions            int
}

type convictionsByCodeSection struct {
	totalConvictionsByCodeSection  map[string]int
	countyConvictionsByCodeSection map[string]int
	DOJEligibilityByCodeSection    map[string]map[string]int
}

type historySummaryStats struct {
	feloniesDismissed int
	feloniesReduced int
	maybeEligibleFelonies int
	notEligibleFelonies int
	misdemeanorsDismissed int
	feloniesDismissedLast7Years int
	feloniesReducedLast7Years int
	maybeEligibleFeloniesLast7Years int
	notEligibleFeloniesLast7Years int
	misdemeanorsDismissedLast7Years int
	totalConvictionsLast7Years int
}

type dataProcessorStats struct {
	nDOJProp64Convictions int
	nDOJSubjects          int
	nDOJFelonies          int
	nDOJMisdemeanors      int
}

func NewDataProcessor(
	dojInformation *data.DOJInformation,
	outputDOJWriter DOJWriter,
	outputCondensedDOJWriter DOJWriter,
) DataProcessor {
	return DataProcessor{
		dojInformation:  dojInformation,
		outputDOJWriter: outputDOJWriter,
		outputCondensedDOJWriter: outputCondensedDOJWriter,
		totalClearanceResults: totalClearanceResults{},
		clearanceByCodeSection: clearanceByCodeSection{
			numberEligibilityByReason:        make(map[string]int),
			numberDismissedByCodeSection:     make(map[string]int),
			numberReducedByCodeSection:       make(map[string]int),
			numberIneligibleByCodeSection:    make(map[string]int),
			numberMaybeEligibleFlagForReviewByCodeSection: make(map[string]int),
		},
		relatedChargeClearanceByCodeSection: clearanceByCodeSection{
			numberEligibilityByReason:        make(map[string]int),
			numberDismissedByCodeSection:     make(map[string]int),
			numberReducedByCodeSection:       make(map[string]int),
			numberIneligibleByCodeSection:    make(map[string]int),
			numberMaybeEligibleFlagForReviewByCodeSection: make(map[string]int),
		},
		totalExistingConvictions: totalExistingConvictions{},
		convictionsByCodeSection: convictionsByCodeSection{
			totalConvictionsByCodeSection:  make(map[string]int),
			countyConvictionsByCodeSection: make(map[string]int),
		},
		relatedConvictionsByCodeSection: convictionsByCodeSection{
			totalConvictionsByCodeSection:  make(map[string]int),
			countyConvictionsByCodeSection: make(map[string]int),
		},
	}
}

func (d *DataProcessor) incrementConvictionsByCodeSection(
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

func (d *DataProcessor) incrementEligibilities(
	eligibility *data.EligibilityInfo,
	conviction *data.DOJRow,
	county string,
	matchedCodeSection string,
	last7years bool,
	historyStats *historySummaryStats,
	clearanceByCodeSection *clearanceByCodeSection,
	) {
	switch eligibility.EligibilityDetermination {
		case "Eligible for Dismissal":
			clearanceByCodeSection.numberDismissedByCodeSection[matchedCodeSection]++
			if conviction.Felony {
				historyStats.feloniesDismissed++
				if last7years {
					historyStats.feloniesDismissedLast7Years++
				}
			} else {
				historyStats.misdemeanorsDismissed++
				if last7years {
					historyStats.misdemeanorsDismissedLast7Years++
				}
			}

		case "Eligible for Reduction":
			clearanceByCodeSection.numberReducedByCodeSection[matchedCodeSection]++
			historyStats.feloniesReduced++
			if last7years {
				historyStats.feloniesReducedLast7Years++
			}

		case "Not eligible":
			clearanceByCodeSection.numberIneligibleByCodeSection[matchedCodeSection]++
			historyStats.notEligibleFelonies++
			if last7years {
				historyStats.notEligibleFeloniesLast7Years++
			}

		case "Maybe Eligible - Flag for Review":
			clearanceByCodeSection.numberMaybeEligibleFlagForReviewByCodeSection[matchedCodeSection]++
			historyStats.maybeEligibleFelonies++
			if last7years {
				historyStats.maybeEligibleFeloniesLast7Years++
			}
		}

	clearanceByCodeSection.numberEligibilityByReason[eligibility.EligibilityReason]++
}

func (d *DataProcessor) incrementConvictionAndClearanceStats(
	county string,
	history *data.DOJHistory,
	historyStats *historySummaryStats,
	convictionStats *totalExistingConvictions,
	clearanceStats *totalClearanceResults,
	) {

	convictionStats.totalConvictions += len(history.Convictions)
	convictionStats.totalCountyConvictions += history.NumberOfConvictionsInCounty(county)

	if history.NumberOfFelonies() > 0 {
		convictionStats.totalHasFelony++

		if history.NumberOfFelonies() == (historyStats.feloniesDismissed + historyStats.feloniesReduced) {
			clearanceStats.numberNoLongerHaveFelony++
		}

		if history.NumberOfFelonies() == (historyStats.feloniesDismissed + historyStats.feloniesReduced + historyStats.maybeEligibleFelonies + historyStats.notEligibleFelonies) {
			clearanceStats.numberNoLongerHaveFelonyIfAllSealed++
		}
	}

	if len(history.Convictions) > 0 {
		convictionStats.totalHasConvictions++

		if len(history.Convictions) == (historyStats.feloniesDismissed + historyStats.misdemeanorsDismissed) {
			clearanceStats.numberNoMoreConvictions++
		}

		if len(history.Convictions) == (historyStats.feloniesDismissed + historyStats.feloniesReduced + historyStats.maybeEligibleFelonies + historyStats.notEligibleFelonies + historyStats.misdemeanorsDismissed) {
			clearanceStats.numberNoMoreConvictionsIfAllSealed++
		}
	}

	if historyStats.totalConvictionsLast7Years > 0 {
		convictionStats.totalHasConvictionLast7Years++

		if historyStats.totalConvictionsLast7Years == (historyStats.feloniesDismissedLast7Years + historyStats.misdemeanorsDismissedLast7Years) {
			clearanceStats.numberClearedRecordsLast7Years++
		}

		if historyStats.totalConvictionsLast7Years == (historyStats.feloniesDismissedLast7Years + historyStats.feloniesReducedLast7Years + historyStats.maybeEligibleFeloniesLast7Years + historyStats.notEligibleFeloniesLast7Years + historyStats.misdemeanorsDismissedLast7Years) {
			clearanceStats.numberClearedRecordsLast7YearsIfAllSealed++
		}
	}
}

func condenseRow(row []string) []string{
	var condensedRow  []string
	condensedRow = append(condensedRow, row[11])
	condensedRow = append(condensedRow, row[12])
	condensedRow = append(condensedRow, row[13])
	condensedRow = append(condensedRow, row[14])
	condensedRow = append(condensedRow, row[22])
	condensedRow = append(condensedRow, row[37])
	condensedRow = append(condensedRow, row[40])
	condensedRow = append(condensedRow, row[46])
	condensedRow = append(condensedRow, row[48])
	condensedRow = append(condensedRow, row[51])
	condensedRow = append(condensedRow, row[52])
	condensedRow = append(condensedRow, row[54])
	condensedRow = append(condensedRow, row[80])
	condensedRow = append(condensedRow, row[82])
	condensedRow = append(condensedRow, row[86])
	condensedRow = append(condensedRow, row[87])
	condensedRow = append(condensedRow, row[88])
	condensedRow = append(condensedRow, row[90])
	condensedRow = append(condensedRow, row[93])
	condensedRow = append(condensedRow, row[94])
	return condensedRow
}

func (d *DataProcessor) Process(county string) {
	fmt.Printf("Processing Histories\n")
	for _, history := range d.dojInformation.Histories {

		historyStats := historySummaryStats{}
		relatedChargeHistoryStats := historySummaryStats{}

		for _, conviction := range history.Convictions {
			var last7years = false
			if time.Since(conviction.DispositionDate).Hours() <= 61320 {
				last7years = true
				historyStats.totalConvictionsLast7Years++
			}

			matchedCodeSection := data.EligibilityFlows[county].MatchedCodeSection(conviction.CodeSection)
			matchedRelatedCodeSection := data.EligibilityFlows[county].MatchedRelatedCodeSection(conviction.CodeSection)

			d.incrementConvictionsByCodeSection(conviction, county, matchedCodeSection, &d.convictionsByCodeSection)
			d.incrementConvictionsByCodeSection(conviction, county, matchedRelatedCodeSection, &d.relatedConvictionsByCodeSection)

			eligibility, ok := d.dojInformation.Eligibilities[conviction.Index]

			if ok && matchedCodeSection != "" {
				d.incrementEligibilities(eligibility, conviction, county, matchedCodeSection, last7years, &historyStats, &d.clearanceByCodeSection)
			}
			if ok && matchedRelatedCodeSection != "" {
				d.incrementEligibilities(eligibility, conviction, county, matchedRelatedCodeSection, last7years, &relatedChargeHistoryStats, &d.relatedChargeClearanceByCodeSection)
			}
		}

		d.incrementConvictionAndClearanceStats(county, history, &historyStats, &d.totalExistingConvictions, &d.totalClearanceResults)
	}

	for i, row := range d.dojInformation.Rows {
		d.outputDOJWriter.WriteDOJEntry(row, d.dojInformation.Eligibilities[i])
		var condensedRow = condenseRow(row)
		d.outputCondensedDOJWriter.WriteDOJEntry(condensedRow, d.dojInformation.Eligibilities[i])
	}

	d.outputDOJWriter.Flush()
	d.outputCondensedDOJWriter.Flush()

	fmt.Println()
	fmt.Println("----------- Overall summary of DOJ file --------------------")
	fmt.Printf("Found %d Total rows in DOJ file\n", len(d.dojInformation.Rows))
	fmt.Printf("Found %d Total individuals in DOJ file\n", len(d.dojInformation.Histories))
	fmt.Printf("Found %d Total convictions in DOJ file\n", d.totalExistingConvictions.totalConvictions)
	fmt.Printf("Found %d convictions in this county\n", d.totalExistingConvictions.totalCountyConvictions)

	fmt.Println()
	fmt.Printf("----------- Prop64 Convictions Overall--------------------")
	printSummaryByCodeSection("total", d.convictionsByCodeSection.totalConvictionsByCodeSection)
	fmt.Println()
	fmt.Printf("----------- Prop64 Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", d.convictionsByCodeSection.countyConvictionsByCodeSection)
	printSummaryByCodeSection("that are eligible for dismissal", d.clearanceByCodeSection.numberDismissedByCodeSection)
	printSummaryByCodeSection("that are eligible for reduction", d.clearanceByCodeSection.numberReducedByCodeSection)
	printSummaryByCodeSection("that are flagged for review", d.clearanceByCodeSection.numberMaybeEligibleFlagForReviewByCodeSection)
	printSummaryByCodeSection("that are not eligible", d.clearanceByCodeSection.numberIneligibleByCodeSection)

	fmt.Println()
	fmt.Println("----------- Eligibility Reasons --------------------")
	for key, val := range d.clearanceByCodeSection.numberEligibilityByReason {
		fmt.Printf("Found %d convictions in this county with eligibility reason: %s\n", val, key)
	}
	fmt.Println()
	fmt.Printf("----------- Prop64 Related Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", d.relatedConvictionsByCodeSection.countyConvictionsByCodeSection)
	printSummaryByCodeSection("that are eligible for dismissal", d.relatedChargeClearanceByCodeSection.numberDismissedByCodeSection)
	printSummaryByCodeSection("that are flagged for review", d.relatedChargeClearanceByCodeSection.numberMaybeEligibleFlagForReviewByCodeSection)
	printSummaryByCodeSection("that are not eligible", d.relatedChargeClearanceByCodeSection.numberIneligibleByCodeSection)

	fmt.Println()
	fmt.Println("----------- Impact to individuals --------------------")
	fmt.Printf("%d individuals currently have a felony on their record\n", d.totalExistingConvictions.totalHasFelony)
	fmt.Printf("%d individuals currently have convictions on their record\n", d.totalExistingConvictions.totalHasConvictions)
	fmt.Printf("%d individuals currently have convictions on their record in the last 7 years\n", d.totalExistingConvictions.totalHasConvictionLast7Years)
	fmt.Println()
	fmt.Println("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.totalClearanceResults.numberNoLongerHaveFelony)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.totalClearanceResults.numberNoMoreConvictions)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.totalClearanceResults.numberClearedRecordsLast7Years)
	fmt.Println()
	fmt.Println("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------")

	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.totalClearanceResults.numberNoLongerHaveFelonyIfAllSealed)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.totalClearanceResults.numberNoMoreConvictionsIfAllSealed)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.totalClearanceResults.numberClearedRecordsLast7YearsIfAllSealed )

}

func printSummaryByCodeSection(description string, resultsByCodeSection map[string]int) {
	fmt.Printf("\nFound %d convictions %s\n", sumValues(resultsByCodeSection), description)
	for codeSection, number := range resultsByCodeSection {
		fmt.Printf("Found %d %s convictions %s\n", number, codeSection, description)
	}
}

func sumValues(mapOfInts map[string]int) int {
	total := 0
	for _, value := range mapOfInts {
		total += value
	}
	return total
}
