package data

import (
	"strings"
	"time"
)

type EligibilityInfo struct {
	NumberOfConvictionsOnRecord    int
	DateOfConviction               time.Time
	YearsSinceThisConviction       float64
	YearsSinceMostRecentConviction float64
	NumberOfProp64Convictions      int
	comparisonTime                 time.Time
	Superstrikes                   string
	PC290CodeSections              string
	PC290Registration              string
	EligibilityDetermination       string
	EligibilityReason              string
	CaseNumber                     string
	Deceased                       string
}

func NewEligibilityInfo(row *DOJRow, subject *Subject, comparisonTime time.Time, county string) *EligibilityInfo {
	info := new(EligibilityInfo)
	info.comparisonTime = comparisonTime

	if (row.DispositionDate == time.Time{}) {
		info.YearsSinceThisConviction = -1.0
	} else {
		info.YearsSinceThisConviction = info.yearsSinceEvent(row.DispositionDate)
	}

	if subject.IsDeceased {
		info.Deceased = "Deceased"
	} else {
		info.Deceased = "-"
	}

	if subject.PC290Registration {
		info.PC290Registration = "Yes"
	} else {
		info.PC290Registration = "-"
	}

	if len(subject.PC290CodeSections()) > 0 {
		info.PC290CodeSections = strings.Join(subject.PC290CodeSections(), ";")
	} else {
		info.PC290CodeSections = "-"
	}

	if len(subject.SuperstrikeCodeSections()) > 0 {
		info.Superstrikes = strings.Join(subject.SuperstrikeCodeSections(), ";")
	} else {
		info.Superstrikes = "-"
	}

	mostRecentConvictionDate := subject.MostRecentConvictionDate()
	if (mostRecentConvictionDate == time.Time{}) {
		info.YearsSinceMostRecentConviction = -1.0
	} else {
		info.YearsSinceMostRecentConviction = info.yearsSinceEvent(mostRecentConvictionDate)
	}

	info.NumberOfConvictionsOnRecord = len(subject.Convictions)
	info.NumberOfProp64Convictions = subject.NumberOfProp64Convictions(county)
	info.DateOfConviction = row.DispositionDate
	info.CaseNumber = strings.Join(subject.CaseNumbers[row.CountOrder[0:6]], "; ")

	EligibilityFlows[county].BeginEligibilityFlow(info, row)

	if info.EligibilityReason != "" {
		return info
	} else {
		return nil
	}
}

func (info *EligibilityInfo) yearsSinceEvent(date time.Time) float64 {
	hours := info.comparisonTime.Sub(date).Hours()
	years := hours / (24 * 365.25)
	return years
}

func (info *EligibilityInfo) hasSuperstrikes() bool {
	return info.Superstrikes != "-"
}
