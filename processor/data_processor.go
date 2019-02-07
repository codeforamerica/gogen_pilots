package processor

import (
	"fmt"
	"gogen/data"
	"regexp"
	"strings"
	"time"
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
	numberDismissedCounts                     int
	numberReducedCounts                       int
	numberIneligibleCounts                    int
	numberClearedRecordsLast7Years            int
	numberHistoriesWithConvictionInLast7Years int
	numberRecordsNoFelonies                   int
	numberHistoriesWithFelonies               int
	numberEligibilityByReason                 map[string]int
	numberDismissedByCodeSection              map[string]int
	numberReducedByCodeSection                map[string]int
	numberNoLongerHaveFelony                  int
	numberNoMoreConvictions                   int
}

var Prop64CodeSections = []string{"11357", "11358", "11359", "11360"}

type convictionStats struct {
	totalConvictions               int
	totalCountyConvictions         int
	totalCountyProp64Convictions   int
	totalProp64Convictions         int
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
) DataProcessor {
	return DataProcessor{
		dojInformation:  dojInformation,
		outputDOJWriter: outputDOJWriter,
		clearanceStats: clearanceStats{
			numberEligibilityByReason:    make(map[string]int),
			numberDismissedByCodeSection: make(map[string]int),
			numberReducedByCodeSection:   make(map[string]int),
		},
		convictionStats: convictionStats{
			totalConvictionsByCodeSection:  make(map[string]int),
			countyConvictionsByCodeSection: make(map[string]int),
		},
	}
}

func (d *DataProcessor) incrementConvictions(conviction *data.DOJRow, county string) {
	for _, codeSection := range Prop64CodeSections {
		if strings.HasPrefix(conviction.CodeSection, codeSection) {
			d.convictionStats.totalConvictionsByCodeSection[codeSection]++
			if conviction.County == county {
				d.convictionStats.countyConvictionsByCodeSection[codeSection]++
			}
		}
	}

}

func (d *DataProcessor) incrementClearanceStats(conviction *data.DOJRow, determination string) {
	for _, codeSection := range Prop64CodeSections {
		if strings.HasPrefix(conviction.CodeSection, codeSection) {
			if determination == "Eligible for Dismissal" {
				d.clearanceStats.numberDismissedByCodeSection[codeSection]++
			}
			if determination == "Eligible for Reduction" {
				d.clearanceStats.numberReducedByCodeSection[codeSection]++
			}
		}
	}
}

func (d *DataProcessor) Process(county string) {
	fmt.Printf("Processing Histories\n")
	for _, history := range d.dojInformation.Histories {
		var feloniesDismissed = 0
		var feloniesReduced = 0
		var misdemeanorsDismissed = 0
		var feloniesDismissedLast7Years = 0
		var misdemeanorsDismissedLast7Years = 0
		var totalConvictionsLast7Years = 0

		d.convictionStats.totalConvictions += len(history.Convictions)
		d.convictionStats.totalCountyConvictions += history.NumberOfConvictionsInCounty(county)

		for _, conviction := range history.Convictions {
			var last7years = false
			eligibility, ok := d.dojInformation.Eligibilities[conviction.Index]

			d.incrementConvictions(conviction, county)
			if ok {
				d.incrementClearanceStats(conviction, eligibility.EligibilityDetermination)

				if time.Since(conviction.DispositionDate).Hours() <= 61320 {
					last7years = true
					totalConvictionsLast7Years++
				}

				switch eligibility.EligibilityDetermination {
				case "Eligible for Dismissal":
					d.clearanceStats.numberDismissedCounts++
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
					d.clearanceStats.numberReducedCounts++
					feloniesReduced++

				case "Not eligible":
					d.clearanceStats.numberIneligibleCounts++
				}

				d.clearanceStats.numberEligibilityByReason[eligibility.EligibilityReason]++
			}
		}

		if history.NumberOfFelonies() == (feloniesDismissed + feloniesReduced) {
			d.clearanceStats.numberNoLongerHaveFelony++
		}
		if len(history.Convictions) == (feloniesDismissed + misdemeanorsDismissed) {
			d.clearanceStats.numberNoMoreConvictions++
		}
		if totalConvictionsLast7Years == (feloniesDismissedLast7Years + misdemeanorsDismissedLast7Years) {
			d.clearanceStats.numberClearedRecordsLast7Years++
		}
	}

	for i, row := range d.dojInformation.Rows {
		d.outputDOJWriter.WriteDOJEntry(row, d.dojInformation.Eligibilities[i])
	}

	d.outputDOJWriter.Flush()

	for _, val := range d.convictionStats.totalConvictionsByCodeSection {
		d.convictionStats.totalProp64Convictions += val
	}

	for _, val := range d.convictionStats.countyConvictionsByCodeSection {
		d.convictionStats.totalCountyProp64Convictions += val
	}

	fmt.Printf("Found %d Total Convictions in DOJ file\n", d.convictionStats.totalConvictions)
	fmt.Printf("Found %d Total Prop64 Convictions in DOJ file\n", d.convictionStats.totalProp64Convictions)
	for _, codeSection := range Prop64CodeSections {
		fmt.Printf("Found %d HS %s Convictions total in DOJ file\n", d.convictionStats.totalConvictionsByCodeSection[codeSection], codeSection)
	}

	fmt.Printf("Found %d County Convictions in DOJ file\n", d.convictionStats.totalCountyConvictions)
	fmt.Printf("Found %d County Prop64 Convictions in DOJ file\n", d.convictionStats.totalCountyProp64Convictions)

	for _, codeSection := range Prop64CodeSections {
		fmt.Printf("Found %d HS %s Convictions in this county in DOJ file\n", d.convictionStats.countyConvictionsByCodeSection[codeSection], codeSection)
	}

	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissedCounts)

	for _, codeSection := range Prop64CodeSections {
		fmt.Printf("Found %d HS %s Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissedByCodeSection[codeSection], codeSection)
	}

	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReducedCounts)

	for _, codeSection := range Prop64CodeSections {
		fmt.Printf("Found %d HS %s Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReducedByCodeSection[codeSection], codeSection)
	}

	fmt.Printf("Found %d Prop64 Convictions in this county that are not eligible in DOJ file\n", d.clearanceStats.numberIneligibleCounts)

	for key, val := range d.clearanceStats.numberEligibilityByReason {
		fmt.Printf("Found %d Prop64 Convictions in this county with eligibility reason: %s\n", val, key)
	}

	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.clearanceStats.numberNoLongerHaveFelony)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.clearanceStats.numberNoMoreConvictions)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.clearanceStats.numberClearedRecordsLast7Years)
}
