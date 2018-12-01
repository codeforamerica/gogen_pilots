package data

import (
	"encoding/csv"
	"fmt"
)

type SummaryMatchData struct {
	ambiguousMatches int
	matchCountByType map[string]int
	matchStrengths   map[int]int
}

type DOJInformation struct {
	Rows             [][]string
	Histories        map[string]DOJHistory
	SummaryMatchData SummaryMatchData
}

type MatchData struct {
	history       *DOJHistory
	matchResults  map[string]bool
	matchStrength int
}

func (information *DOJInformation) findDOJHistory(entry CMSEntry) *DOJHistory {
	var matches []MatchData
	for _, history := range (information.Histories) {
		matchData := history.match(entry)
		if matchData.matchStrength > 0 {
			matches = append(matches, matchData)
		}
	}

	if len(matches) == 0 {
		return nil
	}

	bestMatch := matches[0]
	if len(matches) > 1 {
		information.SummaryMatchData.ambiguousMatches++
		fmt.Print("Ambiguous match!")
		for _, match := range (matches) {
			//TODO better printing for ambiguous matches
			if match.matchStrength > bestMatch.matchStrength {
				bestMatch = match
			}
		}
	}

	information.summarizeMatchData(bestMatch)
	return bestMatch.history
}

func (information *DOJInformation) summarizeMatchData(data MatchData) {
	for key, val := range data.matchResults {
		if val {
			information.SummaryMatchData.matchCountByType[key]++
		}
		information.SummaryMatchData.matchStrengths[data.matchStrength]++
	}
}

func NewDOJInformation(sourceCSV *csv.Reader) (*DOJInformation, error) {
	const SubjectIDIndex int = 0

	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}

	info := DOJInformation{
		Rows:      rows,
		Histories: make(map[string]DOJHistory),
	}

	for _, row := range rows {
		info.Histories[row[SubjectIDIndex]].PushRow(NewDOJRow(row))
	}

	return &info, nil
}
