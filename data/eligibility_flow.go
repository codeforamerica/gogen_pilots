package data

type EligibilityFlow interface {
	BeginEligibilityFlow(info *EligibilityInfo, conviction *DOJRow)
}

var eligibilityFlows = map[string]EligibilityFlow{
	"SACRAMENTO":  sacramentoEligibilityFlow{},
	"SAN JOAQUIN": sanJoaquinEligibilityFlow{},
}
