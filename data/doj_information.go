package data

import (
	"encoding/csv"
	"fmt"
	"gogen/utilities"
	"time"
)

type DOJInformation struct {
	Rows      [][]string
	Histories map[string]*DOJHistory
}

func (information *DOJInformation) generateHistories() {
	currentRowIndex := 0.0
	totalRows := 486481.0

	fmt.Println("Reading DOJ Data Into Memory")

	var totalTime time.Duration = 0

	for _, row := range information.Rows {
		startTime := time.Now()
		dojRow := NewDOJRow(row)
		if information.Histories[dojRow.SubjectID] == nil {
			information.Histories[dojRow.SubjectID] = new(DOJHistory)
		}
		information.Histories[dojRow.SubjectID].PushRow(dojRow)
		currentRowIndex++

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(currentRowIndex, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
}

func NewDOJInformation(sourceCSV *csv.Reader) (*DOJInformation, error) {
	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}
	rows = rows[1:]
	info := DOJInformation{
		Rows:      rows,
		Histories: make(map[string]*DOJHistory),
	}

	info.generateHistories()
	return &info, nil
}
