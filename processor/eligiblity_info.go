package processor

import (
	"fmt"
	"gogen/data"
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
}

const (
	eligible      = "eligible"
	ineligible    = "ineligible"
	notApplicable = "n/a"
	noMatch       = "no match"
	notFound      = "not found"
)

func (info *EligibilityInfo) checkWeight(entry data.CMSEntry, weightInfo data.WeightsEntry) {
	var eligibleString string

	if strings.HasPrefix(entry.Charge, "11357") || entry.Level == "M" {
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

func (info *EligibilityInfo) checkDOJHistory(entry data.CMSEntry, history *data.DOJHistory) {
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
		hours := time.Since(mostRecentConvictionDate).Hours()
		years := hours / (24 * 265.25)
		info.YearsSinceMostRecentConviction = fmt.Sprintf("%.1f", years)
	}

	if entry.Level != "F" || strings.HasPrefix(entry.Charge, "11357") {
		result = notApplicable

		info.PC290Registration = result
		info.PC290Charges = result
		info.PC290CodeSections = result
		info.Superstrikes = result
		info.SuperstrikeCodeSections = result
		info.TwoPriors = result
		info.AgeAtConviction = result
		info.YearsSinceEvent = result
		info.YearsSinceMostRecentConviction = result
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
	}

	superstrikes := history.SuperstrikesCodeSections()
	if len(superstrikes) > 0 {
		info.Superstrikes = ineligible
		info.SuperstrikeCodeSections = strings.Join(superstrikes, "; ")
	} else {
		info.Superstrikes = eligible
	}

	if history.ThreeConvictionsSameCode(entry.Charge) {
		info.TwoPriors = ineligible
	} else {
		info.TwoPriors = eligible
	}
}

func (info *EligibilityInfo) computeFinalEligibility() {
	disqualifiers := info.Over1Lb == ineligible ||
		info.PC290Registration == ineligible ||
		info.PC290Charges == ineligible ||
		info.Superstrikes == ineligible ||
		info.TwoPriors == ineligible

	if disqualifiers {
		info.FinalRecommendation = ineligible
		return
	}

	needsReview := info.Over1Lb == noMatch || info.PC290Registration == noMatch

	if needsReview {
		info.FinalRecommendation = "needs review"
		return
	}

	info.FinalRecommendation = eligible
}

func NewEligibilityInfo(entry data.CMSEntry, weightInfo data.WeightsEntry, history *data.DOJHistory) *EligibilityInfo {
	eligibilityInfo := new(EligibilityInfo)

	if (entry.DateOfBirth == time.Time{} || entry.DispositionDate == time.Time{}) {
		eligibilityInfo.AgeAtConviction = notFound
	} else {
		hours := entry.DispositionDate.Sub(entry.DateOfBirth).Hours()
		years := int(hours / (24 * 265.25))
		eligibilityInfo.YearsSinceMostRecentConviction = fmt.Sprintf("%d", years)
	}

	eligibilityInfo.checkWeight(entry, weightInfo)
	eligibilityInfo.checkDOJHistory(entry, history)
	eligibilityInfo.computeFinalEligibility()
	return eligibilityInfo
}
