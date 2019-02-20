package data

import "regexp"

type EligibilityFlow interface {
	BeginEligibilityFlow(info *EligibilityInfo, conviction *DOJRow)
	IsProp64Charge(codeSection string) (result bool)
}

var eligibilityFlows = map[string]EligibilityFlow{
	"SACRAMENTO": sacramentoEligibilityFlow{
		prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360).*`),
	},
	"SAN JOAQUIN": sanJoaquinEligibilityFlow{
		prop64Matcher:        regexp.MustCompile(`(11357|11358|11359|11360).*`),
		relatedChargeMatcher: regexp.MustCompile(`(647(f)\s*PC|602\s*PC|466\s*PC|148.9\s*PC|148\s*PC|11364\s*HS|11550\s*HS|4140\s*BP|4149\s*BP|4060\s*BP|40508\s*VC).*`),
	},
}
