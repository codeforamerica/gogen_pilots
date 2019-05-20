package data

import (
	"fmt"
	"strings"
	"time"
)

type impactStats struct {
	NoFelonies          bool
	NoConvictions       bool
	NoConvictions7Years bool
}

type DOJHistory struct {
	SubjectID                          string
	Name                               string
	WeakName                           string
	CII                                string
	DOB                                time.Time
	SSN                                string
	CDL                                string
	Convictions                        []*DOJRow
	seenConvictions                    map[string]bool
	PC290Registration                  bool
	OriginalCII                        string
	CyclesWithProp64Charges            map[string]bool
	CaseNumbers                        map[string][]string
	IsDeceased                         bool
	ImpactStatsCountyLogic             impactStats
	ImpactStatsAllDismissed            impactStats
	ImpactStatsAllDismissedWithRelated impactStats
	infos                              map[int]*EligibilityInfo
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

func (history *DOJHistory) numberDismissedCounty() int {
	result := 0
	for _, info := range history.infos {
		if info.EligibilityDetermination["county"] == "Eligible for Dismissal" {
			result++
		}
	}
	return result
}

func (history *DOJHistory) numberFeloniesDismissedCounty() int {
	result := 0
	for i, info := range history.infos {
		if info.EligibilityDetermination["county"] == "Eligible for Dismissal" && history.Convictions[i].Felony {
			result++
		}
	}
	return result
}

func (history *DOJHistory) numberFeloniesReducedCounty() int {
	result := 0
	for i, info := range history.infos {
		if info.EligibilityDetermination["county"] == "Eligible for Reduction" && history.Convictions[i].Felony {
			result++
		}
	}
	return result
}

func (history *DOJHistory) numberFeloniesDismissedAll() int {
	result := 0
	for i, info := range history.infos {
		if info.EligibilityDetermination["allDismissed"] == "Eligible for Dismissal" && history.Convictions[i].Felony {
			result++
		}
	}
	return result
}

func (history *DOJHistory) numberFeloniesDismissedAllWithRelated() int {
	result := 0
	for i, info := range history.infos {
		if info.EligibilityDetermination["allDismissedRelated"] == "Eligible for Dismissal" && history.Convictions[i].Felony {
			result++
		}
	}
	return result
}

func (history *DOJHistory) numberDismissedAll() int {
	result := 0
	for i, info := range history.infos {
		if info.EligibilityDetermination["allDismissed"] == "Eligible for Dismissal" {
			fmt.Printf("Dismissing conviction for %s under new county logic\n", history.Convictions[i].Name)

			result++
		}
	}
	return result
}

func (history *DOJHistory) numberDismissedAllWithRelated() int {
	result := 0
	for _, info := range history.infos {
		if info.EligibilityDetermination["allDismissedRelated"] == "Eligible for Dismissal" {
			result++
		}
	}
	return result
}

func (history *DOJHistory) computeEligibilities(infos map[int]*EligibilityInfo, comparisonTime time.Time, county string) {
	for i, row := range history.Convictions {
		if row.County == county {
			eligibilityInfo := NewEligibilityInfo(row, history, comparisonTime, county)
			if eligibilityInfo != nil {
				infos[row.Index] = eligibilityInfo
				history.infos[i] = eligibilityInfo
			}
		}
	}

	if history.NumberOfFelonies() > 0 {
		history.ImpactStatsCountyLogic.NoFelonies = history.NumberOfFelonies() == (history.numberFeloniesDismissedCounty() + history.numberFeloniesReducedCounty())
		history.ImpactStatsAllDismissed.NoFelonies = history.NumberOfFelonies() == history.numberFeloniesDismissedAll()
		history.ImpactStatsAllDismissedWithRelated.NoFelonies = history.NumberOfFelonies() == history.numberFeloniesDismissedAllWithRelated()
	}

	if len(history.Convictions) > 0 {
		history.ImpactStatsCountyLogic.NoConvictions = len(history.Convictions) == history.numberDismissedCounty()
		history.ImpactStatsAllDismissed.NoConvictions = len(history.Convictions) == history.numberDismissedAll()
		history.ImpactStatsAllDismissedWithRelated.NoConvictions = len(history.Convictions) == history.numberDismissedAllWithRelated()
	}

	if history.NumberOfConvictionsInLast7Years() > 0 {
		history.ImpactStatsCountyLogic.NoConvictions7Years = history.NumberOfConvictionsInLast7Years() == history.numberDismissedCounty()
		history.ImpactStatsAllDismissed.NoConvictions7Years = history.NumberOfConvictionsInLast7Years() == history.numberDismissedAll()
		history.ImpactStatsAllDismissedWithRelated.NoConvictions7Years = history.NumberOfConvictionsInLast7Years() == history.numberDismissedAllWithRelated()
	}
}
