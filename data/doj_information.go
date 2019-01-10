package data

import (
	"encoding/csv"
	"fmt"
	"gogen/utilities"
	"time"
)

type DOJInformation struct {
	Rows           [][]string
	Histories      map[string]*DOJHistory
	Eligibilities  map[int]*EligibilityInfo
	comparisonTime time.Time
}

func (i *DOJInformation) generateHistories() {
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
		i.Histories[dojRow.SubjectID].PushRow(dojRow)
		currentRowIndex++

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(currentRowIndex, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
}

func (i *DOJInformation) determineEligibility() {
	for _, history := range i.Histories {
		history.computeEligibilities(i.Eligibilities, i.comparisonTime)
	}
}

func NewDOJInformation(sourceCSV *csv.Reader, comparisonTime time.Time) (*DOJInformation, error) {
	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}
	rows = rows[1:]
	info := DOJInformation{
		Rows:           rows,
		Histories:      make(map[string]*DOJHistory),
		Eligibilities:  make(map[int]*EligibilityInfo),
		comparisonTime: comparisonTime,
	}

	info.generateHistories()
	info.determineEligibility()
	return &info, nil
}
