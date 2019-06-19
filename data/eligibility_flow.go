package data

import (
	"time"
)

type EligibilityFlow interface {
	ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo
	BeginEligibilityFlow(info *EligibilityInfo, conviction *DOJRow, subject *Subject)
	ChecksRelatedCharges() (result bool)
}

var EligibilityFlows = map[string]EligibilityFlow{
	"SACRAMENTO":                      sacramentoEligibilityFlow{},
	"SAN JOAQUIN":                     sanJoaquinEligibilityFlow{},
	"CONTRA COSTA":                    contraCostaEligibilityFlow{},
	"LOS ANGELES":                     losAngelesEligibilityFlow{},
	"DISMISS ALL PROP 64":             dismissAllProp64EligibilityFlow{},
	"DISMISS ALL PROP 64 AND RELATED": dismissAllProp64AndRelatedEligibilityFlow{},
}

type EligibilityOptions struct {
	BaselineEligibility BaselineEligibility `json:"baselineEligibility"`
	AdditionalRelief AdditionalRelief `json:"additionalRelief"`
}

type BaselineEligibility struct {
	Dismiss []string `json:"dismiss"`
	Reduce  []string `json:"reduce"`
}

type AdditionalRelief struct {
	Under21 bool `json:"under21"`
}