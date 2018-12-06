package data

import (
	"encoding/csv"
	"fmt"
	"gogen/utilities"
	"time"
)

type SummaryMatchData struct {
	ambiguousMatches int
	matchCountByType map[string]int
	matchStrengths   map[int]int
}

type DOJInformation struct {
	Rows                [][]string
	Histories           map[string]*DOJHistory
	SummaryMatchData    SummaryMatchData
	weakNameAndDOBIndex map[string]*DOJHistory
	nameAndDOBIndex     map[string]*DOJHistory
	ciiIndex            map[string]*DOJHistory
	ssnIndex            map[string]*DOJHistory
	cdlIndex            map[string]*DOJHistory
	courtNumberIndex    map[string]*DOJHistory
}

type MatchData struct {
	History       *DOJHistory
	MatchResults  map[string]bool
	MatchStrength int
}

func (information *DOJInformation) FindDOJHistory(entry CMSEntry) (*DOJHistory, time.Duration) {
	var matches []MatchData
	var totalTime time.Duration = 0

	for _, history := range information.Histories {
		startTime := time.Now()
		matchData := history.Match(entry)
		if matchData.MatchStrength > 0 {
			matches = append(matches, matchData)
		}
		totalTime += time.Since(startTime)
	}

	if len(matches) == 0 {
		return nil, utilities.AverageTime(totalTime, float64(len(information.Histories)))
	}

	bestMatch := matches[0]
	if len(matches) > 1 {
		information.SummaryMatchData.ambiguousMatches++
		fmt.Println(fmt.Sprintf("Ambiguous match for `%s`", entry.FormattedName))
		for _, match := range matches {
			//TODO better printing for ambiguous matches
			fmt.Println(fmt.Sprintf("(name: `%s`, matches: %t): %+v", match.History.Name, entry.FormattedName == match.History.Name, match))
			if match.MatchStrength > bestMatch.MatchStrength {
				bestMatch = match
			}
		}
	}

	information.summarizeMatchData(bestMatch)
	return bestMatch.History, utilities.AverageTime(totalTime, float64(len(information.Histories)))
}

func (information *DOJInformation) summarizeMatchData(data MatchData) {
	for key, val := range data.MatchResults {
		if val {
			information.SummaryMatchData.matchCountByType[key]++
		}
	}
	information.SummaryMatchData.matchStrengths[data.MatchStrength]++
}

func (information *DOJInformation) generateIndexes() {
	information.weakNameAndDOBIndex = make(map[string]*DOJHistory)
	information.nameAndDOBIndex = make(map[string]*DOJHistory)
	information.ciiIndex = make(map[string]*DOJHistory)
	information.ssnIndex = make(map[string]*DOJHistory)
	information.cdlIndex = make(map[string]*DOJHistory)
	information.courtNumberIndex = make(map[string]*DOJHistory)


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

	info := DOJInformation{
		Rows:      rows,
		Histories: make(map[string]*DOJHistory),
	}

	info.generateHistories()
	info.generateIndexes()

	info.SummaryMatchData.matchCountByType = make(map[string]int)
	info.SummaryMatchData.matchStrengths = make(map[int]int)

	return &info, nil
}
