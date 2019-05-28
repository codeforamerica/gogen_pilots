package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	. "gogen/matchers"
	"gogen/utilities"
	"os"
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

func (i *DOJInformation) generateHistories(county string) {
	currentRowIndex := 0.0
	totalRows := float64(len(i.Rows))

	fmt.Println("Reading DOJ Data Into Memory")

	var totalTime time.Duration = 0

	for index, row := range i.Rows {
		startTime := time.Now()
		dojRow := NewDOJRow(row, index)
		if i.Histories[dojRow.SubjectID] == nil {
			i.Histories[dojRow.SubjectID] = new(DOJHistory)
		}
		i.Histories[dojRow.SubjectID].PushRow(dojRow, county)
		currentRowIndex++

		totalTime += time.Since(startTime)

		utilities.PrintProgressBar(currentRowIndex, totalRows, totalTime, "")
	}
	fmt.Println("\nComplete...")
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
	//prop64ConvictionsByCodeSection := make(map[string]int)
	for _, history := range i.Histories {
		fmt.Println("in history")
		for convictionIndex, conviction := range history.Convictions {
			fmt.Println("in convictions")
			if conviction.County == county {
				fmt.Println("in county match")
				ok, codeSection := Prop64Matcher(conviction.CodeSection)
				if ok {
					fmt.Println("in okay")
					eligibilityDetermination := history.EligibilityInfos[convictionIndex].EligibilityDetermination
					//if prop64ConvictionsInCountyByCodeSectionByEligibility[eligibilityDetermination] == nil {
					//	prop64ConvictionsInCountyByCodeSectionByEligibility[eligibilityDetermination] = map[string]int{codeSection: 1}
					//} else {
					prop64ConvictionsInCountyByCodeSectionByEligibility[eligibilityDetermination][codeSection]++
				}
			}
		}
	}
	fmt.Printf("map value: %v", prop64ConvictionsInCountyByCodeSectionByEligibility)
	return prop64ConvictionsInCountyByCodeSectionByEligibility
}

func (i *DOJInformation) determineEligibility(county string) {
	for _, history := range i.Histories {
		infos := EligibilityFlows[county].ProcessHistory(history, i.comparisonTime)

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

func NewDOJInformation(dojFileName string, comparisonTime time.Time, county string) (*DOJInformation, error) {
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
		Histories:      make(map[string]*DOJHistory),
		Eligibilities:  make(map[int]*EligibilityInfo),
		comparisonTime: comparisonTime,
	}

	info.generateHistories(county)
	info.determineEligibility(county)

	return &info, nil
}
