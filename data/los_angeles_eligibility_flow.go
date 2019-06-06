package data

import (
	"gogen/matchers"
	"regexp"
	"strings"
	"time"
)

type losAngelesEligibilityFlow struct {
	prop64Matcher *regexp.Regexp
}

func (ef losAngelesEligibilityFlow) ProcessHistory(history *DOJHistory, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range history.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, history, comparisonTime, "LOS ANGELES")
			ef.BeginEligibilityFlow(info, conviction, history)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef losAngelesEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

func (ef losAngelesEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == "LOS ANGELES" && ef.IsProp64Charge(codeSection)
}

func (ef losAngelesEligibilityFlow) IsProp64Charge(codeSection string) bool {
	ok, _ := matchers.Prop64Matcher(codeSection)
	return ok
}

func (ef losAngelesEligibilityFlow) MatchedCodeSection(codeSection string) string {
	matches := ef.prop64Matcher.FindStringSubmatch(codeSection)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func (ef losAngelesEligibilityFlow) MatchedRelatedCodeSection(codeSection string) string {
	return ""
}

func (ef losAngelesEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = reason
}

func (ef losAngelesEligibilityFlow) EligibleReduction(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Reduction"
	info.EligibilityReason = reason
}

func (ef losAngelesEligibilityFlow) NotEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Not eligible"
	info.EligibilityReason = reason
}

func (ef losAngelesEligibilityFlow) MaybeEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Maybe Eligible - Flag for Review"
	info.EligibilityReason = reason
}

func (ef losAngelesEligibilityFlow) HandReview(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Hand Review"
	info.EligibilityReason = reason
}

func (ef losAngelesEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if ef.IsProp64Charge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.ConvictionIs11357(info, row, history)
	} else {
		ef.NotEligible(info, "Occurred after 11/09/2016")
	}
}

func (ef losAngelesEligibilityFlow) ConvictionIs11357(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if ef.MatchedCodeSection(row.CodeSection) == "11357" {
		if strings.HasPrefix(row.CodeSection, "11357(A)") || strings.HasPrefix(row.CodeSection, "11357(B)") {
			ef.EligibleDismissal(info, "11357(a) or 11357(b)")
		} else {
			ef.MaybeEligible(info, "Other 11357")
		}
	} else {
		ef.HasSuperstrikes(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) HasSuperstrikes(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.hasSuperstrikes() {
		ef.NotEligible(info, "PC 667(e)(2)(c)(iv)")
	} else {
		ef.TwoPriors(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) TwoPriors(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.hasTwoPriors(row, history) {
		ef.NotEligible(info, "Two priors")
	} else {
		ef.OlderThanFifty(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) OlderThanFifty(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.olderThanFifty(row, history) {
		ef.EligibleDismissal(info, "50 years or older")
	} else {
		ef.YoungerThanTwentyOne(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) YoungerThanTwentyOne(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.youngerThanTwentyOne(row, history) {
		ef.EligibleDismissal(info, "21 years or younger")
	} else {
		ef.Prop64OnlyWithCompletedSentences(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) Prop64OnlyWithCompletedSentences(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.onlyProp64Convictions(row, history) && info.allSentencesCompleted(row, history) {
		ef.EligibleDismissal(info, "Only has 11357-60 charges and completed sentence")
	} else {
		ef.NoConvictionsPastTenYears(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) NoConvictionsPastTenYears(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if info.noConvictionsPastTenYears(row, history) {
		ef.EligibleDismissal(info, "No convictions in past 10 years")
	} else {
		ef.ServingSentence(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) ServingSentence(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if !info.allSentencesCompleted(row, history) {
		ef.HandReview(info, "Currently serving sentence")
	} else {
		ef.IsDeceased(info, row, history)
	}
}

func (ef losAngelesEligibilityFlow) IsDeceased(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if history.IsDeceased  {
		ef.EligibleDismissal(info, "Deceased")
	} else {
		ef.HandReview(info, "????")
	}
}

