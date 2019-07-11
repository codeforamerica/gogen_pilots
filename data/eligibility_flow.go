package data

import (
	"time"
)

type EligibilityFlow interface {
	ProcessSubject(subject *Subject, comparisonTime time.Time, flowCounty string) map[int]*EligibilityInfo
	ChecksRelatedCharges() (result bool)
}

var EligibilityFlows = map[string]EligibilityFlow{
	"LOS ANGELES":                     losAngelesEligibilityFlow{},
	"DISMISS ALL PROP 64":             dismissAllProp64EligibilityFlow{},
	"DISMISS ALL PROP 64 AND RELATED": dismissAllProp64AndRelatedEligibilityFlow{},
}

type EligibilityOptions struct {
	BaselineEligibility BaselineEligibility `json:"baselineEligibility"`
	AdditionalRelief    AdditionalRelief    `json:"additionalRelief"`
}

type BaselineEligibility struct {
	Dismiss []string `json:"dismiss"`
	Reduce  []string `json:"reduce"`
}

type AdditionalRelief struct {
	SubjectUnder21AtConviction    bool `json:"subjectUnder21AtConviction"`
	SubjectHasOnlyProp64Charges   bool `json:"subjectHasOnlyProp64Charges"`
	SubjectAgeThreshold           int  `json:"subjectAgeThreshold"`
	YearsSinceConvictionThreshold int  `json:"yearsSinceConvictionThreshold"`
}
