package processor

import (
	"fmt"
	"time"
	"gogen/data"
	"regexp"
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
	numberNoMoreConvictions int
	numberClearedRecordsLast7Years            int
	numberHistoriesWithConvictionInLast7Years int
	numberRecordsNoFelonies                   int
	numberHistoriesWithFelonies               int
}

type convictionStats struct {
	totalConvictions            int
	totalCountyConvictions      int
	totalProp64Convictions      int
	numDOJConvictions           map[string]int
	DOJEligibilityByCodeSection map[string]map[string]int
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
		d.convictionStats.totalConvictions += len(history.Convictions)
		d.convictionStats.totalCountyConvictions += history.NumberOfConvictionsInCounty(county)
		d.convictionStats.totalProp64Convictions += history.NumberOfProp64Convictions()
		var feloniesDismissed = 0
		var feloniesReduced = 0
		var misdemeanorsDismissed = 0
		var feloniesDismissedLast7Years = 0
		var misdemeanorsDismissedLast7Years = 0
		var totalConvictionsLast7Years = 0

		for _, conviction := range history.Convictions {
			var last7years = false
			eligibility, ok := d.dojInformation.Eligibilities[conviction.Index]

			if ok {
				if time.Since(conviction.DispositionDate).Hours() <= 61320 {
					last7years = true
					totalConvictionsLast7Years ++
					fmt.Printf("in last 7 years %#v", last7years)
				}

				switch eligibility.EligibilityReason {
				case "Misdemeanor or Infraction":
					misdemeanorsDismissed ++
					if last7years {
						misdemeanorsDismissedLast7Years ++
					}

				case "HS 11357(b)":
					feloniesDismissed ++
					if last7years {
						feloniesDismissedLast7Years ++
					}

				case "Final Conviction older than 10 years":
					feloniesDismissed ++
					if last7years {
						feloniesDismissedLast7Years ++
					}

				case "Later Convictions":
					feloniesReduced ++

				case "Sentence not Completed":
					feloniesReduced ++

				case "Sentence Completed":
					feloniesDismissed ++
					if last7years {
						feloniesDismissedLast7Years ++
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

	fmt.Printf("Found %d Total Convictions in DOJ file\n", d.convictionStats.totalConvictions)
	fmt.Printf("Found %d SAN FRANCISCO County Convictions in DOJ file\n", d.convictionStats.totalCountyConvictions)
	fmt.Printf("Found %d SAN FRANCISCO County Prop64 Convictions in DOJ file\n", d.convictionStats.totalProp64Convictions)

	//for _, value := range d.dojInformation.Eligibilities{
	//	fmt.Printf("eligibilities map %#v \n", value.EligibilityDetermination)
	//}
	for i, row := range d.dojInformation.Rows {
		numberDimissedFelonies := 0
		numberDimissedMisdemeanors := 0
		numberReducedFelonies := 0

		d.outputDOJWriter.WriteDOJEntry(row, d.dojInformation.Eligibilities[i])

		val, ok := d.dojInformation.Eligibilities[i]
		if ok {
			switch val.EligibilityDetermination {
			case "Eligible for Dismissal":
				d.clearanceStats.numberDismissedCounts ++

			case "Eligible for Reduction":
				d.clearanceStats.numberReducedCounts ++

			case "Not eligible":
				d.clearanceStats.numberIneligibleCounts ++
			}

			switch val.EligibilityReason {
			case "Misdemeanor or Infraction":
				numberDimissedMisdemeanors ++

				d.clearanceStats.numberDismissedMisdemeanor ++

			case "Occurred after 11/09/2016":
				d.clearanceStats.numberNotEligibleNovNine16 ++

			case "HS 11357(b)":
				numberDimissedFelonies ++
				d.clearanceStats.numberDismissed11357b ++

			case "Final Conviction older than 10 years":
				numberDimissedFelonies ++
				d.clearanceStats.numberDismissedOlderThan10Years ++

			case "Later Convictions":
				numberReducedFelonies ++
				d.clearanceStats.numberReducedLaterConvictions ++

			case "Sentence not Completed":
				numberReducedFelonies ++
				d.clearanceStats.numberReducedIncompleteSentence ++
				d.clearanceStats.numberCheckSentencingData ++

			case "Sentence Completed":
				numberDimissedFelonies ++
				d.clearanceStats.numberDismissedCompletedSentence ++
				d.clearanceStats.numberCheckSentencingData ++

			}

		}
	}
	d.outputDOJWriter.Flush()

	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissedCounts)
	fmt.Printf("Found %d Prop64 Convictions in this county that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReducedCounts)
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
