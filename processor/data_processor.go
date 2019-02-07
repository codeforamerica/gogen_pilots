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
	numberDismissedMisdemeanor                int
	numberDismissed11357b                     int
	numberDismissedOlderThan10Years           int
	numberReducedLaterConvictions             int
	numberReducedIncompleteSentence           int
	numberDismissedCompletedSentence          int
	numberNotEligibleNovNine16                int
	numberNoLongerHaveFelony                  int
	numberCheckSentencingData                 int
	numberNoMoreConvictions                   int
	numberClearedRecordsLast7Years            int
	numberHistoriesWithConvictionInLast7Years int
	numberRecordsNoFelonies                   int
	numberHistoriesWithFelonies               int
	numberEligibilityByReason                 map[string]int
	numberDismissedByCodeSection              map[string]int
	numberReducedByCodeSection                map[string]int
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
	if determination == "Eligible for Dismissal" {
		d.clearanceStats.numberDismissedByCodeSection[conviction.CodeSection[:5]]++
	}
	if determination == "Eligible for Reduction" {
		d.clearanceStats.numberReducedByCodeSection[conviction.CodeSection[:5]]++
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

	d.convictionStats.totalProp64Convictions = d.convictionStats.total11357Convictions + d.convictionStats.total11358Convictions + d.convictionStats.total11359Convictions + d.convictionStats.total11360Convictions
	d.convictionStats.totalCountyProp64Convictions = d.convictionStats.county11357Convictions + d.convictionStats.county11358Convictions + d.convictionStats.county11359Convictions + d.convictionStats.county11360Convictions
	fmt.Printf("Found %d Total Convictions in DOJ file\n", d.convictionStats.totalConvictions)
	fmt.Printf("Found %d Total Prop64 Convictions in DOJ file\n", d.convictionStats.totalProp64Convictions)
	fmt.Printf("Found %d HS 11357 Convictions total in DOJ file\n", d.convictionStats.total11357Convictions)
	fmt.Printf("Found %d HS 11358 Convictions total in DOJ file\n", d.convictionStats.total11358Convictions)
	fmt.Printf("Found %d HS 11359 Convictions total in DOJ file\n", d.convictionStats.total11359Convictions)
	fmt.Printf("Found %d HS 11360 Convictions total in DOJ file\n", d.convictionStats.total11360Convictions)

	fmt.Printf("Found %d County Convictions in DOJ file\n", d.convictionStats.totalCountyConvictions)
	fmt.Printf("Found %d County Prop64 Convictions in DOJ file\n", d.convictionStats.totalCountyProp64Convictions)
	fmt.Printf("Found %d HS 11357 Convictions in this county in DOJ file\n", d.convictionStats.county11357Convictions)
	fmt.Printf("Found %d HS 11358 Convictions in this county in DOJ file\n", d.convictionStats.county11358Convictions)
	fmt.Printf("Found %d HS 11359 Convictions in this county in DOJ file\n", d.convictionStats.county11359Convictions)
	fmt.Printf("Found %d HS 11360 Convictions in this county in DOJ file\n", d.convictionStats.county11360Convictions)

	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissedCounts)
	fmt.Printf("Found %d HS 11357 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissed11357)
	fmt.Printf("Found %d HS 11358 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissed11358)
	fmt.Printf("Found %d HS 11359 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissed11359)
	fmt.Printf("Found %d HS 11360 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissed11360)

	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReducedCounts)
	fmt.Printf("Found %d HS 11357 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReduced11357)
	fmt.Printf("Found %d HS 11358 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReduced11358)
	fmt.Printf("Found %d HS 11359 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReduced11359)
	fmt.Printf("Found %d HS 11360 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReduced11360)

	fmt.Printf("Found %d Prop64 Convictions in this county that are not eligible in DOJ file\n", d.clearanceStats.numberIneligibleCounts)

	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file because of Misdemeanor or Infraction\n", d.clearanceStats.numberDismissedMisdemeanor)
	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file because of HS 11357b\n", d.clearanceStats.numberDismissed11357b)
	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file because final conviction older than 10 years\n", d.clearanceStats.numberDismissedOlderThan10Years)
	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for reduction in DOJ file because there are later convictions\n", d.clearanceStats.numberReducedLaterConvictions)
	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for reduction in DOJ file because they did not complete their sentence\n", d.clearanceStats.numberReducedIncompleteSentence)
	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file because they completed their sentence\n", d.clearanceStats.numberDismissedCompletedSentence)
	fmt.Printf("Found %d Prop64 Convictions in this county that are not eligible because after November 9 2016\n", d.clearanceStats.numberNotEligibleNovNine16)

	fmt.Printf("Found %d Prop64 Convictions in this county that need sentence data checked\n", d.clearanceStats.numberCheckSentencingData)

	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.clearanceStats.numberNoLongerHaveFelony)
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.clearanceStats.numberNoMoreConvictions)
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.clearanceStats.numberClearedRecordsLast7Years)
}
