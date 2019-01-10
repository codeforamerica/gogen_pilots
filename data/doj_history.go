package data

import (
	"sort"
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
	seenConvictions   map[string]bool
	OriginalCII       string
}

func (history *DOJHistory) PushRow(row DOJRow) {
	if history.SubjectID == "" {
		history.SubjectID = row.SubjectID
		history.Name = row.Name
		history.WeakName = strings.Split(row.Name, " ")[0]
		history.CII = row.CII
		history.OriginalCII = row.CII
		history.DOB = row.DOB
		history.SSN = row.SSN
		history.CDL = row.CDL
		history.seenConvictions = make(map[string]bool)
	}

	if row.Convicted && !history.seenConvictions[row.CountOrder] {
		history.Convictions = append(history.Convictions, &row)
		history.seenConvictions[row.CountOrder] = true
	}

	if row.PC290Registration {
		history.PC290Registration = true
	}
}

func (history *DOJHistory) MostRecentConvictionDate() time.Time {
	if len(history.Convictions) == 0 {
		return time.Time{}
	}
	convictions := history.Convictions
	sort.Slice(convictions, func(i, j int) bool {
		return convictions[i].DispositionDate.Before(convictions[j].DispositionDate)
	})
	return convictions[len(convictions)-1].DispositionDate
}

func (history *DOJHistory) NumberOfProp64Convictions() int {
	result := 0
	for _, row := range history.Convictions {
		if IsProp64Charge(row.CodeSection) {
			result++
		}
	}
	return result
}

func (history *DOJHistory) computeEligibilities(infos map[int]*EligibilityInfo, comparisonTime time.Time) {
	for _, row := range history.Convictions {
		if IsProp64Charge(row.CodeSection) {
			infos[row.Index] = NewEligibilityInfo(row, history, comparisonTime)
		}
	}
}
