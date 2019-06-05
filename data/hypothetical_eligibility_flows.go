package data

import (
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
