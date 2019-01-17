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
	numberDismissed11357                      int
	numberDismissed11358                      int
	numberDismissed11359                      int
	numberDismissed11360                      int
	numberReduced11357                        int
	numberReduced11358                        int
	numberReduced11359                        int
	numberReduced11360                        int
}

type convictionStats struct {
	totalConvictions             int
	totalCountyConvictions       int
	totalCountyProp64Convictions int
	totalProp64Convictions       int
	total11357Convictions        int
	total11358Convictions        int
	total11359Convictions        int
	total11360Convictions        int
	county11357Convictions       int
	county11358Convictions       int
	county11359Convictions       int
	county11360Convictions       int
	numDOJConvictions            map[string]int
	DOJEligibilityByCodeSection  map[string]map[string]int
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
		prop64Matcher:   regexp.MustCompile(`(11357|11358|11359|11360).*`),
		convictionStats: convictionStats{
			numDOJConvictions:           make(map[string]int),
			DOJEligibilityByCodeSection: make(map[string]map[string]int),
		},
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

			if strings.HasPrefix(conviction.CodeSection, "11357") {
				d.convictionStats.total11357Convictions++
				if conviction.County == county {
					d.convictionStats.county11357Convictions++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Dismissal" {
					d.clearanceStats.numberDismissed11357++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Reduction" {
					d.clearanceStats.numberReduced11357++
				}
			}

			if strings.HasPrefix(conviction.CodeSection, "11358") {
				d.convictionStats.total11358Convictions++
				if conviction.County == county {
					d.convictionStats.county11358Convictions++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Dismissal" {
					d.clearanceStats.numberDismissed11358++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Reduction" {
					d.clearanceStats.numberReduced11358++
				}
			}

			if strings.HasPrefix(conviction.CodeSection, "11359") {
				d.convictionStats.total11359Convictions++
				if conviction.County == county {
					d.convictionStats.county11359Convictions++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Dismissal" {
					d.clearanceStats.numberDismissed11359++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Reduction" {
					d.clearanceStats.numberReduced11359++
				}
			}

			if strings.HasPrefix(conviction.CodeSection, "11360") {
				d.convictionStats.total11360Convictions++
				if conviction.County == county {
					d.convictionStats.county11360Convictions++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Dismissal" {
					d.clearanceStats.numberDismissed11360++
				}
				if ok && eligibility.EligibilityDetermination == "Eligible for Reduction" {
					d.clearanceStats.numberReduced11360++
				}

			}

			if ok {
				switch eligibility.EligibilityDetermination {
				case "Eligible for Dismissal":
					d.clearanceStats.numberDismissedCounts++

				case "Eligible for Reduction":
					d.clearanceStats.numberReducedCounts++

				case "Not eligible":
					d.clearanceStats.numberIneligibleCounts++
				}

				if time.Since(conviction.DispositionDate).Hours() <= 61320 {
					last7years = true
					totalConvictionsLast7Years++
				}

				switch eligibility.EligibilityReason {
				case "Misdemeanor or Infraction":
					misdemeanorsDismissed++
					d.clearanceStats.numberDismissedMisdemeanor++

					if last7years {
						misdemeanorsDismissedLast7Years++
					}

				case "Occurred after 11/09/2016":
					d.clearanceStats.numberNotEligibleNovNine16++

				case "HS 11357(b)":
					d.clearanceStats.numberDismissed11357b++

					feloniesDismissed++
					if last7years {
						feloniesDismissedLast7Years++
					}

				case "Final Conviction older than 10 years":
					d.clearanceStats.numberDismissedOlderThan10Years++

					feloniesDismissed++
					if last7years {
						feloniesDismissedLast7Years++
					}

				case "Later Convictions":
					d.clearanceStats.numberReducedLaterConvictions++

					feloniesReduced++

				case "Sentence not Completed":
					d.clearanceStats.numberReducedIncompleteSentence++
					d.clearanceStats.numberCheckSentencingData++

					feloniesReduced++

				case "Sentence Completed":
					d.clearanceStats.numberDismissedCompletedSentence++
					d.clearanceStats.numberCheckSentencingData++

					feloniesDismissed++
					if last7years {
						feloniesDismissedLast7Years++
					}
				}
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
