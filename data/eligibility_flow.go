package data

import (
	"regexp"
	"time"
)

type EligibilityFlow interface {
	ProcessHistory(history *DOJHistory, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo
	BeginEligibilityFlow(info *EligibilityInfo, conviction *DOJRow, history *DOJHistory)
	ChecksRelatedCharges() (result bool)
	IsProp64Charge(codeSection string) (result bool)
	MatchedCodeSection(codeSection string) (matchedCodeSection string)
	MatchedRelatedCodeSection(codeSection string) (matchedCodeSection string)

}

var EligibilityFlows = map[string]EligibilityFlow{
	"SACRAMENTO": sacramentoEligibilityFlow{
		prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360)`),
	},
	"SAN JOAQUIN": sanJoaquinEligibilityFlow{
		prop64Matcher:        regexp.MustCompile(`(11357|11358|11359|11360)`),
		relatedChargeMatcher: regexp.MustCompile(`(647\(f\)\s*PC|602\s*PC|466\s*PC|148\.9\s*PC|148\s*PC|11364\s*HS|11550\s*HS|4140\s*BP|4149\s*BP|4060\s*BP|40508\s*VC|1320[^\d\.][^\.]*PC).*`),
	},
	"CONTRA COSTA": contraCostaEligibilityFlow{
		prop64Matcher:        regexp.MustCompile(`(11357|11358|11359|11360)`),
		relatedChargeMatcher: regexp.MustCompile(`(647\(f\)\s*PC|602\s*PC|466\s*PC|148\.9\s*PC|148\s*PC|11364\s*HS|11550\s*HS|4140\s*BP|4149\s*BP|4060\s*BP|40508\s*VC|1320[^\d\.][^\.]*PC).*`),
	},
	"LOS ANGELES": losAngelesEligibilityFlow{
		prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360)`),
	},
	"DISMISS ALL PROP 64": dismissAllProp64EligibilityFlow{
		prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360)`),
	},
	"DISMISS ALL PROP 64 AND RELATED": dismissAllProp64AndRelatedEligibilityFlow{
		prop64Matcher:        regexp.MustCompile(`(11357|11358|11359|11360)`),
		relatedChargeMatcher: regexp.MustCompile(`(647\(f\)\s*PC|602\s*PC|466\s*PC|148\.9\s*PC|148\s*PC|11364\s*HS|11550\s*HS|4140\s*BP|4149\s*BP|4060\s*BP|40508\s*VC|1320[^\d\.][^\.]*PC).*`),
	},
}
