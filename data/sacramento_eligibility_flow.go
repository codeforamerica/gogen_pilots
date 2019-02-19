package data

import (
	"strings"
	"time"
)

type sacramentoEligibilityFlow struct{}

func (ef sacramentoEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow) {
	ef.ConvictionBeforeNovNine2016(info, row)
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
	if !row.Felony {
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
		ef.ThisConvictionOlderThan10Years(info, row)
	} else {
		ef.CurrentlyServingSentence(info, row)
	}
}

func (ef sacramentoEligibilityFlow) ThisConvictionOlderThan10Years(info *EligibilityInfo, row *DOJRow) {
	if info.YearsSinceThisConviction > 10 {
		ef.FinalConvictionOnRecord(info, row)
	} else {
		ef.EligibleReduction(info, "Occurred in last 10 years")
	}
}

func (ef sacramentoEligibilityFlow) CurrentlyServingSentence(info *EligibilityInfo, row *DOJRow) {
	if row.SentenceEndDate.Before(info.comparisonTime) {
		ef.EligibleDismissal(info, "Sentence Completed")
	} else {
		ef.EligibleReduction(info, "Sentence not Completed")
	}
}

func (ef sacramentoEligibilityFlow) FinalConvictionOnRecord(info *EligibilityInfo, row *DOJRow) {
	if info.YearsSinceMostRecentConviction == info.YearsSinceThisConviction {
		ef.EligibleDismissal(info, "Final Conviction older than 10 years")
	} else {
		ef.EligibleReduction(info, "Later Convictions")
	}
}