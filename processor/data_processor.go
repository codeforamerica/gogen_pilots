package processor

import (
	"fmt"
	"gogen/data"
	"regexp"
	"time"
)

type DataProcessor struct {
	dojInformation           *data.DOJInformation
	outputDOJWriter          DOJWriter
	outputCondensedDOJWriter DOJWriter
	prop64Matcher            *regexp.Regexp
	stats                    dataProcessorStats
	clearanceStats           clearanceStats
	convictionStats          convictionStats
}

type clearanceStats struct {
	numberFullyClearedRecords                     int
	numberClearedRecordsLast7Years                int
	numberClearedRecordsLast7YearsIfAllSealed     int
	numberHistoriesWithConvictionInLast7Years     int
	numberRecordsNoFelonies                       int
	numberHistoriesWithFelonies                   int
	numberEligibilityByReason                     map[string]int
	numberDismissedByCodeSection                  map[string]int
	numberReducedByCodeSection                    map[string]int
	numberIneligibleByCodeSection                 map[string]int
	numberMaybeEligibleFlagForReviewByCodeSection map[string]int
	numberNoLongerHaveFelony                      int
	numberNoLongerHaveFelonyIfAllSealed           int
	numberNoMoreConvictions                       int
	numberNoMoreConvictionsIfAllSealed            int
}

type convictionStats struct {
	totalConvictions               int
	totalCountyConvictions         int
	totalHasFelony                 int
	totalHasConvictionLast7Years   int
	totalHasConvictions            int
	totalConvictionsByCodeSection  map[string]int
	countyConvictionsByCodeSection map[string]int
	DOJEligibilityByCodeSection    map[string]map[string]int
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
		clearanceStats: clearanceStats{
			numberEligibilityByReason:        make(map[string]int),
			numberDismissedByCodeSection:     make(map[string]int),
			numberReducedByCodeSection:       make(map[string]int),
			numberIneligibleByCodeSection:    make(map[string]int),
			numberMaybeEligibleFlagForReviewByCodeSection: make(map[string]int),
		},
		convictionStats: convictionStats{
			totalConvictionsByCodeSection:  make(map[string]int),
			countyConvictionsByCodeSection: make(map[string]int),
		},
	}
}

func (d *DataProcessor) incrementConvictions(conviction *data.DOJRow, county string, matchedCodeSection string) {
	if matchedCodeSection != "" {
		d.convictionStats.totalConvictionsByCodeSection[matchedCodeSection]++
		if conviction.County == county {
			d.convictionStats.countyConvictionsByCodeSection[matchedCodeSection]++
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
		var feloniesDismissed = 0
		var feloniesReduced = 0
		var maybeEligibleFelonies = 0
		var notEligibleFelonies = 0
		var misdemeanorsDismissed = 0
		var feloniesDismissedLast7Years = 0
		var feloniesReducedLast7Years = 0
		var maybeEligibleFeloniesLast7Years = 0
		var notEligibleFeloniesLast7Years = 0
		var misdemeanorsDismissedLast7Years = 0
		var totalConvictionsLast7Years = 0

		d.convictionStats.totalConvictions += len(history.Convictions)
		d.convictionStats.totalCountyConvictions += history.NumberOfConvictionsInCounty(county)

		for _, conviction := range history.Convictions {
			var last7years = false
			if time.Since(conviction.DispositionDate).Hours() <= 61320 {
				last7years = true
				totalConvictionsLast7Years++
			}

			matchedCodeSection := data.EligibilityFlows[county].MatchedCodeSection(conviction.CodeSection)

			eligibility, ok := d.dojInformation.Eligibilities[conviction.Index]

			d.incrementConvictions(conviction, county, matchedCodeSection)
			if ok && matchedCodeSection != "" {
				switch eligibility.EligibilityDetermination {
				case "Eligible for Dismissal":
					d.clearanceStats.numberDismissedByCodeSection[matchedCodeSection]++
					if conviction.Felony {
						feloniesDismissed++
						if last7years {
							feloniesDismissedLast7Years++
						}
					} else {
						misdemeanorsDismissed++
						if last7years {
							misdemeanorsDismissedLast7Years++
						}
					}

				case "Eligible for Reduction":
					d.clearanceStats.numberReducedByCodeSection[matchedCodeSection]++
					feloniesReduced++
					if last7years {
						feloniesReducedLast7Years++
					}

				case "Not eligible":
					d.clearanceStats.numberIneligibleByCodeSection[matchedCodeSection]++
					notEligibleFelonies++
					if last7years {
						notEligibleFeloniesLast7Years++
					}

				case "Maybe Eligible - Flag for Review":
					d.clearanceStats.numberMaybeEligibleFlagForReviewByCodeSection[matchedCodeSection]++
					maybeEligibleFelonies++
					if last7years {
						maybeEligibleFeloniesLast7Years++
					}
				}

				d.clearanceStats.numberEligibilityByReason[eligibility.EligibilityReason]++
			}
		}

		if history.NumberOfFelonies() > 0 {
			d.convictionStats.totalHasFelony++

			if history.NumberOfFelonies() == (feloniesDismissed + feloniesReduced) {
				d.clearanceStats.numberNoLongerHaveFelony++
			}

			if history.NumberOfFelonies() == (feloniesDismissed + feloniesReduced + maybeEligibleFelonies + notEligibleFelonies) {
				d.clearanceStats.numberNoLongerHaveFelonyIfAllSealed++
			}
		}

		if len(history.Convictions) > 0 {
			d.convictionStats.totalHasConvictions++

			if len(history.Convictions) == (feloniesDismissed + misdemeanorsDismissed) {
				d.clearanceStats.numberNoMoreConvictions++
			}

			if len(history.Convictions) == (feloniesDismissed + feloniesReduced + maybeEligibleFelonies + notEligibleFelonies + misdemeanorsDismissed) {
				d.clearanceStats.numberNoMoreConvictionsIfAllSealed++
			}
		}

		if totalConvictionsLast7Years > 0 {
			d.convictionStats.totalHasConvictionLast7Years++

			if totalConvictionsLast7Years == (feloniesDismissedLast7Years + misdemeanorsDismissedLast7Years) {
				d.clearanceStats.numberClearedRecordsLast7Years++
			}

			if totalConvictionsLast7Years == (feloniesDismissedLast7Years + feloniesReducedLast7Years + maybeEligibleFeloniesLast7Years + notEligibleFeloniesLast7Years + misdemeanorsDismissedLast7Years) {
				d.clearanceStats.numberClearedRecordsLast7YearsIfAllSealed++
			}
		}
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
	fmt.Printf("Found %d Total convictions in DOJ file\n", d.convictionStats.totalConvictions)
	fmt.Printf("Found %d convictions in this county\n", d.convictionStats.totalCountyConvictions)

	fmt.Println()
	fmt.Printf("----------- Prop64 and Related Convictions Overall--------------------")
	printSummaryByCodeSection("total", d.convictionStats.totalConvictionsByCodeSection)
	fmt.Println()
	fmt.Printf("----------- Prop64 and Related Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", d.convictionStats.countyConvictionsByCodeSection)
	printSummaryByCodeSection("that are eligible for dismissal", d.clearanceStats.numberDismissedByCodeSection)
	printSummaryByCodeSection("that are eligible for reduction", d.clearanceStats.numberReducedByCodeSection)
	printSummaryByCodeSection("that are flagged for review", d.clearanceStats.numberMaybeEligibleFlagForReviewByCodeSection)
	printSummaryByCodeSection("that are not eligible", d.clearanceStats.numberIneligibleByCodeSection)

	fmt.Println()
	fmt.Println("----------- Eligibility Reasons --------------------")
	for key, val := range d.clearanceStats.numberEligibilityByReason {
		fmt.Printf("Found %d convictions in this county with eligibility reason: %s\n", val, key)
	}
	fmt.Println()
	fmt.Println("----------- Impact to individuals --------------------")
	fmt.Printf("%d individuals currently have a felony on their record\n", d.convictionStats.totalHasFelony)
	fmt.Printf("%d individuals currently have convictions on their record\n", d.convictionStats.totalHasConvictions)
	fmt.Printf("%d individuals currently have convictions on their record in the last 7 years\n", d.convictionStats.totalHasConvictionLast7Years)
	fmt.Println()
	fmt.Println("----------- Based on the current eligibility logic --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.clearanceStats.numberNoLongerHaveFelony)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.clearanceStats.numberNoMoreConvictions)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.clearanceStats.numberClearedRecordsLast7Years)
	fmt.Println()
	fmt.Println("----------- If all convictions are dismissed and sealed--------------------")

	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.clearanceStats.numberNoLongerHaveFelonyIfAllSealed)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.clearanceStats.numberNoMoreConvictionsIfAllSealed)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.clearanceStats.numberClearedRecordsLast7YearsIfAllSealed )

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
