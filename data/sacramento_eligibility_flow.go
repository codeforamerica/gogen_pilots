package data

import (
	"gogen/matchers"
	"strings"
	"time"
)

type sacramentoEligibilityFlow struct {
}

func (ef sacramentoEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, "SACRAMENTO")
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef sacramentoEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

func (ef sacramentoEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == "SACRAMENTO" && matchers.IsProp64Charge(codeSection)
}

func (ef sacramentoEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row)
	}
}

func (ef sacramentoEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.ConvictionIsNotFelony(info, row)
	} else {
		info.SetNotEligible("Occurred after 11/09/2016")
	}
}

func (ef sacramentoEligibilityFlow) ConvictionIsNotFelony(info *EligibilityInfo, row *DOJRow) {
	if !row.IsFelony {
		info.SetEligibleForDismissal("Misdemeanor or Infraction")
	} else {
		ef.Is11357b(info, row)
	}
}

func (ef sacramentoEligibilityFlow) Is11357b(info *EligibilityInfo, row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357(b)") {
		info.SetEligibleForDismissal("HS 11357(b)")
	} else {
		ef.MoreThanOneConviction(info, row)
	}
}

func (ef sacramentoEligibilityFlow) MoreThanOneConviction(info *EligibilityInfo, row *DOJRow) {
	if info.NumberOfConvictionsOnRecord > 1 {
		ef.HasConvictionsInPast10Years(info, row)
	} else {
		ef.CurrentlyServingSentence(info, row)
	}
}

func (ef sacramentoEligibilityFlow) HasConvictionsInPast10Years(info *EligibilityInfo, row *DOJRow) {
	if info.YearsSinceMostRecentConviction > 10 {
		info.SetEligibleForDismissal("No convictions in past 10 years")
	} else {
		info.SetEligibleForReduction("Has convictions in past 10 years")
	}
}

func (ef sacramentoEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		info.SetEligibleForDismissal("Sentence Completed")
	} else {
		info.SetEligibleForReduction("Sentence not Completed")
	}
}
