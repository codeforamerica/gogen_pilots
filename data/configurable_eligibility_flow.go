package data

import (
	"fmt"
	"gogen/matchers"
	"regexp"
	"time"
)

type configurableEligibilityFlow struct {
	county                               string
	dismissMatcher                       []*regexp.Regexp
	dismissConvictionsUnderAgeOf21       bool
	dismissIfSubjectHasOnlyProp64Charges bool
	subjectAgeThreshold                  int
	yearsSinceConvictionThreshold        int
}

func NewConfigurableEligibilityFlow(options EligibilityOptions, county string) configurableEligibilityFlow {
	dismissMatcherRegex := makeRegexes(options.BaselineEligibility.Dismiss)

	if options.AdditionalRelief.SubjectAgeThreshold != 0 {
		if options.AdditionalRelief.SubjectAgeThreshold > 65 || options.AdditionalRelief.SubjectAgeThreshold < 40 {
			panic("SubjectAgeThreshold should be between 40 and 65, or 0")
		}
	}

	if options.AdditionalRelief.YearsSinceConvictionThreshold != 0 {
		if options.AdditionalRelief.YearsSinceConvictionThreshold > 15 || options.AdditionalRelief.YearsSinceConvictionThreshold < 1 {
			panic("YearsSinceConvictionThreshold should be between 1 and 15, or 0")
		}
	}

	return configurableEligibilityFlow{
		county:                               county,
		dismissMatcher:                       dismissMatcherRegex,
		dismissConvictionsUnderAgeOf21:       options.AdditionalRelief.SubjectUnder21AtConviction,
		dismissIfSubjectHasOnlyProp64Charges: options.AdditionalRelief.SubjectHasOnlyProp64Charges,
		subjectAgeThreshold:                  options.AdditionalRelief.SubjectAgeThreshold,
		yearsSinceConvictionThreshold:        options.AdditionalRelief.YearsSinceConvictionThreshold,
	}
}

func makeRegexes(source []string) []*regexp.Regexp {
	result := make([]*regexp.Regexp, len(source))
	for i, s := range source {
		result[i] = regexp.MustCompile(regexp.QuoteMeta(s) + ".*HS")
	}
	return result
}

func (ef configurableEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo {
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

func (ef configurableEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

func (ef configurableEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return county == ef.county && matchers.IsProp64Charge(codeSection)
}

func (ef configurableEligibilityFlow) EvaluateEligibility(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if !row.IsFelony {
		info.SetEligibleForDismissal("Misdemeanor or Infraction")
		return
	}
	if ef.isDismissedCodeSection(row.CodeSection) {
		info.SetEligibleForDismissal(fmt.Sprintf("Dismiss all %s convictions", row.CodeSection))
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
	if ef.dismissIfSubjectHasOnlyProp64Charges && info.onlyProp64Convictions(row, subject) {
		info.SetEligibleForDismissal("Only has 11357-60 charges")
		return
	}

	info.SetEligibleForReduction(fmt.Sprintf("Reduce all %s convictions", row.CodeSection))
}

func (ef configurableEligibilityFlow) isDismissedCodeSection(codeSection string) bool {
	if len(ef.dismissMatcher) > 0 {
		for _, regex := range ef.dismissMatcher {
			if regex.MatchString(codeSection) {
				return true
			}
		}
	}
	return false
}
