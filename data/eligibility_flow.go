package data

import "regexp"

type EligibilityFlow interface {
	BeginEligibilityFlow(info *EligibilityInfo, conviction *DOJRow)
	IsProp64Charge(codeSection string) (result bool)
}

var eligibilityFlows = map[string]EligibilityFlow{
	"SACRAMENTO":  sacramentoEligibilityFlow{prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360).*`)},
	"SAN JOAQUIN": sanJoaquinEligibilityFlow{prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360).*`)},
}
