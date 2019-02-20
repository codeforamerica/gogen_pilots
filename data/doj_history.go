package data

import (
	"fmt"
	"strings"
	"time"
)

type DOJHistory struct {
	SubjectID       string
	Name            string
	WeakName        string
	CII             string
	DOB             time.Time
	SSN             string
	CDL             string
	Convictions     []*DOJRow
	seenConvictions map[string]bool
	OriginalCII     string
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
	}
	if row.Convicted && history.seenConvictions[row.CountOrder] {
		lastConviction := history.Convictions[len(history.Convictions)-1]
		newEndDate := lastConviction.SentenceEndDate.Add(row.SentencePartDuration)
		lastConviction.SentenceEndDate = newEndDate
	}

	if row.Convicted && !history.seenConvictions[row.CountOrder] {
		history.Convictions = append(history.Convictions, &row)
		history.seenConvictions[row.CountOrder] = true
	}

	fmt.Printf("pushing row %s with cycle=%s (%s)", row.CountOrder, row.CountOrder[0:3], row.CodeSection)

	// hash key=cycle val=hasProp64 charge

	// if conviction
	//    if prevProp64Hash[cycle]
	//         conviction.HasProp64ChargeInCycle = true
	if eligibilityFlows[county].IsProp64Charge(row.CodeSection) {
		fmt.Printf("     this is a prop64 charge!\n")
		for _, conviction := range history.Convictions {

			if conviction.CountOrder[0:3] == row.CountOrder[0:3] {
				fmt.Printf("!!found related charge for %s\n", conviction.CodeSection)
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

func (history *DOJHistory) NumberOfProp64Convictions(county string) int {
	result := 0
	for _, row := range history.Convictions {
		if eligibilityFlows[county].IsProp64Charge(row.CodeSection) {
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
