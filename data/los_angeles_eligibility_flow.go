package data

import (
	"regexp"
	"strings"
	"time"
)

type losAngelesEligibilityFlow struct {
	prop64Matcher *regexp.Regexp
}

func (ef losAngelesEligibilityFlow) IsProp64Charge(codeSection string) bool {
	return ef.prop64Matcher.Match([]byte(codeSection))
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

func (ef losAngelesEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow) {
	if ef.IsProp64Charge(row.CodeSection) {
		ef.ConvictionBeforeNovNine2016(info, row)
	}
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

func (ef losAngelesEligibilityFlow) ConvictionBeforeNovNine2016(info *EligibilityInfo, row *DOJRow) {
	if info.DateOfConviction.Before(time.Date(2016, 11, 9, 0, 0, 0, 0, time.UTC)) {
		ef.ConvictionIs11357(info, row)
	} else {
		ef.NotEligible(info, "Occurred after 11/09/2016")
	}
}

func (ef losAngelesEligibilityFlow) ConvictionIsNotFelony(info *EligibilityInfo, row *DOJRow) {
	if !row.Felony {
		ef.EligibleDismissal(info, "Misdemeanor or Infraction")
	} else {
		ef.Is11357b(info, row)
	}
}

func (ef losAngelesEligibilityFlow) Is11357b(info *EligibilityInfo, row *DOJRow) {
	if strings.HasPrefix(row.CodeSection, "11357(b)") {
		ef.EligibleDismissal(info, "HS 11357(b)")
	} else {
		ef.MoreThanOneConviction(info, row)
	}
}

func (ef losAngelesEligibilityFlow) MoreThanOneConviction(info *EligibilityInfo, row *DOJRow) {
	if info.NumberOfConvictionsOnRecord > 1 {
		ef.HasConvictionsInPast10Years(info, row)
	} else {
		ef.CurrentlyServingSentence(info, row)
	}
}

func (ef losAngelesEligibilityFlow) HasConvictionsInPast10Years(info *EligibilityInfo, row *DOJRow) {
	if info.YearsSinceMostRecentConviction > 10 {
		ef.EligibleDismissal(info, "No convictions in past 10 years")
	} else {
		ef.EligibleReduction(info, "Has convictions in past 10 years")
	}
}

func (ef losAngelesEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		ef.EligibleDismissal(info, "Sentence Completed")
	} else {
		ef.EligibleReduction(info, "Sentence not Completed")
	}
}

func (ef losAngelesEligibilityFlow) ConvictionIs11357(info *EligibilityInfo, row *DOJRow) {
	if ef.MatchedCodeSection(row.CodeSection) == "11357" {
		if strings.HasPrefix(row.CodeSection, "11357(A)") || strings.HasPrefix(row.CodeSection, "11357(B)") {
			ef.EligibleDismissal(info,"11357(a) or 11357(b)")
		} else {
			ef.MaybeEligible(info, "Other 11357")
		}
	} else {
		ef.HasSuperstrikes(info, row)
	}
}

func (ef losAngelesEligibilityFlow) HasSuperstrikes(info *EligibilityInfo, row *DOJRow) {
	if info.hasSuperstrikes() {
		ef.NotEligible(info, "PC 667(e)(2)(c)(iv)")
	} else {
		ef.ConvictionIsNotFelony(info, row)
	}
}
