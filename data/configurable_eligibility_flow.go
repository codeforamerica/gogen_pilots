package data

import (
	"fmt"
	"gogen/matchers"
	"regexp"
	"strings"
	"time"
)

type configurableEligibilityFlow struct {
	county                         string
	dismissMatcher                 *regexp.Regexp
	dismissConvictionsUnderAgeOf21 bool
}

func NewConfigurableEligibilityFlow(options EligibilityOptions, county string) configurableEligibilityFlow {
	var dismissMatcherRegexSource string
	var dismissMatcherRegex *regexp.Regexp
	dismissMatcherRegexSource = strings.Join(escapeRegexMetaChars(options.BaselineEligibility.Dismiss), "|")
	if dismissMatcherRegexSource != "" {
		dismissMatcherRegexSource = ".*(" + dismissMatcherRegexSource + ").*HS"
		dismissMatcherRegex = regexp.MustCompile(dismissMatcherRegexSource)
	}

	return configurableEligibilityFlow{
		county:                         county,
		dismissMatcher:                 dismissMatcherRegex,
		dismissConvictionsUnderAgeOf21: options.AdditionalRelief.Under21,
	}
}

func escapeRegexMetaChars(source []string) []string {
	result := make([]string, len(source))
	for i, s := range source {
		result[i] = regexp.QuoteMeta(s)
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
	info.SetEligibleForReduction(fmt.Sprintf("Reduce all %s convictions", row.CodeSection))
}

func (ef configurableEligibilityFlow) isDismissedCodeSection(codeSection string) bool {
	return ef.dismissMatcher != nil && ef.dismissMatcher.MatchString(codeSection)
}