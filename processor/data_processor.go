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
	numberClearedRecordsLast7Years            int
	numberHistoriesWithConvictionInLast7Years int
	numberRecordsNoFelonies                   int
	numberHistoriesWithFelonies               int
}

type convictionStats struct {
	totalConvictions int
	numDOJConvictions           map[string]int
	DOJEligibilityByCodeSection map[string]map[string]int
}

type dataProcessorStats struct {
	nDOJProp64Convictions    int
	nDOJSubjects             int
	nDOJFelonies             int
	nDOJMisdemeanors         int
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
func (d *DataProcessor) Process() {
	for _, history := range d.dojInformation.Histories {
		fmt.Printf(history.Name + " ")
		fmt.Println(len(history.Convictions))
		d.convictionStats.totalConvictions += len(history.Convictions)
	}

	for i, row := range d.dojInformation.Rows {
		d.outputDOJWriter.WriteDOJEntry(row, d.dojInformation.Eligibilities[i])
	}
	d.outputDOJWriter.Flush()

	fmt.Printf("Found %d Total Convictions in DOJ file\n", d.convictionStats.totalConvictions)
}

