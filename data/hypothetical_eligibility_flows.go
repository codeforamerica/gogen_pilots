package data

import (
	"gogen/matchers"
	"time"
)

type dismissAllProp64EligibilityFlow struct {
}

func (ef dismissAllProp64EligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County, county) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, county)
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef dismissAllProp64EligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef dismissAllProp64EligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) {
		info.SetEligibleForDismissal("Dismiss all Prop 64 charges")
	}
}

func (ef dismissAllProp64EligibilityFlow) checkRelevancy(codeSection string, convictionCounty string, flowCounty string) bool {
	return convictionCounty == flowCounty && matchers.IsProp64Charge(codeSection)
}

type dismissAllProp64AndRelatedEligibilityFlow struct {
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County, county) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, county)
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) ChecksRelatedCharges() bool {
	return true
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	if matchers.IsProp64Charge(row.CodeSection) || matchers.IsRelatedCharge(row.CodeSection) {
		info.SetEligibleForDismissal("Dismiss all Prop 64 and related charges")
	}
}

func (ef dismissAllProp64AndRelatedEligibilityFlow) checkRelevancy(codeSection string, convictionCounty string, flowCounty string) bool {
	return convictionCounty == flowCounty && (matchers.IsProp64Charge(codeSection) || matchers.IsRelatedCharge(codeSection))
}
