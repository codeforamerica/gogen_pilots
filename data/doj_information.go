package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gogen/matchers"
	"gogen/utilities"
	"os"
	"time"
)

type DOJInformation struct {
	Rows                     [][]string
	Subjects                 map[string]*Subject
	Eligibilities            map[int]*EligibilityInfo
	comparisonTime           time.Time
	TotalConvictions         int
	TotalConvictionsInCounty int
}

func (i *DOJInformation) generateHistories(county string) {
	totalRows := len(i.Rows)

	fmt.Println("Reading DOJ Data Into Memory")

	var totalTime time.Duration = 0

	for index, row := range i.Rows {
		startTime := time.Now()
		dojRow := NewDOJRow(row, index)
		if i.Subjects[dojRow.SubjectID] == nil {
			i.Subjects[dojRow.SubjectID] = new(Subject)
		}
		i.Subjects[dojRow.SubjectID].PushRow(dojRow, county)

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(index + 1, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
}

func (i *DOJInformation) TotalIndividuals() int {
	return len(i.Subjects)
}

func (i *DOJInformation) TotalRows() int {
	return len(i.Rows)
}

func (i *DOJInformation) OverallProp64ConvictionsByCodeSection() map[string]int {
	allProp64Convictions := make(map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			ok, codeSection := matchers.Prop64Matcher(conviction.CodeSection)
			if ok {
				allProp64Convictions[codeSection]++
			}
		}
	}
	return allProp64Convictions
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSection(county string) map[string]int {
	prop64ConvictionsInCounty := make(map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, codeSection := matchers.Prop64Matcher(conviction.CodeSection)
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
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, codeSection := matchers.Prop64Matcher(conviction.CodeSection)
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
	relatedConvictionsInThisCountyByCodeSectionByEligibility := make(map[string]map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, codeSection := matchers.RelatedChargeMatcher(conviction.CodeSection)
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
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, _ := matchers.Prop64Matcher(conviction.CodeSection)
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

func (i *DOJInformation) CountIndividualsNoLongerHaveFelony() int {
	countIndividuals := 0
	for _, subject := range i.Subjects {
		countFelonies := 0
		countFeloniesReducedOrDismissed := 0
		for _, conviction := range subject.Convictions {
			if conviction.IsFelony {
				countFelonies++
				if i.Eligibilities[conviction.Index] != nil {
					if determination := i.Eligibilities[conviction.Index].EligibilityDetermination;
						determination == "Eligible for Dismissal" {
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
	for _, subject := range i.Subjects {
		countConvictionsDismissed := 0
		for _, conviction := range subject.Convictions {
			if i.Eligibilities[conviction.Index] != nil {
				if determination := i.Eligibilities[conviction.Index].EligibilityDetermination;
					determination == "Eligible for Dismissal" {
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

func (i *DOJInformation) CountIndividualsNoLongerHaveConvictionInLast7Years() int {
	countIndividuals := 0
	for _, subject := range i.Subjects {
		convictionsInLast7Years := 0
		countConvictionsDismissedInLast7years := 0
		for _, conviction := range subject.Convictions {
			if conviction.OccurredInLast7Years() {
				convictionsInLast7Years++
				if i.Eligibilities[conviction.Index] != nil {
					if determination := i.Eligibilities[conviction.Index].EligibilityDetermination;
						determination == "Eligible for Dismissal" {
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

func (i *DOJInformation) determineEligibility(county string) {
	for _, subject := range i.Subjects {
		infos := EligibilityFlows[county].ProcessHistory(subject, i.comparisonTime)

		i.TotalConvictions += len(subject.Convictions)
		for _, conviction := range subject.Convictions {

			if conviction.County == county {
				i.TotalConvictionsInCounty++
			}
		}

		for index, info := range infos {
			i.Eligibilities[index] = info
		}
	}
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

func NewDOJInformation(dojFileName string, comparisonTime time.Time, county string) *DOJInformation {
	dojFile, err := os.Open(dojFileName)
	if err != nil {
		panic(err)
	}

	bufferedReader := bufio.NewReader(dojFile)
	bufferedReader.ReadLine() // read and discard header row

	sourceCSV := csv.NewReader(bufferedReader)
	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}

	info := DOJInformation{
		Rows:           rows,
		Subjects:       make(map[string]*Subject),
		Eligibilities:  make(map[int]*EligibilityInfo),
		comparisonTime: comparisonTime,
	}

	info.generateHistories(county)
	info.determineEligibility(county)

	return &info
}
