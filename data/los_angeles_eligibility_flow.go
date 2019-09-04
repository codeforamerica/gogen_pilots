package data

import (
	"fmt"
	"gogen_pilots/matchers"
	"strings"
	"time"
)

type losAngelesEligibilityFlow struct {
}

func (ef losAngelesEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string, age float64, yearsConvictionFree int) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, "LOS ANGELES")
			ef.BeginEligibilityFlow(info, conviction, subject, age, yearsConvictionFree)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef losAngelesEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

func (ef losAngelesEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == "LOS ANGELES" && matchers.IsProp64Charge(codeSection)
}

func (ef losAngelesEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if matchers.IsProp64Charge(row.CodeSection) {
		ef.ConvictionIsMisdemeanorOrInfraction(info, row, subject, age, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) ConvictionIsMisdemeanorOrInfraction(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if row.IsFelony {
		ef.ConvictionBeforeNovNine2016(info, row, subject, age, yearsConvictionFree)
	} else {
		info.SetCityAttorneyReview("Misdemeanor or Infraction")
	}
}

func (ef losAngelesEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.ConvictionIs11357(info, row, subject, age, yearsConvictionFree)
	} else {
		info.SetNotEligible("Occurred after 11/09/2016")
	}
}


func (ef losAngelesEligibilityFlow) ConvictionIs11357(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	ok, codeSection := matchers.ExtractProp64Section(row.CodeSection)
	if ok && codeSection == "11357" {
		if strings.HasPrefix(row.CodeSection, "11357(A)") || strings.HasPrefix(row.CodeSection, "11357(B)") {
			info.SetEligibleForDismissal("11357(a) or 11357(b)")
		} else {
			info.SetHandReview("Other 11357")
		}
	} else {
		ef.HasSuperstrikes(info, row, subject, age, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) HasSuperstrikes(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if info.hasSuperstrikes() {
		info.SetNotEligible("PC 667(e)(2)(c)(iv)")
	} else {
		ef.HasPC290(info, row, subject, age, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) HasPC290(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if info.hasPC290() {
		info.SetNotEligible("PC 290")
	} else {
		ef.TwoPriors(info, row, subject, age, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) TwoPriors(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if info.hasTwoPriors(row, subject) {
		info.SetNotEligible("Two priors")
	} else {
		ef.OlderThanGivenAge(info, row, subject, age, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) OlderThanGivenAge(info *EligibilityInfo, row *DOJRow, subject *Subject, age float64, yearsConvictionFree int) {
	if info.olderThanFifty(row, subject, age) {
		info.SetEligibleForDismissal(fmt.Sprintf("%v years or older", age))
	} else {
		ef.YoungerThanTwentyOne(info, row, subject, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) YoungerThanTwentyOne(info *EligibilityInfo, row *DOJRow, subject *Subject, yearsConvictionFree int) {
	if info.youngerThanTwentyOne(row, subject) {
		info.SetEligibleForDismissal("21 years or younger")
	} else {
		ef.Prop64OnlyWithCompletedSentences(info, row, subject, yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) Prop64OnlyWithCompletedSentences(info *EligibilityInfo, row *DOJRow, subject *Subject, yearsConvictionFree int) {
	if info.onlyProp64Convictions(row, subject) && info.allSentencesCompleted(row, subject) {
		info.SetEligibleForDismissal("Only has 11357-60 charges and completed sentence")
	} else {
		ef.NoConvictionsInGivenTimePeriod(info, row, subject,yearsConvictionFree)
	}
}

func (ef losAngelesEligibilityFlow) NoConvictionsInGivenTimePeriod(info *EligibilityInfo, row *DOJRow, subject *Subject, yearsConvictionFree int) {
	if info.noConvictionsPastTenYears(row, subject, yearsConvictionFree) {
		info.SetEligibleForDismissal(fmt.Sprintf("No convictions in past %v years", yearsConvictionFree))
	} else {
		ef.ServingSentence(info, row, subject)
	}
}

func (ef losAngelesEligibilityFlow) ServingSentence(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if !info.allSentencesCompleted(row, subject) {
		info.SetHandReview("Currently serving sentence")
	} else {
		ef.IsDeceased(info, row, subject)
	}
}

func (ef losAngelesEligibilityFlow) IsDeceased(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if subject.IsDeceased {
		info.SetEligibleForDismissal("Deceased")
	} else {
		info.SetHandReview("????")
	}
}
