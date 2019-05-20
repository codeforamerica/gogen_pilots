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
	EligibilityDetermination       map[string]string
	EligibilityReason              string
	CaseNumber                     string
	Deceased                       string
}

type genericAllDismissedElegibility struct {
	county string
}

func (g genericAllDismissedElegibility) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow) {
	if EligibilityFlows[g.county].IsProp64Charge(row.CodeSection) {
		g.EligibleDismissal(info)
	}
}

func (g genericAllDismissedElegibility) EligibleDismissal(info *EligibilityInfo) {
	info.EligibilityDetermination["allDismissed"] = "Eligible for Dismissal"
}

type genericAllDismissedWithRelatedElegibility struct {
	county string
}

func (g genericAllDismissedWithRelatedElegibility) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow) {
	if EligibilityFlows[g.county].IsProp64Charge(row.CodeSection) || EligibilityFlows[g.county].IsRelatedCharge(row.CodeSection) {
		g.EligibleDismissal(info)
	}
}

func (g genericAllDismissedWithRelatedElegibility) EligibleDismissal(info *EligibilityInfo) {
	info.EligibilityDetermination["allDismissedRelated"] = "Eligible for Dismissal"
}

func NewEligibilityInfo(row *DOJRow, history *DOJHistory, comparisonTime time.Time, county string) *EligibilityInfo {
	info := new(EligibilityInfo)
	info.comparisonTime = comparisonTime
	info.EligibilityDetermination = make(map[string]string)

	if (row.DispositionDate == time.Time{}) {
		info.YearsSinceThisConviction = -1.0
	} else {
		info.YearsSinceThisConviction = info.yearsSinceEvent(row.DispositionDate)
	}

	if history.IsDeceased {
		info.Deceased = "Deceased"
	} else {
		info.Deceased = "-"
	}

	if history.PC290Registration {
		info.PC290Registration = "Yes"
	} else {
		info.PC290Registration = "-"
	}

	if len(history.PC290CodeSections()) > 0 {
		info.PC290CodeSections = strings.Join(history.PC290CodeSections(), ";")
	} else {
		info.PC290CodeSections = "-"
	}

	if len(history.SuperstrikeCodeSections()) > 0 {
		info.Superstrikes = strings.Join(history.SuperstrikeCodeSections(), ";")
	} else {
		info.Superstrikes = "-"
	}

	mostRecentConvictionDate := history.MostRecentConvictionDate()
	if (mostRecentConvictionDate == time.Time{}) {
		info.YearsSinceMostRecentConviction = -1.0
	} else {
		info.YearsSinceMostRecentConviction = info.yearsSinceEvent(mostRecentConvictionDate)
	}

	info.NumberOfConvictionsOnRecord = len(history.Convictions)
	info.NumberOfProp64Convictions = history.NumberOfProp64Convictions(county)
	info.DateOfConviction = row.DispositionDate
	info.CaseNumber = strings.Join(history.CaseNumbers[row.CountOrder[0:6]], "; ")

	EligibilityFlows[county].BeginEligibilityFlow(info, row)
	genericAllDismissedElegibility{county: county}.BeginEligibilityFlow(info, row)
	genericAllDismissedWithRelatedElegibility{county: county}.BeginEligibilityFlow(info, row)

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
