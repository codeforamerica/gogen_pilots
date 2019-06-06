package data

import (
	"fmt"
	.	"gogen/matchers"
	"regexp"
	"time"
)

type dismissAllProp64EligibilityFlow struct {
	prop64Matcher *regexp.Regexp
}

func (ef dismissAllProp64EligibilityFlow) ProcessHistory(history *DOJHistory, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range history.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County, county) {
			info := NewEligibilityInfo(conviction, history, comparisonTime, county)
			ef.BeginEligibilityFlow(info, conviction, history)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef dismissAllProp64EligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef dismissAllProp64EligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if ef.IsProp64Charge(row.CodeSection) {
		ef.EligibleDismissal(info, "Dismiss all Prop 64 charges")
	}
}

func (ef dismissAllProp64EligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = reason
}

func (ef dismissAllProp64EligibilityFlow) checkRelevancy(codeSection string, convictionCounty string, flowCounty string) bool {
	return convictionCounty == flowCounty && ef.IsProp64Charge(codeSection)
}

func (ef dismissAllProp64EligibilityFlow) IsProp64Charge(codeSection string) bool {
	ok, _ := Prop64Matcher(codeSection)
	return ok
}

func (ef dismissAllProp64EligibilityFlow) MatchedCodeSection(codeSection string) string {
	matches := ef.prop64Matcher.FindStringSubmatch(codeSection)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func (ef dismissAllProp64EligibilityFlow) MatchedRelatedCodeSection(codeSection string) string {
	return ""
}

type dismissAllProp64AndRelatedEligibilityFlow struct {
	prop64Matcher *regexp.Regexp
	relatedChargeMatcher *regexp.Regexp
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) ProcessHistory(history *DOJHistory, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range history.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County, county) {
			info := NewEligibilityInfo(conviction, history, comparisonTime, county)
			ef.BeginEligibilityFlow(info, conviction, history)
			infos[conviction.Index] = info
			fmt.Printf("\n\ndetermination for code section %v for is: %v \n\n", conviction.CodeSection, infos[conviction.Index].EligibilityDetermination)
		}
	}
	return infos
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, history *DOJHistory) {
	if ef.IsProp64Charge(row.CodeSection) || ef.IsRelatedCharge(row.CodeSection){
		ef.EligibleDismissal(info, "Dismiss all Prop 64 and related charges")
	}
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Eligible for Dismissal"
	info.EligibilityReason = reason
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) checkRelevancy(codeSection string, convictionCounty string, flowCounty string) bool {
	return convictionCounty == flowCounty && ef.IsProp64Charge(codeSection)
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) IsProp64Charge(codeSection string) bool {
	ok, _ := Prop64Matcher(codeSection)
	return ok
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) MatchedCodeSection(codeSection string) string {
	matches := ef.prop64Matcher.FindStringSubmatch(codeSection)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) MatchedRelatedCodeSection(codeSection string) string {
	relatedChargeMatches := ef.relatedChargeMatcher.FindStringSubmatch(codeSection)
	if len(relatedChargeMatches) > 0 {
		return relatedChargeMatches[1]
	}

	return ""
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) IsRelatedCharge(codeSection string) bool {
	ok, _ := RelatedChargeMatcher(codeSection)
	return ok
}

