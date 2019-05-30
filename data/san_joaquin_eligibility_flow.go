package data

import (
	"gogen/matchers"
	"regexp"
	"strings"
	"time"
)

type sanJoaquinEligibilityFlow struct {
	prop64Matcher        *regexp.Regexp
	relatedChargeMatcher *regexp.Regexp
}

func (ef sanJoaquinEligibilityFlow) ProcessHistory(history *DOJHistory, comparisonTime time.Time) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range history.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, history, comparisonTime, "SAN JOAQUIN")
			ef.BeginEligibilityFlow(info, conviction)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef sanJoaquinEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == "SAN JOAQUIN" && (ef.IsProp64Charge(codeSection) || ef.isRelatedCharge(codeSection))
}

func (ef sanJoaquinEligibilityFlow) IsProp64Charge(codeSection string) bool {
	ok, _ := matchers.Prop64Matcher(codeSection)
	return ok
}

func (ef sanJoaquinEligibilityFlow) MatchedCodeSection(codeSection string) string {
	prop64Matches := ef.prop64Matcher.FindStringSubmatch(codeSection)
	if len(prop64Matches) > 0 {
		return prop64Matches[1]
	}

	return ""
}

func (ef sanJoaquinEligibilityFlow) MatchedRelatedCodeSection(codeSection string) string {
	relatedChargeMatches := ef.relatedChargeMatcher.FindStringSubmatch(codeSection)
	if len(relatedChargeMatches) > 0 {
		return relatedChargeMatches[1]
	}

	return ""
}

func (ef sanJoaquinEligibilityFlow) isRelatedCharge(codeSection string) bool {
	ok, _ := matchers.RelatedChargeMatcher(codeSection)
	return ok
}

func (ef sanJoaquinEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow) {
	if ef.IsProp64Charge(row.CodeSection) || ef.isRelatedCharge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef sanJoaquinEligibilityFlow) EligibleReduction(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Reduction"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef sanJoaquinEligibilityFlow) MaybeEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Maybe Eligible - Flag for Review"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef sanJoaquinEligibilityFlow) NotEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Not eligible"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef sanJoaquinEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.convictionIsRelatedCharge(info, row)
	} else {
		ef.NotEligible(info, "Occurred after 11/09/2016")
	}
}

func (ef sanJoaquinEligibilityFlow) convictionIsRelatedCharge(info *EligibilityInfo, row *DOJRow) {
	if ef.isRelatedCharge(row.CodeSection) {
		ef.hasProp64ChargeInCycle(info, row)
	} else {
		ef.ConvictionIsNotFelony(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) ConvictionIsNotFelony(info *EligibilityInfo, row *DOJRow) {
	if !row.IsFelony {
		ef.EligibleDismissal(info, "Misdemeanor or Infraction")
	} else {
		ef.Is11357(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) Is11357(info *EligibilityInfo, row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357") {
		ef.EligibleDismissal(info, "11357 HS")
	} else {
		ef.MoreThanOneConviction(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) MoreThanOneConviction(info *EligibilityInfo, row *DOJRow) {
	if info.NumberOfConvictionsOnRecord > 1 {
		ef.HasConvictionsInPast5Years(info, row)
	} else {
		ef.CurrentlyServingSentence(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) HasConvictionsInPast5Years(info *EligibilityInfo, row *DOJRow) {
	if info.YearsSinceMostRecentConviction > 5 {
		ef.EligibleDismissal(info, "No convictions in past 5 years")
	} else {
		ef.MaybeEligible(info, "Has convictions in past 5 years")
	}
}

func (ef sanJoaquinEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		ef.EligibleDismissal(info, "Sentence Completed")
	} else {
		ef.MaybeEligible(info, "Sentence not Completed")
	}
}

func (ef sanJoaquinEligibilityFlow) hasProp64ChargeInCycle(info *EligibilityInfo, row *DOJRow) {
	if row.HasProp64ChargeInCycle {
		ef.isBP4060Charge(info, row)
	} else {
		ef.MaybeEligible(info, "No Related Prop64 Charges")
	}
}

func (ef sanJoaquinEligibilityFlow) isBP4060Charge(info *EligibilityInfo, row *DOJRow) {
	if row.CodeSection == "4060 BP" {
		ef.MaybeEligible(info, "4060 BP")
	} else {
		ef.EligibleDismissal(info, row.CodeSection)
	}
}
