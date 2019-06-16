package data

import (
	"gogen/matchers"
	"regexp"
	"strings"
	"time"
)

type sacramentoEligibilityFlow struct {
	prop64Matcher *regexp.Regexp
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
	return county == "SACRAMENTO" && ef.IsProp64Charge(codeSection)
}

func (ef sacramentoEligibilityFlow) IsProp64Charge(codeSection string) bool {
	ok, _ := matchers.Prop64Matcher(codeSection)
	return ok
}

func (ef sacramentoEligibilityFlow) MatchedCodeSection(codeSection string) string {
	matches := ef.prop64Matcher.FindStringSubmatch(codeSection)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func (ef sacramentoEligibilityFlow) MatchedRelatedCodeSection(codeSection string) string {
	return ""
}

func (ef sacramentoEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if ef.IsProp64Charge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row)
	}
}

func (ef sacramentoEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = reason
}

func (ef sacramentoEligibilityFlow) EligibleReduction(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Reduction"
	info.EligibilityReason = reason
}

func (ef sacramentoEligibilityFlow) NotEligible(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Not eligible"
	info.EligibilityReason = reason
}

func (ef sacramentoEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.ConvictionIsNotFelony(info, row)
	} else {
		ef.NotEligible(info, "Occurred after 11/09/2016")
	}
}

func (ef sacramentoEligibilityFlow) ConvictionIsNotFelony(info *EligibilityInfo, row *DOJRow) {
	if !row.IsFelony {
		ef.EligibleDismissal(info, "Misdemeanor or Infraction")
	} else {
		ef.Is11357b(info, row)
	}
}

func (ef sacramentoEligibilityFlow) Is11357b(info *EligibilityInfo, row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357(b)") {
		ef.EligibleDismissal(info, "HS 11357(b)")
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
		ef.EligibleDismissal(info, "No convictions in past 10 years")
	} else {
		ef.EligibleReduction(info, "Has convictions in past 10 years")
	}
}

func (ef sacramentoEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		ef.EligibleDismissal(info, "Sentence Completed")
	} else {
		ef.EligibleReduction(info, "Sentence not Completed")
	}
}
