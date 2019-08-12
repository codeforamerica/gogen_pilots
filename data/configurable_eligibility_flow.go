package data

import (
	"errors"
	"fmt"
	"gogen/matchers"
	"time"
)

type ConfigurableEligibilityFlow struct {
	county                               string
	DismissSections                      []string
	ReduceSections                       []string
	dismissConvictionsUnderAgeOf21       bool
	dismissIfSubjectHasOnlyProp64Charges bool
	dismissIfSubjectIsDeceased           bool
	subjectAgeThreshold                  int
	yearsSinceConvictionThreshold        int
	yearsCrimeFreeThreshold              int
}

func NewConfigurableEligibilityFlow(options EligibilityOptions, county string) (ConfigurableEligibilityFlow, error) {

	if options.AdditionalRelief.SubjectAgeThreshold != 0 {
		if options.AdditionalRelief.SubjectAgeThreshold > 65 || options.AdditionalRelief.SubjectAgeThreshold < 40 {
			return ConfigurableEligibilityFlow{}, errors.New("SubjectAgeThreshold should be between 40 and 65, or 0")
		}
	}

	if options.AdditionalRelief.YearsSinceConvictionThreshold != 0 {
		if options.AdditionalRelief.YearsSinceConvictionThreshold > 15 || options.AdditionalRelief.YearsSinceConvictionThreshold < 1 {
			return ConfigurableEligibilityFlow{}, errors.New("YearsSinceConvictionThreshold should be between 1 and 15, or 0")
		}
	}

	if options.AdditionalRelief.YearsCrimeFreeThreshold != 0 {
		if options.AdditionalRelief.YearsCrimeFreeThreshold > 15 || options.AdditionalRelief.YearsCrimeFreeThreshold < 1 {
			return ConfigurableEligibilityFlow{}, errors.New("YearsCrimeFreeThreshold should be between 1 and 15, or 0")
		}
	}

	return ConfigurableEligibilityFlow{
		county:                               county,
		DismissSections:                      options.BaselineEligibility.Dismiss,
		ReduceSections:                       options.BaselineEligibility.Reduce,
		dismissConvictionsUnderAgeOf21:       options.AdditionalRelief.SubjectUnder21AtConviction,
		dismissIfSubjectIsDeceased:           options.AdditionalRelief.SubjectIsDeceased,
		dismissIfSubjectHasOnlyProp64Charges: options.AdditionalRelief.SubjectHasOnlyProp64Charges,
		subjectAgeThreshold:                  options.AdditionalRelief.SubjectAgeThreshold,
		yearsSinceConvictionThreshold:        options.AdditionalRelief.YearsSinceConvictionThreshold,
		yearsCrimeFreeThreshold:              options.AdditionalRelief.YearsCrimeFreeThreshold,
	},
	nil
}

func (ef ConfigurableEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, ef.county)
			ef.EvaluateEligibility(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef ConfigurableEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

func (ef ConfigurableEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == ef.county && matchers.IsProp64Charge(codeSection)
}

func (ef ConfigurableEligibilityFlow) EvaluateEligibility(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if !row.IsFelony {
		info.SetEligibleForDismissal("Misdemeanor or Infraction")
		return
	}
	if matched, canonicalCodeSection := ef.isDismissedCodeSection(row.CodeSection); matched {
		info.SetEligibleForDismissal(composeEligibilityReason(canonicalCodeSection, true))
		return
	}
	if ef.dismissConvictionsUnderAgeOf21 && row.wasConvictionUnderAgeOf21(subject) {
		info.SetEligibleForDismissal("21 years or younger")
		return
	}
	if ef.subjectAgeThreshold != 0 && subject.olderThan(ef.subjectAgeThreshold, info.comparisonTime) {
		info.SetEligibleForDismissal(fmt.Sprintf("%d years or older", ef.subjectAgeThreshold))
		return
	}
	if ef.yearsSinceConvictionThreshold != 0 && row.convictionBefore(ef.yearsSinceConvictionThreshold, info.comparisonTime) {
		info.SetEligibleForDismissal(fmt.Sprintf("Conviction occurred %d or more years ago", ef.yearsSinceConvictionThreshold))
		return
	}
	if ef.yearsCrimeFreeThreshold != 0 && subject.MostRecentConvictionDate().Before(info.comparisonTime.AddDate(-ef.yearsCrimeFreeThreshold, 0, 0)) {
		info.SetEligibleForDismissal(fmt.Sprintf("No convictions in the past %d years", ef.yearsCrimeFreeThreshold))
		return
	}
	if ef.dismissIfSubjectHasOnlyProp64Charges && info.onlyProp64Convictions(row, subject) {
		info.SetEligibleForDismissal("Only has 11357-60 charges")
		return
	}

	if ef.dismissIfSubjectIsDeceased && subject.IsDeceased {
		info.SetEligibleForDismissal("Individual is deceased")
		return
	}

	if matched, canonicalCodeSection := ef.isReducedCodeSection(row.CodeSection); matched {
		info.SetEligibleForReduction(composeEligibilityReason(canonicalCodeSection, false))
		return
	}
}

func (ef ConfigurableEligibilityFlow) isDismissedCodeSection(candidateCodeSection string) (bool, string) {
	for _, codeSection := range ef.DismissSections {
		if matchers.Prop64MatchersByCodeSection[codeSection].MatchString(candidateCodeSection) {
			return true, codeSection
		}
	}
	return false, ""
}

func (ef ConfigurableEligibilityFlow) isReducedCodeSection(candidateCodeSection string) (bool, string) {
	for _, codeSection := range ef.ReduceSections {
		if matchers.Prop64MatchersByCodeSection[codeSection].MatchString(candidateCodeSection) {
			return true, codeSection
		}
	}
	return false, ""
}

func composeEligibilityReason(canonicalCodeSection string, isDismiss bool) string {
	var verb string
	if isDismiss {
		verb = "Dismiss"
	} else {
		verb = "Reduce"
	}
	return fmt.Sprintf("%s all HS %s convictions", verb, canonicalCodeSection)
}
