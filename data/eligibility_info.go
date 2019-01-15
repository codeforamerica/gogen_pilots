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
	EligibilityDetermination       string
	EligibilityReason              string
}

func NewEligibilityInfo(row *DOJRow, history *DOJHistory, comparisonTime time.Time) *EligibilityInfo {
	info := new(EligibilityInfo)
	info.comparisonTime = comparisonTime

	if (row.DispositionDate == time.Time{}) {
		info.YearsSinceThisConviction = -1.0
	} else {
		info.YearsSinceThisConviction = info.yearsSinceEvent(row.DispositionDate)
	}

	mostRecentConvictionDate := history.MostRecentConvictionDate()
	if (mostRecentConvictionDate == time.Time{}) {
		info.YearsSinceMostRecentConviction = -1.0
	} else {
		info.YearsSinceMostRecentConviction = info.yearsSinceEvent(mostRecentConvictionDate)
	}

	info.NumberOfConvictionsOnRecord = len(history.Convictions)
	info.NumberOfProp64Convictions = history.NumberOfProp64Convictions()
	info.DateOfConviction = row.DispositionDate

	info.BeginEligibilityFlow(row)

	return info
}

func (info *EligibilityInfo) yearsSinceEvent(date time.Time) float64 {
	hours := info.comparisonTime.Sub(date).Hours()
	years := hours / (24 * 365.25)
	return years
}

func (info *EligibilityInfo) BeginEligibilityFlow(row *DOJRow) {
	info.ConvictionBeforeNovNine2016(row)
}

func (info *EligibilityInfo) EligibleDismissal(reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = reason
}

func (info *EligibilityInfo) EligibleReduction(reason string) {
	info.EligibilityDetermination = "Eligible for Reduction"
	info.EligibilityReason = reason
}

func (info *EligibilityInfo) NotEligible(reason string) {
	info.EligibilityDetermination = "Not eligible"
	info.EligibilityReason = reason
}

func (info *EligibilityInfo) ConvictionBeforeNovNine2016(row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		info.ConvictionIsNotFelony(row)
	} else {
		info.NotEligible("Occurred after 11/09/2016")
	}
}

func (info *EligibilityInfo) ConvictionIsNotFelony(row *DOJRow) {
	if !row.Felony {
		info.EligibleDismissal("Misdemeanor or Infraction")
	} else {
		info.Is11357b(row)
	}
}

func (info *EligibilityInfo) Is11357b(row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357(b)") {
		info.EligibleDismissal("HS 11357(b)")
	} else {
		info.MoreThanOneConviction(row)
	}
}

func (info *EligibilityInfo) MoreThanOneConviction(row *DOJRow) {
	if info.NumberOfConvictionsOnRecord > 1 {
		info.ThisConvictionOlderThan10Years(row)
	} else {
		info.CurrentlyServingSentence(row)
	}
}

func (info *EligibilityInfo) ThisConvictionOlderThan10Years(row *DOJRow) {
	if info.YearsSinceThisConviction > 10 {
		info.FinalConvictionOnRecord(row)
	} else {
		info.EligibleReduction("Occurred in last 10 years")
	}
}

func (info *EligibilityInfo) CurrentlyServingSentence(row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		info.EligibleDismissal("Sentence Completed")
	} else {
		info.EligibleReduction("Sentence not Completed")
	}
}

func (info *EligibilityInfo) FinalConvictionOnRecord(row *DOJRow) {
	if info.YearsSinceMostRecentConviction == info.YearsSinceThisConviction {
		info.EligibleDismissal("Final Conviction older than 10 years")
	} else {
		info.EligibleReduction("Later Convictions")
	}
}
