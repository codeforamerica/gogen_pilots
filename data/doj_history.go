package data

import (
	"regexp"
	"strings"
	"time"
)

type DOJHistory struct {
	SubjectID               string
	Name                    string
	WeakName                string
	CII                     string
	DOB                     time.Time
	SSN                     string
	CDL                     string
	Convictions             []*DOJRow
	seenConvictions         map[string]bool
	PC290Registration       bool
	OriginalCII             string
	CyclesWithProp64Charges map[string]bool
	CaseNumbers             map[string][]string
	IsDeceased              bool
}

func (history *DOJHistory) PushRow(row DOJRow, county string) {
	if history.SubjectID == "" {
		history.SubjectID = row.SubjectID
		history.Name = row.Name
		history.WeakName = strings.Split(row.Name, " ")[0]
		history.CII = row.CII
		history.OriginalCII = row.CII
		history.DOB = row.DOB
		history.SSN = row.SSN
		history.CDL = row.CDL
		history.seenConvictions = make(map[string]bool)
		history.CyclesWithProp64Charges = make(map[string]bool)
		history.CaseNumbers = make(map[string][]string)
	}
	if row.Convicted && history.seenConvictions[row.CountOrder] {
		lastConviction := history.Convictions[len(history.Convictions)-1]
		newEndDate := lastConviction.SentenceEndDate.Add(row.SentencePartDuration)
		lastConviction.SentenceEndDate = newEndDate
	}

	if row.Type == "DECEASED" {
		history.IsDeceased = true
	}

	if row.Type == "COURT ACTION" && row.OFN != "" {
		history.CaseNumbers[row.CountOrder[0:6]] = setAppend(history.CaseNumbers[row.CountOrder[0:6]], row.OFN)
	}
	if row.PC290Registration {
		history.PC290Registration = true
	}
	if row.Convicted && !history.seenConvictions[row.CountOrder] {
		row.HasProp64ChargeInCycle = history.CyclesWithProp64Charges[row.CountOrder[0:3]]
		history.Convictions = append(history.Convictions, &row)
		history.seenConvictions[row.CountOrder] = true
	}

	if EligibilityFlows[county].IsProp64Charge(row.CodeSection) {
		history.CyclesWithProp64Charges[row.CountOrder[0:3]] = true
		for _, conviction := range history.Convictions {
			if conviction.CountOrder[0:3] == row.CountOrder[0:3] {
				conviction.HasProp64ChargeInCycle = true
			}
		}
	}
}

func (history *DOJHistory) MostRecentConvictionDate() time.Time {

	var latestDate time.Time

	for _, conviction := range history.Convictions {
		if conviction.DispositionDate.After(latestDate) {
			latestDate = conviction.DispositionDate
		}
	}

	return latestDate
}

func (history *DOJHistory) SuperstrikeCodeSections() []string {
	var result []string
	for _, row := range history.Convictions {
		for _, pattern := range superstrikesPatterns {
			if pattern == row.CodeSection {
				result = append(result, row.CodeSection)
			}
		}
	}
	return result
}

func (history *DOJHistory) PC290CodeSections() []string {
	var result []string

	for _, row := range history.Convictions {
		for _, pattern := range pc290Patterns {
			if pattern.MatchString(row.CodeSection) {
				result = append(result, row.CodeSection)
			}
		}
	}
	return result
}

func (history *DOJHistory) NumberOfProp64Convictions(county string) int {
	result := 0
	for _, row := range history.Convictions {
		if EligibilityFlows[county].IsProp64Charge(row.CodeSection) {
			result++
		}
	}
	return result
}

func (history *DOJHistory) NumberOfConvictionsInCounty(county string) int {
	result := 0
	for _, row := range history.Convictions {
		if row.County == county {
			result++
		}
	}
	return result
}

func (history *DOJHistory) NumberOfFelonies() int {
	felonies := 0
	for _, row := range history.Convictions {
		if row.Felony {
			felonies++
		}
	}
	return felonies
}

func (history *DOJHistory) NumberOfConvictionsInLast7Years() int {
	convictionsInRange := 0

	for _, conviction := range history.Convictions {
		if conviction.OccurredInLast7Years() {
			convictionsInRange++
		}
	}

	return convictionsInRange
}

func setAppend(arr []string, item string) []string {
	for _, el := range arr {
		if el == item {
			return arr
		}
	}
	return append(arr, item)
}

func (history *DOJHistory) computeEligibilities(infos map[int]*EligibilityInfo, comparisonTime time.Time, county string) {
	for _, row := range history.Convictions {
		if row.County == county {
			eligibilityInfo := NewEligibilityInfo(row, history, comparisonTime, county)
			if eligibilityInfo != nil {
				infos[row.Index] = eligibilityInfo
			}
		}
	}
}

var superstrikesPatterns = []string{
	"187 PC",
	"191.5 PC",
	"187-664 PC",
	"191.5-664 PC",
	"209 PC",
	"220 PC",
	"245(D)(3) PC",
	"261(A)(2) PC",
	"261(A)(6) PC",
	"262(A)(2) PC",
	"262(A)(4) PC",
	"264.1 PC",
	"269 PC",
	"286(C)(1) PC",
	"286(C)(2)(A) PC",
	"286(C)(2)(B) PC",
	"286(C)(2)(C) PC",
	"286(C)(3) PC",
	"286(D)(1) PC",
	"286(D)(2) PC",
	"286(D)(3) PC",
	"288(A) PC",
	"288(B)(1) PC",
	"288(B)(2) PC",
	"288A(C)(1) PC",
	"288A(C)(2)(A) PC",
	"288A(C)(2)(B) PC",
	"288A(C)(2)(C) PC",
	"288A(D) PC",
	"288.5(A) PC",
	"289(A)(1)(A) PC",
	"289(A)(1)(B) PC",
	"289(A)(1)(C) PC",
	"289(A)(2)(C) PC",
	"289(J) PC",
	"653F PC",
	"11418(A)(1) PC",
}

var pc290Patterns = []*regexp.Regexp{
	regexp.MustCompile(`290 PC`),//
	regexp.MustCompile(`290(\(|\.|[a-zA-Z])+(.*) PC`),//
	regexp.MustCompile(`236\.1\([BC]\)(.*) PC`),//
	regexp.MustCompile(`243\.4(.*) PC`),//
	regexp.MustCompile(`261 PC`),//
	regexp.MustCompile(`261(\(|\.|[a-zA-Z])+(.*) PC`),//
	regexp.MustCompile(`262\(A\)\(1\) PC`),
	regexp.MustCompile(`264\.1(.*) PC`),//
	regexp.MustCompile(`266 PC`),//
	regexp.MustCompile(`266C PC`),//
	regexp.MustCompile(`266H\(B\)(.*) PC`),
	regexp.MustCompile(`266I\(B\)(.*) PC`),
	regexp.MustCompile(`266J(.*) PC`),//
	regexp.MustCompile(`267 PC`),//
	regexp.MustCompile(`269 PC`),//
	regexp.MustCompile(`269(\(|\.|[a-zA-Z])+(.*) PC`),//
	regexp.MustCompile(`285 PC`),//
	regexp.MustCompile(`286([^\.]*) PC`),
	regexp.MustCompile(`288 PC`),//
	regexp.MustCompile(`288\(([^\.]*) PC`),
	regexp.MustCompile(`288A(.*) PC`),
	regexp.MustCompile(`288\.[23457](.*) PC`),
	regexp.MustCompile(`289([^\.]*) PC`),
	regexp.MustCompile(`311\.1(.*) PC`),//
	regexp.MustCompile(`311\.2\([BCD]\) PC`),
	regexp.MustCompile(`311\.([34]|10|11)(.*) PC`),
	regexp.MustCompile(`314 PC`),//
	regexp.MustCompile(`314(\(|\.|[a-zA-Z])+(.*) PC`),//
	regexp.MustCompile(`647\.6(.*) PC`),//
	regexp.MustCompile(`647A(.*) PC`),//
	regexp.MustCompile(`653F\(C\) PC`),
}
