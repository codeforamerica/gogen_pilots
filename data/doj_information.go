package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gogen/utilities"
	"os"
	"time"
)

type DOJInformation struct {
	Rows             [][]string
	Histories        map[string]*DOJHistory
	Eligibilities    map[int]*EligibilityInfo
	comparisonTime   time.Time
	TotalConvictions int
	TotalConvictionsInCounty int
}

func (i *DOJInformation) generateHistories(county string) {
	currentRowIndex := 0.0
	totalRows := float64(len(i.Rows))

	fmt.Println("Reading DOJ Data Into Memory")

	var totalTime time.Duration = 0

	for index, row := range i.Rows {
		startTime := time.Now()
		dojRow := NewDOJRow(row, index)
		if i.Histories[dojRow.SubjectID] == nil {
			i.Histories[dojRow.SubjectID] = new(DOJHistory)
		}
		i.Histories[dojRow.SubjectID].PushRow(dojRow, county)
		currentRowIndex++

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(currentRowIndex, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
}

func (i *DOJInformation) determineEligibility(county string) {
	for _, history := range i.Histories {
		//history.computeEligibilities(i.Eligibilities, i.comparisonTime, county)
		infos := EligibilityFlows[county].ProcessHistory(history, i.comparisonTime)

		i.TotalConvictions += len(history.Convictions)
		i.TotalConvictionsInCounty += len(history.EligibilityInfos)

		for index, info := range infos {
			i.Eligibilities[index] = info
		}
	}
}

func NewDOJInformation(dojFileName string, comparisonTime time.Time, county string) (*DOJInformation, error) {
	dojFile, err := os.Open(dojFileName)
	if err != nil {
		panic(err)
	}

	bufferedReader := bufio.NewReader(dojFile)
	bufferedReader.ReadLine() // read and discard header row

	sourceCSV := csv.NewReader(bufferedReader)
	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}
	info := DOJInformation{
		Rows:           rows,
		Histories:      make(map[string]*DOJHistory),
		Eligibilities:  make(map[int]*EligibilityInfo),
		comparisonTime: comparisonTime,
	}

	info.generateHistories(county)
	info.determineEligibility(county)

	return &info, nil
}
