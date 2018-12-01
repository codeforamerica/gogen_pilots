package data

import "time"

type DOJHistory struct {
	SubjectID         string
	Name              string
	CII               string
	DOB               time.Time
	SSN               string
	CDL               string
	PC290Registration bool
	Convictions       []DOJRow
}

func (history DOJHistory) PushRow(row DOJRow) {
	if history.SubjectID == "" {
		history.SubjectID = row.SubjectID
		history.Name = row.Name
		history.CII = row.CII
		history.DOB = row.DOB
		history.SSN = row.SSN
		history.CDL = row.CDL
	}

	if (row.Convicted) {
		history.Convictions = append(history.Convictions, row)
	}

	if (row.PC290Registration) {
		history.PC290Registration = true
	}
}

func (history *DOJHistory) match(entry CMSEntry) MatchData {
	var results = make(map[string]bool)
	results["ssn"] = entry.SSN != "" && entry.SSN == history.SSN
	results["cdl"] = entry.CDL != "" && entry.CDL == history.CDL

	matchStrength := 0
	for _, val := range results {
		if val {
			matchStrength++
		}
	}

	return MatchData{
		history:       history,
		matchResults:  results,
		matchStrength: matchStrength,
	}
}
