package processor

import (
	"fmt"
	"gogen/data"
	"strings"
)

type EligibilityInfo struct {
	QFinalSum string
	Over1Lb string
	PC290Registration string
	PC290Charges string
	PC290CodeSections string
	Superstrikes string
	SuperstrikeCodeSections string
	TwoPriors string
	AgeAtConviction string
	YearsSinceEvent string
	YearsSinceMostRecentConviction string
	FinalRecommendation string
}

func checkWeight(entry data.CMSEntry, weightInfo data.WeightsEntry, info *EligibilityInfo) {
	var eligibleString string

	if strings.HasPrefix(entry.Charge, "11357") || entry.Level == "M" {
		info.QFinalSum = "n/a"
		info.Over1Lb = "n/a"
		return
	}

	if !weightInfo.Found {
		info.QFinalSum = "not found"
		info.Over1Lb = "not found"
		return
	}

	if weightInfo.Weight <= 453.592 {
		eligibleString = "eligible"
	} else {
		eligibleString = "ineligible"
	}
	info.QFinalSum = fmt.Sprintf("%.1f", weightInfo.Weight)
	info.Over1Lb = eligibleString
}

func checkDOJHistory(entry data.CMSEntry, history *data.DOJHistory, info *EligibilityInfo) {
	result := ""
	if(history == nil) {
		result = "no match"

		info.PC290Registration = result
		info.PC290Charges = result
		info.PC290CodeSections = result
		return
	}

	if entry.Level != "F" || strings.HasPrefix(entry.Charge, "11357") {
		result = "n/a"

		info.PC290Registration = result
		info.PC290Charges = result
		info.PC290CodeSections = result
		return
	}



	//def check_pc290_registration(level, charge, doj_history)
	//return 'n/a' unless level == 'F'
	//return 'n/a' if charge.start_with?('11357')
	//return 'no match' if doj_history.nil?
	//return 'ineligible' if doj_history.pc290_registration?
	//'eligible'
  //end
  //
	//def check_pc290_charges(level, charge, doj_history)
	//return 'n/a' unless level == 'F'
	//return 'n/a' if charge.start_with?('11357')
	//return 'no match' if doj_history.nil?
	//return 'ineligible' if doj_history.pc290_charges?
	//'eligible'
  //end
  //
	//def check_pc290_code_sections(level, charge, doj_history)
	//return 'n/a' unless level == 'F'
	//return 'n/a' if charge.start_with?('11357')
	//return 'no match' if doj_history.nil?
	//return '-' if doj_history.pc290_code_sections.length == 0
	//doj_history.pc290_code_sections.map { |count| count.code_section }.join("; ")
	//end
  //
	//def check_superstrikes(level, charge, doj_history)
	//return 'n/a' unless level == 'F'
	//return 'n/a' if charge.start_with?('11357')
	//return 'no match' if doj_history.nil?
	//return 'ineligible' if doj_history.superstrikes?
	//'eligible'
  //end
  //
	//def check_superstrike_code_sections(level, charge, doj_history)
	//return 'n/a' unless level == 'F'
	//return 'n/a' if charge.start_with?('11357')
	//return 'no match' if doj_history.nil?
	//return '-' if doj_history.superstrike_code_sections.length == 0
	//doj_history.superstrike_code_sections.map { |count| count.code_section }.join("; ")
	//end
  //
	//def check_two_priors(level, charge, doj_history)
	//return 'n/a' unless level == 'F'
	//return 'n/a' if charge.start_with?('11357')
	//return 'no match' if doj_history.nil?
	//return 'ineligible' if doj_history.three_convictions_same_code?(charge)
	//'eligible'
  //end
  //
	//def check_age_at_conviction(dob, disposition_date)
	//return 'not found' if dob.nil? || disposition_date.nil?
	//days = (disposition_date - dob).to_i
	//(days / 365.25).to_i
	//end
  //
	//def check_years_since_event(disposition_date)
	//return 'not found' if disposition_date.nil?
	//days = Date.today - disposition_date
	//(days / 365.25).round(1)
	//end
  //
	//def check_years_since_recent_conviction(doj_history)
	//return 'no match' if doj_history.nil?
	//disposition_date = doj_history.most_recent_conviction_date
	//check_years_since_event(disposition_date)
	//end
  //

}

func computeFinalEligibility(info *EligibilityInfo) {
	//def final_recommendation(eligibility_checks)
	//return 'ineligible' if eligibility_checks.include?('ineligible')
	//return 'needs review' if eligibility_checks.include?('no match')
	//'eligible'
	//end
}

func ComputeEligibility(entry data.CMSEntry, weightInfo data.WeightsEntry, history *data.DOJHistory) *EligibilityInfo {
	eligibilityInfo := new(EligibilityInfo)
	checkWeight(entry, weightInfo, eligibilityInfo)
	checkDOJHistory(entry, history, eligibilityInfo)
	computeFinalEligibility(eligibilityInfo)
	return eligibilityInfo
}

