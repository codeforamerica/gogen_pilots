package processor

import (
	"fmt"
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
	}
	fmt.Printf("Found %d Total Convictions in DOJ file\n", d.convictionStats.totalConvictions)
	fmt.Printf("Found %d SAN FRANCISCO County Convictions in DOJ file\n", d.convictionStats.totalCountyConvictions)
	fmt.Printf("Found %d SAN FRANCISCO County Prop64 Convictions in DOJ file\n", d.convictionStats.totalProp64Convictions)

	//for _, value := range d.dojInformation.Eligibilities{
	//	fmt.Printf("eligibilities map %#v \n", value.EligibilityDetermination)
	//}
	for i, row := range d.dojInformation.Rows {
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
		}
	}
	d.outputDOJWriter.Flush()

	fmt.Printf("Found %d SAN FRANCISCO County Prop64 Convictions that are eligible for dismissal in DOJ file\n", d.clearanceStats.numberDismissedCounts)
	fmt.Printf("Found %d SAN FRANCISCO County Prop64 Convictions that are eligible for reduction in DOJ file\n", d.clearanceStats.numberReducedCounts)
	fmt.Printf("Found %d SAN FRANCISCO County Prop64 Convictions that are not eligible in DOJ file\n", d.clearanceStats.numberIneligibleCounts)
}
