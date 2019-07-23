package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gogen/matchers"
	"gogen/utilities"
	"os"
	"strings"
	"time"
)

type DOJInformation struct {
	Rows                 [][]string
	Subjects             map[string]*Subject
	comparisonTime       time.Time
	checksRelatedCharges bool
}

func (i *DOJInformation) aggregateSubjects(eligibilityFlow EligibilityFlow) {
	totalRows := len(i.Rows)

	fmt.Println("Reading DOJ Data Into Memory")

	var totalTime time.Duration = 0

	for index, row := range i.Rows {
		startTime := time.Now()
		dojRow := NewDOJRow(row, index)
		if i.Subjects[dojRow.SubjectID] == nil {
			i.Subjects[dojRow.SubjectID] = new(Subject)
		}
		i.Subjects[dojRow.SubjectID].PushRow(dojRow, eligibilityFlow)

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(index+1, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
}

func (i *DOJInformation) DetermineEligibility(county string, eligibilityFlow EligibilityFlow) map[int]*EligibilityInfo {
	eligibilities := make(map[int]*EligibilityInfo)
	for _, subject := range i.Subjects {
		infos := eligibilityFlow.ProcessSubject(subject, i.comparisonTime, county)
		for index, info := range infos {
			eligibilities[index] = info
		}
	}
	return eligibilities
}

func (i *DOJInformation) TotalIndividuals() int {
	return len(i.Subjects)
}

func (i *DOJInformation) TotalRows() int {
	return len(i.Rows)
}

func (i *DOJInformation) TotalConvictions() int {
	totalConvictions := 0
	for _, subject := range i.Subjects {
		totalConvictions += len(subject.Convictions)
	}
	return totalConvictions
}

func (i *DOJInformation) TotalConvictionsInCounty(county string) int {
	totalConvictionsInCounty := 0
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				totalConvictionsInCounty++
			}
		}
	}
	return totalConvictionsInCounty
}

func (i *DOJInformation) OverallProp64ConvictionsByCodeSection() map[string]int {
	allProp64Convictions := make(map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			ok, codeSection := matchers.ExtractProp64Section(conviction.CodeSection)
			if ok {
				allProp64Convictions[codeSection]++
			}
		}
	}
	return allProp64Convictions
}

func (i *DOJInformation) OverallRelatedConvictionsByCodeSection() map[string]int {
	relatedConvictions := make(map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			ok, codeSection := matchers.ExtractRelatedChargeSection(conviction.CodeSection)
			if ok {
				relatedConvictions[codeSection]++
			}
		}
	}
	return relatedConvictions
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSection(county string) map[string]int {
	prop64ConvictionsInCounty := make(map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, codeSection := matchers.ExtractProp64Section(conviction.CodeSection)
				if ok {
					prop64ConvictionsInCounty[codeSection]++
				}
			}
		}
	}
	return prop64ConvictionsInCounty
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county string, eligibilities map[int]*EligibilityInfo) map[string]map[string]int {
	prop64ConvictionsInCountyByCodeSectionByEligibility := make(map[string]map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, codeSection := matchers.ExtractProp64Section(conviction.CodeSection)
				if ok {
					eligibilityDetermination := eligibilities[conviction.Index].EligibilityDetermination
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
func (i *DOJInformation) RelatedConvictionsInThisCountyByCodeSectionByEligibility(county string, eligibilities map[int]*EligibilityInfo) map[string]map[string]int {
	if !i.checksRelatedCharges {
		emptyMap := make(map[string]map[string]int)
		return emptyMap
	}
	relatedConvictionsInThisCountyByCodeSectionByEligibility := make(map[string]map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, codeSection := matchers.ExtractRelatedChargeSection(conviction.CodeSection)
				if ok {
					eligibilityDetermination := eligibilities[conviction.Index].EligibilityDetermination
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

func (i *DOJInformation) Prop64ConvictionsInThisCountyByEligibilityByReason(county string, eligibilities map[int]*EligibilityInfo) map[string]map[string]int {
	prop64ConvictionsInCountyByEligibilityByReason := make(map[string]map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				if matchers.IsProp64Charge(conviction.CodeSection) {
					eligibilityDetermination := eligibilities[conviction.Index].EligibilityDetermination
					eligibilityReason := eligibilities[conviction.Index].EligibilityReason
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
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
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

	for _, subject := range i.Subjects {
		if len(subject.Convictions) > 0 {
			countIndividuals++
		}
	}
	return countIndividuals

}

func (i *DOJInformation) CountIndividualsWithConvictionInLast7Years() int {
	countIndividuals := 0

OuterLoop:
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.OccurredInLast7Years() {
				countIndividuals++
				continue OuterLoop
			}
		}
	}
	return countIndividuals
}

func (i *DOJInformation) CountIndividualsNoLongerHaveFelony(eligibilities map[int]*EligibilityInfo) int {
	countIndividuals := 0
	for _, subject := range i.Subjects {
		countFelonies := 0
		countFeloniesReducedOrDismissed := 0
		for _, conviction := range subject.Convictions {
			if conviction.IsFelony {
				countFelonies++
				if eligibilities[conviction.Index] != nil {
					if determination := eligibilities[conviction.Index].EligibilityDetermination; determination == "Eligible for Dismissal" || determination == "Eligible for Reduction" {
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

func (i *DOJInformation) CountIndividualsNoLongerHaveConviction(eligibilities map[int]*EligibilityInfo) int {
	countIndividuals := 0
	for _, subject := range i.Subjects {
		countConvictionsDismissed := 0
		for _, conviction := range subject.Convictions {
			if eligibilities[conviction.Index] != nil {
				if determination := eligibilities[conviction.Index].EligibilityDetermination; determination == "Eligible for Dismissal" {
					countConvictionsDismissed++
				}
			}
		}
		if len(subject.Convictions) != 0 &&
			(len(subject.Convictions) == countConvictionsDismissed) {
			countIndividuals++
		}
	}
	return countIndividuals
}

func (i *DOJInformation) CountIndividualsNoLongerHaveConvictionInLast7Years(eligibilities map[int]*EligibilityInfo) int {
	countIndividuals := 0
	for _, subject := range i.Subjects {
		convictionsInLast7Years := 0
		countConvictionsDismissedInLast7years := 0
		for _, conviction := range subject.Convictions {
			if conviction.OccurredInLast7Years() {
				convictionsInLast7Years++
				if eligibilities[conviction.Index] != nil {
					if determination := eligibilities[conviction.Index].EligibilityDetermination; determination == "Eligible for Dismissal" {
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

func NewDOJInformation(dojFileName string, comparisonTime time.Time, eligibilityFlow EligibilityFlow) *DOJInformation {
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
		Rows:                 rows,
		Subjects:             make(map[string]*Subject),
		comparisonTime:       comparisonTime,
		checksRelatedCharges: eligibilityFlow.ChecksRelatedCharges(),
	}

	info.aggregateSubjects(eligibilityFlow)

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
