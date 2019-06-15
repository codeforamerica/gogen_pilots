package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	. "gogen/matchers"
	"gogen/utilities"
	"os"
	"strings"
	"time"
)

type DOJInformation struct {
	Rows                     [][]string
	Histories                map[string]*DOJHistory
	Eligibilities            map[int]*EligibilityInfo
	comparisonTime           time.Time
	TotalConvictions         int
	TotalConvictionsInCounty int
}

func (i *DOJInformation) generateHistories(eligibilityFlow EligibilityFlow) {
	totalRows := len(i.Rows)

	fmt.Println("Reading DOJ Data Into Memory")

	var totalTime time.Duration = 0

	for index, row := range i.Rows {
		startTime := time.Now()
		dojRow := NewDOJRow(row, index)
		if i.Histories[dojRow.SubjectID] == nil {
			i.Histories[dojRow.SubjectID] = new(DOJHistory)
		}
		i.Histories[dojRow.SubjectID].PushRow(dojRow, eligibilityFlow)

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(index+1, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
}

func (i *DOJInformation) determineEligibility(county string, eligibilityFlow EligibilityFlow) {
	for _, history := range i.Histories {
		infos := eligibilityFlow.ProcessHistory(history, i.comparisonTime, county)

		i.TotalConvictions += len(history.Convictions)
		for _, conviction := range history.Convictions {

			if conviction.County == county {
				i.TotalConvictionsInCounty++
			}
		}

		for index, info := range infos {
			i.Eligibilities[index] = info
		}
	}
}

func (i *DOJInformation) TotalIndividuals() int {
	return len(i.Histories)
}

func (i *DOJInformation) TotalRows() int {
	return len(i.Rows)
}

func (i *DOJInformation) OverallProp64ConvictionsByCodeSection() map[string]int {
	allProp64Convictions := make(map[string]int)
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			ok, codeSection := Prop64Matcher(conviction.CodeSection)
			if ok {
				allProp64Convictions[codeSection]++
			}
		}
	}
	return allProp64Convictions
}

func (i *DOJInformation) OverallRelatedConvictionsByCodeSection() map[string]int {
	relatedConvictions := make(map[string]int)
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			ok, codeSection := RelatedChargeMatcher(conviction.CodeSection)
			if ok {
				relatedConvictions[codeSection]++
			}
		}
	}
	return relatedConvictions
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSection(county string) map[string]int {
	prop64ConvictionsInCounty := make(map[string]int)
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			if conviction.County == county {
				ok, codeSection := Prop64Matcher(conviction.CodeSection)
				if ok {
					prop64ConvictionsInCounty[codeSection]++
				}
			}
		}
	}
	return prop64ConvictionsInCounty
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county string) map[string]map[string]int {
	prop64ConvictionsInCountyByCodeSectionByEligibility := make(map[string]map[string]int)
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			if conviction.County == county {
				ok, codeSection := Prop64Matcher(conviction.CodeSection)
				if ok {
					eligibilityDetermination := i.Eligibilities[conviction.Index].EligibilityDetermination
					if prop64ConvictionsInCountyByCodeSectionByEligibility[eligibilityDetermination] == nil {
						prop64ConvictionsInCountyByCodeSectionByEligibility[eligibilityDetermination] = make(map[string]int)
					}
					prop64ConvictionsInCountyByCodeSectionByEligibility[eligibilityDetermination][codeSection]++
				}

			}
		}
	}
	return prop64ConvictionsInCountyByCodeSectionByEligibility
}
func (i *DOJInformation) RelatedConvictionsInThisCountyByCodeSectionByEligibility(county string) map[string]map[string]int {
	if EligibilityFlows[county].ChecksRelatedCharges() == false {
		emptyMap := make(map[string]map[string]int)
		return emptyMap
	}
	relatedConvictionsInThisCountyByCodeSectionByEligibility := make(map[string]map[string]int)
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			if conviction.County == county {
				ok, codeSection := RelatedChargeMatcher(conviction.CodeSection)
				if ok {
					eligibilityDetermination := i.Eligibilities[conviction.Index].EligibilityDetermination
					if relatedConvictionsInThisCountyByCodeSectionByEligibility[eligibilityDetermination] == nil {
						relatedConvictionsInThisCountyByCodeSectionByEligibility[eligibilityDetermination] = make(map[string]int)
					}

					relatedConvictionsInThisCountyByCodeSectionByEligibility[eligibilityDetermination][codeSection]++
				}

			}
		}
	}
	return relatedConvictionsInThisCountyByCodeSectionByEligibility
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByEligibilityByReason(county string) map[string]map[string]int {
	prop64ConvictionsInCountyByEligibilityByReason := make(map[string]map[string]int)
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			if conviction.County == county {
				ok, _ := Prop64Matcher(conviction.CodeSection)
				if ok {
					eligibilityDetermination := i.Eligibilities[conviction.Index].EligibilityDetermination
					eligibilityReason := i.Eligibilities[conviction.Index].EligibilityReason
					if prop64ConvictionsInCountyByEligibilityByReason[eligibilityDetermination] == nil {
						prop64ConvictionsInCountyByEligibilityByReason[eligibilityDetermination] = make(map[string]int)
					}
					prop64ConvictionsInCountyByEligibilityByReason[eligibilityDetermination][eligibilityReason]++
				}

			}
		}
	}
	return prop64ConvictionsInCountyByEligibilityByReason
}

func (i *DOJInformation) CountIndividualsWithFelony() int {
	countIndividuals := 0
OuterLoop:
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			if conviction.IsFelony {
				countIndividuals++
				continue OuterLoop
			}
		}
	}
	return countIndividuals

}

func (i *DOJInformation) CountIndividualsWithConviction() int {
	countIndividuals := 0

	for _, history := range i.Histories {
		if len(history.Convictions) > 0 {
			countIndividuals++
		}
	}
	return countIndividuals

}

func (i *DOJInformation) CountIndividualsWithConvictionInLast7Years() int {
	countIndividuals := 0

OuterLoop:
	for _, history := range i.Histories {
		for _, conviction := range history.Convictions {
			if conviction.OccurredInLast7Years() {
				countIndividuals++
				continue OuterLoop
			}
		}
	}
	return countIndividuals
}

func (i *DOJInformation) CountIndividualsNoLongerHaveFelony() int {
	countIndividuals := 0
	for _, history := range i.Histories {
		countFelonies := 0
		countFeloniesReducedOrDismissed := 0
		for _, conviction := range history.Convictions {
			if conviction.IsFelony {
				countFelonies++
				if i.Eligibilities[conviction.Index] != nil {
					if determination := i.Eligibilities[conviction.Index].EligibilityDetermination; determination == "Eligible for Dismissal" || determination == "Eligible for Reduction" {
						countFeloniesReducedOrDismissed++
					}
				}
			}
		}
		if countFelonies != 0 && (countFelonies == countFeloniesReducedOrDismissed) {
			countIndividuals++
		}
	}
	return countIndividuals
}

func (i *DOJInformation) CountIndividualsNoLongerHaveConviction() int {
	countIndividuals := 0
	for _, history := range i.Histories {
		countConvictionsDismissed := 0
		for _, conviction := range history.Convictions {
			if i.Eligibilities[conviction.Index] != nil {
				if determination := i.Eligibilities[conviction.Index].EligibilityDetermination; determination == "Eligible for Dismissal" {
					countConvictionsDismissed++
				}
			}
		}
		if len(history.Convictions) != 0 &&
			(len(history.Convictions) == countConvictionsDismissed) {
			countIndividuals++
		}
	}
	return countIndividuals
}

func (i *DOJInformation) CountIndividualsNoLongerHaveConvictionInLast7Years() int {
	countIndividuals := 0
	for _, history := range i.Histories {
		convictionsInLast7Years := 0
		countConvictionsDismissedInLast7years := 0
		for _, conviction := range history.Convictions {
			if conviction.OccurredInLast7Years() {
				convictionsInLast7Years++
				if i.Eligibilities[conviction.Index] != nil {
					if determination := i.Eligibilities[conviction.Index].EligibilityDetermination; determination == "Eligible for Dismissal" {
						countConvictionsDismissedInLast7years++
					}
				}
			}
		}
		if convictionsInLast7Years != 0 &&
			(convictionsInLast7Years == countConvictionsDismissedInLast7years) {
			countIndividuals++
		}
	}
	return countIndividuals
}

func NewDOJInformation(dojFileName string, comparisonTime time.Time, county string, eligibilityFlow EligibilityFlow) *DOJInformation {
	dojFile, err := os.Open(dojFileName)
	if err != nil {
		panic(err)
	}

	bufferedReader := bufio.NewReader(dojFile)
	sourceCSV := csv.NewReader(bufferedReader)

	if includesHeaders(bufferedReader) {
		bufferedReader.ReadLine() // read and discard header row
	}

	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}
	info := DOJInformation{
		Rows:           rows,
		Histories:      make(map[string]*DOJHistory),
		Eligibilities:  make(map[int]*EligibilityInfo),
		comparisonTime: comparisonTime,
	}

	info.generateHistories(eligibilityFlow)
	info.determineEligibility(county, eligibilityFlow)

	return &info
}

func isHeaderRow(rowString string) bool {
	return strings.HasPrefix(rowString, "RECORD_ID")
}

func includesHeaders(reader *bufio.Reader) bool {
	firstRowBytes, err := reader.Peek(128)

	if err != nil {
		panic(err)
	}

	firstRow := string(firstRowBytes)

	return isHeaderRow(firstRow)
}
