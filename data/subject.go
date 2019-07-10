package data

import (
	"gogen/matchers"
	"time"
)

type Subject struct {
	ID                      string
	Name                    string
	DOB                     time.Time
	Convictions             []*DOJRow
	seenConvictions         map[string]bool
	PC290Registration       bool
	CyclesWithProp64Charges map[string]bool
	CaseNumbers             map[string][]string
	IsDeceased              bool
}

func (subject *Subject) PushRow(row DOJRow, eligibilityFlow EligibilityFlow) {
	if subject.ID == "" {
		subject.ID = row.SubjectID
		subject.Name = row.Name
		subject.DOB = row.DOB
		subject.seenConvictions = make(map[string]bool)
		subject.CyclesWithProp64Charges = make(map[string]bool)
		subject.CaseNumbers = make(map[string][]string)
	}
	if row.WasConvicted && subject.seenConvictions[row.CountOrder] {
		lastConviction := subject.Convictions[len(subject.Convictions)-1]
		newEndDate := lastConviction.SentenceEndDate.Add(row.SentencePartDuration)
		lastConviction.SentenceEndDate = newEndDate
	}

	if row.Type == "DECEASED" {
		subject.IsDeceased = true
	}

	if row.Type == "COURT ACTION" && row.OFN != "" {
		subject.CaseNumbers[row.CountOrder[0:6]] = setAppend(subject.CaseNumbers[row.CountOrder[0:6]], row.OFN)
	}
	if row.IsPC290Registration {
		subject.PC290Registration = true
	}
	if row.WasConvicted && !subject.seenConvictions[row.CountOrder] {
		row.HasProp64ChargeInCycle = subject.CyclesWithProp64Charges[row.CountOrder[0:3]]
		subject.Convictions = append(subject.Convictions, &row)
		subject.seenConvictions[row.CountOrder] = true
	}

	if matchers.IsProp64Charge(row.CodeSection) {
		subject.CyclesWithProp64Charges[row.CountOrder[0:3]] = true
		for _, conviction := range subject.Convictions {
			if conviction.CountOrder[0:3] == row.CountOrder[0:3] {
				conviction.HasProp64ChargeInCycle = true
			}
		}
	}
}

func (subject *Subject) MostRecentConvictionDate() time.Time {

	var latestDate time.Time

	for _, conviction := range subject.Convictions {
		if conviction.DispositionDate.After(latestDate) {
			latestDate = conviction.DispositionDate
		}
	}

	return latestDate
}

func (subject *Subject) SuperstrikeCodeSections() []string {
	var result []string
	for _, row := range subject.Convictions {
		if IsSuperstrike(row.CodeSection) {
			result = append(result, row.CodeSection)
		}
	}
	return result
}

func (subject *Subject) PC290CodeSections() []string {
	var result []string

	for _, row := range subject.Convictions {
		if IsPC290(row.CodeSection) {
			result = append(result, row.CodeSection)
		}
	}
	return result
}

func (subject *Subject) Prop64ConvictionsBySection() (int, int, int, int, int) {
	convictionCountBySection := make(map[string]int)

	for _, conviction := range subject.Convictions {
		matched, codeSection := matchers.ExtractProp64Section(conviction.CodeSection)
		if matched {
			convictionCountBySection[codeSection]++
		}
	}

	return convictionCountBySection["11357"] + convictionCountBySection["11358"] + convictionCountBySection["11359"] + convictionCountBySection["11360"],
		convictionCountBySection["11357"],
		convictionCountBySection["11358"],
		convictionCountBySection["11359"],
		convictionCountBySection["11360"]
}

func (subject *Subject) NumberOfConvictionsInCounty(county string) int {
	result := 0
	for _, row := range subject.Convictions {
		if row.County == county {
			result++
		}
	}
	return result
}

func (subject *Subject) NumberOfFelonies() int {
	felonies := 0
	for _, row := range subject.Convictions {
		if row.IsFelony {
			felonies++
		}
	}
	return felonies
}

func (subject *Subject) NumberOfConvictionsInLast7Years() int {
	convictionsInRange := 0

	for _, conviction := range subject.Convictions {
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

func (subject *Subject) computeEligibilities(infos map[int]*EligibilityInfo, comparisonTime time.Time, county string) {
	for _, row := range subject.Convictions {
		if row.County == county {
			eligibilityInfo := NewEligibilityInfo(row, subject, comparisonTime, county)
			if eligibilityInfo != nil {
				infos[row.Index] = eligibilityInfo
			}
		}
	}
}

func (subject *Subject) olderThan(years int, t time.Time) bool {
	return !subject.DOB.AddDate(years, 0, 0).After(t)
}
