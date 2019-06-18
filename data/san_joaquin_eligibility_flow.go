package data

import (
	"gogen/matchers"
	"strings"
	"time"
)

type sanJoaquinEligibilityFlow struct {
}

func (ef sanJoaquinEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, "SAN JOAQUIN")
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef sanJoaquinEligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef sanJoaquinEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == "SAN JOAQUIN" && (matchers.IsProp64Charge(codeSection) || matchers.IsRelatedCharge(codeSection))
}

func (ef sanJoaquinEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) || matchers.IsRelatedCharge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.convictionIsRelatedCharge(info, row)
	} else {
		info.SetNotEligible("Occurred after 11/09/2016")
	}
}

func (ef sanJoaquinEligibilityFlow) convictionIsRelatedCharge(info *EligibilityInfo, row *DOJRow) {
	if matchers.IsRelatedCharge(row.CodeSection) {
		ef.hasProp64ChargeInCycle(info, row)
	} else {
		ef.ConvictionIsNotFelony(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) ConvictionIsNotFelony(info *EligibilityInfo, row *DOJRow) {
	if !row.IsFelony {
		info.SetEligibleForDismissal("Misdemeanor or Infraction")
	} else {
		ef.Is11357(info, row)
	}
}

func (ef sanJoaquinEligibilityFlow) Is11357(info *EligibilityInfo, row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357") {
		info.SetEligibleForDismissal("11357 HS")
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
		info.SetEligibleForDismissal("No convictions in past 5 years")
	} else {
		info.SetMaybeEligible("Has convictions in past 5 years")
	}
}

func (ef sanJoaquinEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		info.SetEligibleForDismissal("Sentence Completed")
	} else {
		info.SetMaybeEligible("Sentence not Completed")
	}
}

func (ef sanJoaquinEligibilityFlow) hasProp64ChargeInCycle(info *EligibilityInfo, row *DOJRow) {
	if row.HasProp64ChargeInCycle {
		ef.isBP4060Charge(info, row)
	} else {
		info.SetMaybeEligible("No Related Prop64 Charges")
	}
}

func (ef sanJoaquinEligibilityFlow) isBP4060Charge(info *EligibilityInfo, row *DOJRow) {
	if row.CodeSection == "4060 BP" {
		info.SetMaybeEligible("4060 BP")
	} else {
		info.SetEligibleForDismissal(row.CodeSection)
	}
}
