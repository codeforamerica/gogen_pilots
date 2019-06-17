package data

import (
	"gogen/matchers"
	"regexp"
	"strings"
	"time"
)

type contraCostaEligibilityFlow struct {
	prop64Matcher        *regexp.Regexp
	relatedChargeMatcher *regexp.Regexp
}

func (ef contraCostaEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, "CONTRA COSTA")
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}
func (ef contraCostaEligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef contraCostaEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == "CONTRA COSTA" && (matchers.IsProp64Charge(codeSection) || matchers.IsRelatedCharge(codeSection))
}

func (ef contraCostaEligibilityFlow) IsProp64Charge(codeSection string) bool {
	return matchers.IsProp64Charge(codeSection)
}

func (ef contraCostaEligibilityFlow) MatchedRelatedCodeSection(codeSection string) string {
	relatedChargeMatches := ef.relatedChargeMatcher.FindStringSubmatch(codeSection)
	if len(relatedChargeMatches) > 0 {
		return relatedChargeMatches[1]
	}

	return ""
}

func (ef contraCostaEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) || matchers.IsRelatedCharge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row)
	}
}

func (ef contraCostaEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef contraCostaEligibilityFlow) EligibleReduction(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Reduction"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef contraCostaEligibilityFlow) MaybeEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Maybe Eligible - Flag for Review"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef contraCostaEligibilityFlow) NotEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Not eligible"
	info.EligibilityReason = strings.TrimSpace(reason)
}

func (ef contraCostaEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.convictionIsRelatedCharge(info, row)
	} else {
		ef.NotEligible(info, "Occurred after 11/09/2016")
	}
}

func (ef contraCostaEligibilityFlow) convictionIsRelatedCharge(info *EligibilityInfo, row *DOJRow) {
	if matchers.IsRelatedCharge(row.CodeSection) {
		ef.hasProp64ChargeInCycle(info, row)
	} else {
		ef.ConvictionIsNotFelony(info, row)
	}
}

func (ef contraCostaEligibilityFlow) ConvictionIsNotFelony(info *EligibilityInfo, row *DOJRow) {
	if !row.IsFelony {
		ef.EligibleDismissal(info, "Misdemeanor or Infraction")
	} else {
		ef.Is11357(info, row)
	}
}

func (ef contraCostaEligibilityFlow) Is11357(info *EligibilityInfo, row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357") {
		ef.EligibleDismissal(info, "11357 HS")
	} else {
		ef.MoreThanOneConviction(info, row)
	}
}

func (ef contraCostaEligibilityFlow) MoreThanOneConviction(info *EligibilityInfo, row *DOJRow) {
	if info.NumberOfConvictionsOnRecord > 1 {
		ef.HasConvictionsInPast5Years(info, row)
	} else {
		ef.CurrentlyServingSentence(info, row)
	}
}

func (ef contraCostaEligibilityFlow) HasConvictionsInPast5Years(info *EligibilityInfo, row *DOJRow) {
	if info.YearsSinceMostRecentConviction > 5 {
		ef.EligibleDismissal(info, "No convictions in past 5 years")
	} else {
		ef.MaybeEligible(info, "Has convictions in past 5 years")
	}
}

func (ef contraCostaEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		ef.EligibleDismissal(info, "Sentence Completed")
	} else {
		ef.MaybeEligible(info, "Sentence not Completed")
	}
}

func (ef contraCostaEligibilityFlow) hasProp64ChargeInCycle(info *EligibilityInfo, row *DOJRow) {
	if row.HasProp64ChargeInCycle {
		ef.isBP4060Charge(info, row)
	} else {
		ef.MaybeEligible(info, "No Related Prop64 Charges")
	}
}

func (ef contraCostaEligibilityFlow) isBP4060Charge(info *EligibilityInfo, row *DOJRow) {
	if row.CodeSection == "4060 BP" {
		ef.MaybeEligible(info, "4060 BP")
	} else {
		ef.EligibleDismissal(info, row.CodeSection)
	}
}
