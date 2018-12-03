package data

import (
	"strings"
	"time"
)

type DOJHistory struct {
	SubjectID         string
	Name              string
	WeakName          string
	CII               string
	DOB               time.Time
	SSN               string
	CDL               string
	PC290Registration bool
	Convictions       []*DOJRow
}

func (history *DOJHistory) PushRow(row DOJRow) {
	if history.SubjectID == "" {
		history.SubjectID = row.SubjectID
		history.Name = row.Name
		history.WeakName = strings.Split(row.Name, " ")[0]
		history.CII = row.CII
		history.DOB = row.DOB
		history.SSN = row.SSN
		history.CDL = row.CDL
	}

	if row.Convicted {
		history.Convictions = append(history.Convictions, &row)
	}

	if row.PC290Registration {
		history.PC290Registration = true
	}
}

func (history *DOJHistory) Match(entry CMSEntry) MatchData {
	var results = make(map[string]bool)

	if entry.CII != "" {
		cmsCII := entry.CII
		for len(cmsCII) < 8 {
			cmsCII = "0" + cmsCII
		}
		cmsCII = cmsCII[len(cmsCII)-8:]
		dojCII := history.CII[len(history.CII)-8:]
		results["cii"] = cmsCII == dojCII
	}

	results["ssn"] = entry.SSN != "" && entry.SSN == history.SSN
	results["cdl"] = entry.CDL != "" && entry.CDL == history.CDL

	if (entry.CourtNumber != "") {
		matched := false
		for _, row := range history.Convictions {
			if row.County == "SAN FRANCISCO" && row.MatchingCourtNumber(entry.CourtNumber) {
				matched = true
				break
			}
		}
		results["courtno"] = matched
	}

	name := entry.FormattedName()
	dateOfBirth := entry.DateOfBirth
	if (name != "" && dateOfBirth != time.Time{}) {
		results["nameAndDob"] = name == history.Name && dateOfBirth == history.DOB
		results["weakNameAndDob"] = history.matchWeakName(name) && dateOfBirth == history.DOB
	}

	matchStrength := 0
	for _, val := range results {
		if val {
			matchStrength++
		}
	}

	return MatchData{
		History:       history,
		MatchResults:  results,
		MatchStrength: matchStrength,
	}
}

func (history *DOJHistory) matchWeakName(formattedName string) bool {
	firstLast := strings.Split(formattedName, " ")[0]

	return firstLast == history.WeakName
}
