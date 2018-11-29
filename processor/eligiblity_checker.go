package processor

import (
	"gogen/data"
	"strings"
)

func ComputeEligibility(entry data.CMSEntry, weightInfo data.WeightsEntry) data.EligibilityInfo {
	var eligibleString string

	if strings.HasPrefix(entry.Charge, "11357") || entry.Level == "M" {
		return data.EligibilityInfo{
			QFinalSum: weightInfo.Weight,
			Over1Lb:   "n/a",
		}
	}

	if !weightInfo.Found {
		return data.EligibilityInfo{
			QFinalSum: 0,
			Over1Lb:   "not found",
		}
	}

	if weightInfo.Weight <= 453.592 {
		eligibleString = "eligible"
	} else {
		eligibleString = "ineligible"
	}
	return data.EligibilityInfo{
		QFinalSum: weightInfo.Weight,
		Over1Lb:   eligibleString,
	}
}
