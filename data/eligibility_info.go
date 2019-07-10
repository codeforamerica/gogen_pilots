package data

import (
	"gogen/matchers"
	"strings"
	"time"
)

type EligibilityInfo struct {
	NumberOfConvictionsOnRecord    int
	DateOfConviction               time.Time
	YearsSinceThisConviction       float64
	YearsSinceMostRecentConviction float64
	NumberOfProp64Convictions      int
	NumberOf11357Convictions       int
	NumberOf11358Convictions       int
	NumberOf11359Convictions       int
	NumberOf11360Convictions       int
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
	info.NumberOfProp64Convictions, info.NumberOf11357Convictions, info.NumberOf11358Convictions, info.NumberOf11359Convictions, info.NumberOf11360Convictions = subject.Prop64ConvictionsBySection()
	info.DateOfConviction = row.DispositionDate
	info.CaseNumber = strings.Join(subject.CaseNumbers[row.CountOrder[0:6]], "; ")

	return info
}

func (info *EligibilityInfo) yearsSinceEvent(date time.Time) float64 {
	hours := info.comparisonTime.Sub(date).Hours()
	years := hours / (24 * 365.25)
	return years
}

func (info *EligibilityInfo) hasSuperstrikes() bool {
	return info.Superstrikes != "-"
}

func (info *EligibilityInfo) hasTwoPriors(row *DOJRow, subject *Subject) bool {
	priorConvictionsOfSameCodeSectionPrefix := 0
	codeSectionRune := []rune(row.CodeSection)
	codeSectionPrefix := string(codeSectionRune[0:5])
	for _, conviction := range subject.Convictions {
		if matchers.IsProp64Charge(conviction.CodeSection) {
			if conviction.DispositionDate.Before(row.DispositionDate) {
				if strings.HasPrefix(conviction.CodeSection, codeSectionPrefix) {
					priorConvictionsOfSameCodeSectionPrefix++
				}
			}
		}
	}

	return priorConvictionsOfSameCodeSectionPrefix >= 2
}

func (info *EligibilityInfo) olderThanFifty(row *DOJRow, subject *Subject) bool {
	age := info.yearsSinceEvent(subject.DOB)
	if age >= 50 {
		return true
	}
	return false
}

func (info *EligibilityInfo) youngerThanTwentyOne(row *DOJRow, subject *Subject) bool {
	age := info.yearsSinceEvent(subject.DOB)
	if age <= 21 {
		return true
	}
	return false
}

func (info *EligibilityInfo) onlyProp64Convictions(row *DOJRow, subject *Subject) bool {
	return len(subject.Convictions) == info.NumberOfProp64Convictions
}

func (info *EligibilityInfo) allSentencesCompleted(row *DOJRow, subject *Subject) bool {
	for _, conviction := range subject.Convictions {
		if conviction.SentenceEndDate.After(info.comparisonTime) {
			return false
		}
	}
	return true
}

func (info *EligibilityInfo) noConvictionsPastTenYears(row *DOJRow, subject *Subject) bool {
	for _, conviction := range subject.Convictions {
		if conviction.DispositionDate.After(info.comparisonTime.AddDate(-10, 0, 0)) {
			return false
		}
	}
	return true
}

func (info *EligibilityInfo) SetEligibleForDismissal(reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (info *EligibilityInfo) SetEligibleForReduction(reason string) {
	info.EligibilityDetermination = "Eligible for Reduction"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (info *EligibilityInfo) SetNotEligible(reason string) {
	info.EligibilityDetermination = "Not eligible"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (info *EligibilityInfo) SetMaybeEligible(reason string) {
	info.EligibilityDetermination = "Maybe Eligible - Flag for Review"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (info *EligibilityInfo) SetHandReview(reason string) {
	info.EligibilityDetermination = "Hand Review"
	info.EligibilityReason = strings.TrimSpace(reason)
}
