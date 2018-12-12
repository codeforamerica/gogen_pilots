package processor

import (
	"fmt"
	"gogen/data"
	"regexp"
	"strings"
	"time"
)

type EligibilityInfo struct {
	QFinalSum                      string
	Over1Lb                        string
	PC290Registration              string
	PC290Charges                   string
	PC290CodeSections              string
	Superstrikes                   string
	SuperstrikeCodeSections        string
	TwoPriors                      string
	AgeAtConviction                string
	YearsSinceEvent                string
	YearsSinceMostRecentConviction string
	FinalRecommendation            string
	Prop64Charge                   string
	comparisonTime                 time.Time
}

const (
	eligible      = "eligible"
	ineligible    = "ineligible"
	notApplicable = "n/a"
	noMatch       = "no match"
	notFound      = "not found"
	needsReview   = "needs review"
)

func (info *EligibilityInfo) checkWeight(charge string, level string, weightInfo data.WeightsEntry) {
	var eligibleString string

	if strings.HasPrefix(charge, "11357") || strings.HasPrefix(charge, "11358") || level == "M" {
		info.QFinalSum = notApplicable
		info.Over1Lb = notApplicable
		return
	}

	if !weightInfo.Found {
		info.QFinalSum = noMatch
		info.Over1Lb = noMatch
		return
	}

	if weightInfo.Weight <= 453.592 {
		eligibleString = eligible
	} else {
		eligibleString = ineligible
	}
	info.QFinalSum = fmt.Sprintf("%.1f", weightInfo.Weight)
	info.Over1Lb = eligibleString
}

func (info *EligibilityInfo) checkDOJHistory(charge string, level string, history *data.DOJHistory) {
	result := ""
	if history == nil {
		result = noMatch

		info.PC290Registration = result
		info.PC290Charges = result
		info.PC290CodeSections = result
		info.Superstrikes = result
		info.SuperstrikeCodeSections = result
		info.TwoPriors = result
		info.YearsSinceMostRecentConviction = result
		return
	}

	mostRecentConvictionDate := history.MostRecentConvictionDate()
	if (mostRecentConvictionDate == time.Time{}) {
		info.YearsSinceMostRecentConviction = notFound
	} else {
		info.YearsSinceMostRecentConviction = info.yearsSinceEvent(mostRecentConvictionDate)
	}

	if level != "F" || strings.HasPrefix(charge, "11357") {
		result = notApplicable

		info.PC290Registration = result
		info.PC290Charges = result
		info.PC290CodeSections = result
		info.Superstrikes = result
		info.SuperstrikeCodeSections = result
		info.TwoPriors = result
		return
	}

	if history.PC290Registration {
		info.PC290Registration = ineligible
	} else {
		info.PC290Registration = eligible
	}

	pc290 := history.PC290CodeSections()
	if len(pc290) > 0 {
		info.PC290Charges = ineligible
		info.PC290CodeSections = strings.Join(pc290, "; ")
	} else {
		info.PC290Charges = eligible
		info.PC290CodeSections = "-"
	}

	superstrikes := history.SuperstrikesCodeSections()
	if len(superstrikes) > 0 {
		info.Superstrikes = ineligible
		info.SuperstrikeCodeSections = strings.Join(superstrikes, "; ")
	} else {
		info.Superstrikes = eligible
		info.SuperstrikeCodeSections = "-"
	}

	if history.ThreeConvictionsSameCode(charge) {
		info.TwoPriors = ineligible
	} else {
		info.TwoPriors = eligible
	}
}

func (info *EligibilityInfo) yearsSinceEvent(date time.Time) string {
	hours := info.comparisonTime.Sub(date).Hours()
	years := hours / (24 * 365.25)
	return fmt.Sprintf("%.1f", years)
}

func (info *EligibilityInfo) computeFinalEligibility(charge string, prop64Matcher *regexp.Regexp) {
	if prop64Matcher == nil {
		prop64Matcher = regexp.MustCompile("")
	}

	if !prop64Matcher.Match([]byte(charge)) {
		info.Prop64Charge = ineligible
		info.FinalRecommendation = ineligible
		return
	}

	disqualifiers := info.Over1Lb == ineligible ||
		info.PC290Registration == ineligible ||
		info.PC290Charges == ineligible ||
		info.Superstrikes == ineligible ||
		info.TwoPriors == ineligible

	if disqualifiers {
		info.FinalRecommendation = ineligible
		return
	}

	convictionNeedsReview := info.Over1Lb == noMatch || info.PC290Registration == noMatch

	if convictionNeedsReview {
		info.FinalRecommendation = needsReview
		return
	}

	info.FinalRecommendation = eligible
}

func NewEligibilityInfo(entry data.CMSEntry, weightInfo data.WeightsEntry, history *data.DOJHistory, comparisonTime time.Time, prop64Matcher *regexp.Regexp) *EligibilityInfo {
	eligibilityInfo := new(EligibilityInfo)
	eligibilityInfo.comparisonTime = comparisonTime

	if (entry.DateOfBirth == time.Time{} || entry.DispositionDate == time.Time{}) {
		eligibilityInfo.AgeAtConviction = notFound
	} else {
		hours := entry.DispositionDate.Sub(entry.DateOfBirth).Hours()
		years := int(hours / (24 * 365.25))
		eligibilityInfo.AgeAtConviction = fmt.Sprintf("%d", years)
	}

	if (entry.DispositionDate == time.Time{}) {
		eligibilityInfo.YearsSinceEvent = notFound
	} else {
		eligibilityInfo.YearsSinceEvent = eligibilityInfo.yearsSinceEvent(entry.DispositionDate)
	}

	eligibilityInfo.checkWeight(entry.Charge, entry.Level, weightInfo)
	eligibilityInfo.checkDOJHistory(entry.Charge, entry.Level, history)
	eligibilityInfo.computeFinalEligibility(entry.Charge, prop64Matcher)
	return eligibilityInfo
}

func EligibilityInfoFromDOJRow(row *data.DOJRow, history *data.DOJHistory, comparisonTime time.Time, prop64Matcher *regexp.Regexp) *EligibilityInfo {
	eligibilityInfo := new(EligibilityInfo)
	eligibilityInfo.comparisonTime = comparisonTime

	if (history.DOB == time.Time{} || row.DispositionDate == time.Time{}) {
		eligibilityInfo.AgeAtConviction = notFound
	} else {
		hours := row.DispositionDate.Sub(history.DOB).Hours()
		years := int(hours / (24 * 365.25))
		eligibilityInfo.AgeAtConviction = fmt.Sprintf("%d", years)
	}

	if (row.DispositionDate == time.Time{}) {
		eligibilityInfo.YearsSinceEvent = notFound
	} else {
		eligibilityInfo.YearsSinceEvent = eligibilityInfo.yearsSinceEvent(row.DispositionDate)
	}

	level := "M"
	if row.Felony {
		level = "F"
	}
	eligibilityInfo.checkWeight(row.CodeSection, level, data.WeightsEntry{Weight: 0, Found: false})
	eligibilityInfo.checkDOJHistory(row.CodeSection, level, history)
	eligibilityInfo.computeFinalEligibility(row.CodeSection, prop64Matcher)
	return eligibilityInfo
}
