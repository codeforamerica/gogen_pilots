package processor

import (
	"fmt"
	"gogen/data"
	"regexp"
	"strings"
	"time"
)

type EligibilityInfo struct {
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

	disqualifiers := info.PC290Registration == ineligible ||
		info.PC290Charges == ineligible ||
		info.Superstrikes == ineligible ||
		info.TwoPriors == ineligible

	if disqualifiers {
		info.FinalRecommendation = ineligible
		return
	}

	convictionNeedsReview := info.PC290Registration == noMatch

	if convictionNeedsReview {
		info.FinalRecommendation = needsReview
		return
	}

	info.FinalRecommendation = eligible
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
	eligibilityInfo.checkDOJHistory(row.CodeSection, level, history)
	eligibilityInfo.computeFinalEligibility(row.CodeSection, prop64Matcher)
	return eligibilityInfo
}
