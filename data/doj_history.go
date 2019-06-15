package data

import (
	"time"
)

type DOJHistory struct {
	SubjectID               string
	Name                    string
	DOB                     time.Time
	Convictions             []*DOJRow
	seenConvictions         map[string]bool
	PC290Registration       bool
	CyclesWithProp64Charges map[string]bool
	CaseNumbers             map[string][]string
	IsDeceased              bool
	EligibilityInfos        map[int]*EligibilityInfo
}

func (history *DOJHistory) PushRow(row DOJRow, eligibilityFlow EligibilityFlow) {
	if history.SubjectID == "" {
		history.SubjectID = row.SubjectID
		history.Name = row.Name
		history.DOB = row.DOB
		history.seenConvictions = make(map[string]bool)
		history.CyclesWithProp64Charges = make(map[string]bool)
		history.CaseNumbers = make(map[string][]string)
	}
	if row.WasConvicted && history.seenConvictions[row.CountOrder] {
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
	if row.IsPC290Registration {
		history.PC290Registration = true
	}
	if row.WasConvicted && !history.seenConvictions[row.CountOrder] {
		row.HasProp64ChargeInCycle = history.CyclesWithProp64Charges[row.CountOrder[0:3]]
		history.Convictions = append(history.Convictions, &row)
		history.seenConvictions[row.CountOrder] = true
	}

	if eligibilityFlow.IsProp64Charge(row.CodeSection) {
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
		if IsSuperstrike(row.CodeSection) {
			result = append(result, row.CodeSection)
		}
	}
	return result
}

func (history *DOJHistory) PC290CodeSections() []string {
	var result []string

	for _, row := range history.Convictions {
		if IsPC290(row.CodeSection) {
			result = append(result, row.CodeSection)
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
		if row.IsFelony {
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
