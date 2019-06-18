package data

import (
	"fmt"
	"gogen/matchers"
	"regexp"
	"strings"
	"time"
)

type configurableEligibilityFlow struct {
	county         string
	dismissMatcher *regexp.Regexp
	reduceMatcher  *regexp.Regexp
}

func NewConfigurableEligibilityFlow(options EligibilityOptions, county string) configurableEligibilityFlow {
	fmt.Printf("Dismiss : %v\n", options.BaselineEligibility.Dismiss)
	fmt.Printf("Reduce : %v\n", options.BaselineEligibility.Reduce)

	var dismissMatcherRegexSource string
	var dismissMatcherRegex *regexp.Regexp
	dismissMatcherRegexSource = strings.Join(escapeRegexMetaChars(options.BaselineEligibility.Dismiss), "|")
	if dismissMatcherRegexSource != "" {
		dismissMatcherRegexSource = ".*(" + dismissMatcherRegexSource + ").*HS"
		dismissMatcherRegex = regexp.MustCompile(dismissMatcherRegexSource)
	}

	var reduceMatcherRegexSource string
	var reduceMatcherRegex *regexp.Regexp
	escapeRegexMetaChars(options.BaselineEligibility.Reduce)
	reduceMatcherRegexSource = strings.Join(escapeRegexMetaChars(options.BaselineEligibility.Reduce), "|")
	if reduceMatcherRegexSource != "" {
		reduceMatcherRegexSource = ".*(" + reduceMatcherRegexSource + ").*HS"
		reduceMatcherRegex = regexp.MustCompile(reduceMatcherRegexSource)
	}

	return configurableEligibilityFlow{
		county:         county,
		dismissMatcher: dismissMatcherRegex,
		reduceMatcher:  reduceMatcherRegex,
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
			ef.BeginEligibilityFlow(info, conviction, subject)
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

func (ef configurableEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if ef.isDismissedCodeSection(row.CodeSection) {
		info.SetEligibleForDismissal(fmt.Sprintf("Dismiss all %s convictions", row.CodeSection))
	} else if ef.isReducedCodeSection(row.CodeSection) {
		info.SetEligibleForReduction(fmt.Sprintf("Reduce all %s convictions", row.CodeSection))
	}
}

func (ef configurableEligibilityFlow) isDismissedCodeSection(codeSection string) bool {
	return ef.dismissMatcher != nil && ef.dismissMatcher.MatchString(codeSection)
}

func (ef configurableEligibilityFlow) isReducedCodeSection(codeSection string) bool {
	return ef.reduceMatcher != nil && ef.reduceMatcher.MatchString(codeSection)
}
