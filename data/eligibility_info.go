package data

import (
	. "gogen/matchers"
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

func NewEligibilityInfo(row *DOJRow, history *DOJHistory, comparisonTime time.Time, county string) *EligibilityInfo {
	info := new(EligibilityInfo)
	info.comparisonTime = comparisonTime

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

func (info *EligibilityInfo) hasTwoPriors(row *DOJRow, history *DOJHistory) bool {
	priorConvictionsOfSameCodeSectionPrefix := 0
	codeSectionRune := []rune(row.CodeSection)
	codeSectionPrefix := string(codeSectionRune[0:5])
	for _, conviction := range history.Convictions {
		prop64Conviction, _ := Prop64Matcher(conviction.CodeSection)
		if prop64Conviction {
			if conviction.DispositionDate.Before(row.DispositionDate) {
				if strings.HasPrefix(conviction.CodeSection, codeSectionPrefix) {
					priorConvictionsOfSameCodeSectionPrefix++
				}
			}
		}
	}

	return priorConvictionsOfSameCodeSectionPrefix >= 2
}

func (info *EligibilityInfo) olderThanFifty(row *DOJRow, history *DOJHistory) bool {
	age := info.yearsSinceEvent(history.DOB)
	if age >= 50 {
		return true
	}
	return false
}

func (info *EligibilityInfo) youngerThanTwentyOne(row *DOJRow, history *DOJHistory) bool {
	age := info.yearsSinceEvent(history.DOB)
	if age <= 21 {
		return true
	}
	return false
}

func (info *EligibilityInfo) onlyProp64Convictions(row *DOJRow, history *DOJHistory) bool {
	return len(history.Convictions) == info.NumberOfProp64Convictions
}

func (info *EligibilityInfo) allSentencesCompleted(row *DOJRow, history *DOJHistory) bool {
	for _, conviction := range history.Convictions {
		if conviction.SentenceEndDate.After(info.comparisonTime){
			return false
		}
	}
	return true
}

func (info *EligibilityInfo) noConvictionsPastTenYears(row *DOJRow, history *DOJHistory) bool {
	for _, conviction := range history.Convictions {
		if conviction.DispositionDate.After(info.comparisonTime.AddDate(-10,0,0)) {
			return false
		}
	}
	return true
}