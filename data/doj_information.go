package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gogen/matchers"
	"gogen/utilities"
	"os"
	"sort"
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
	return i.countByCodeSectionFilteredMatchedConvictions("", emptyFilter, matchers.ExtractProp64Section)
}

func (i *DOJInformation) OverallRelatedConvictionsByCodeSection() map[string]int {
	return i.countByCodeSectionFilteredMatchedConvictions("", emptyFilter, matchers.ExtractRelatedChargeSection)
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSection(county string) map[string]int {
	return i.countByCodeSectionFilteredMatchedConvictions(county, countyFilter, matchers.ExtractProp64Section)
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county string, eligibilities map[int]*EligibilityInfo) map[string]map[string]int {
	return i.countByCodeSectionAndEligibilityFilteredMatchedConvictions(county, eligibilities, countyFilter, matchers.ExtractProp64Section, countByCodeSectionAndEligibilityDetermination)
}

func (i *DOJInformation) RelatedConvictionsInThisCountyByCodeSectionByEligibility(county string, eligibilities map[int]*EligibilityInfo) map[string]map[string]int {
	if !i.checksRelatedCharges {
		emptyMap := make(map[string]map[string]int)
		return emptyMap
	}
	return i.countByCodeSectionAndEligibilityFilteredMatchedConvictions(county, eligibilities, countyFilter, matchers.ExtractRelatedChargeSection, countByCodeSectionAndEligibilityDetermination)
}

func (i *DOJInformation) Prop64ConvictionsInThisCountyByEligibilityByReason(county string, eligibilities map[int]*EligibilityInfo) map[string]map[string]int {
	return i.countByCodeSectionAndEligibilityFilteredMatchedConvictions(county, eligibilities, countyFilter, matchers.ExtractProp64Section, countByEligibilityDeterminationAndReason)
}

func (i *DOJInformation) EarliestProp64ConvictionDateInThisCounty(county string) time.Time {
	var convictionDates = TimeSlice{}
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if conviction.County == county {
				ok, _ := matchers.ExtractProp64Section(conviction.CodeSection)
				if ok {
					convictionDates = append(convictionDates, conviction.DispositionDate)
				}
			}
		}
	}
	sort.Sort(convictionDates)
	if len(convictionDates) > 0 {
		return convictionDates[0]
	} else {
		return time.Now()
	}
}

func (i *DOJInformation) CountIndividualsWithProp64ConvictionInCounty(county string) int {
	pro64AndCountyFilter := func(d *DOJRow) bool {
		return d.County == county && matchers.IsProp64Charge(d.CodeSection)
	}
	return i.countIndividualsFilteredByConviction(pro64AndCountyFilter)
}

func (i *DOJInformation) CountIndividualsWithFelony() int {
	return i.countIndividualsFilteredByConviction(IsFelonyFilter)
}

func (i *DOJInformation) CountIndividualsWithConviction() int {
	return i.countIndividualsFilteredByConviction(hasConvictionFilter)
}

func (i *DOJInformation) CountIndividualsWithSomeRelief(eligibilities map[int]*EligibilityInfo) int {
	countIndividuals := 0
OuterLoop:
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if eligibilities[conviction.Index] != nil {
				if reducedOrDismissedFilter(eligibilities[conviction.Index]) {
					countIndividuals++
					continue OuterLoop
				}
			}
		}
	}
	return countIndividuals
}

func (i *DOJInformation) CountIndividualsWithConvictionInLast7Years() int {
	return i.countIndividualsFilteredByConviction(occurredInLast7YearsFilter)
}

func (i *DOJInformation) CountIndividualsNoLongerHaveFelony(eligibilities map[int]*EligibilityInfo) int {
	return i.countIndividualsFilteredByFullRelief(eligibilities, IsFelonyFilter, reducedOrDismissedFilter)
}

func (i *DOJInformation) CountIndividualsNoLongerHaveConviction(eligibilities map[int]*EligibilityInfo) int {
	return i.countIndividualsFilteredByFullRelief(eligibilities, hasConvictionFilter, dismissedFilter)
}

func (i *DOJInformation) CountIndividualsNoLongerHaveConvictionInLast7Years(eligibilities map[int]*EligibilityInfo) int {
	return i.countIndividualsFilteredByFullRelief(eligibilities, occurredInLast7YearsFilter, dismissedFilter)
}

func NewDOJInformation(dojFileName string, comparisonTime time.Time, eligibilityFlow EligibilityFlow) (*DOJInformation, error) {
	dojFile, err := os.Open(dojFileName)
	if err != nil {
		return nil, err
	}

	bufferedReader := bufio.NewReader(dojFile)
	sourceCSV := csv.NewReader(bufferedReader)

	hasHeaders, err := includesHeaders(bufferedReader)
	if err != nil {
		return nil, err
	}
	if hasHeaders {
		bufferedReader.ReadLine() // read and discard header row
	}

	rows, err := sourceCSV.ReadAll()
	if err != nil {
		return nil, err
	}
	info := DOJInformation{
		Rows:                 rows,
		Subjects:             make(map[string]*Subject),
		comparisonTime:       comparisonTime,
		checksRelatedCharges: eligibilityFlow.ChecksRelatedCharges(),
	}

	info.aggregateSubjects(eligibilityFlow)

	return &info, nil
}

func (i *DOJInformation) countByCodeSectionFilteredMatchedConvictions(
	county string,
	filter func(county string, conviction *DOJRow) bool,
	matcher func(codeSection string) (bool, string)) map[string]int {
	convictionMap := make(map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if filter(county, conviction) {
				ok, codeSection := matcher(conviction.CodeSection)
				if ok {
					convictionMap[codeSection]++
				}

			}
		}
	}
	return convictionMap
}

func (i *DOJInformation) countByCodeSectionAndEligibilityFilteredMatchedConvictions(
	county string,
	eligibilities map[int]*EligibilityInfo,
	filter func(county string, conviction *DOJRow) bool,
	matcher func(codeSection string) (bool, string),
	mapper func(conviction *DOJRow, codeSection string, eligibilities map[int]*EligibilityInfo, convictionMap map[string]map[string]int) map[string]map[string]int) map[string]map[string]int {
	convictionMap := make(map[string]map[string]int)
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if filter(county, conviction) {
				ok, codeSection := matcher(conviction.CodeSection)
				if ok {
					convictionMap = mapper(conviction, codeSection, eligibilities, convictionMap)
				}

			}
		}
	}
	return convictionMap
}

func (i *DOJInformation) countIndividualsFilteredByConviction(filter func(conviction *DOJRow) bool) int {
	countIndividuals := 0
OuterLoop:
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if filter(conviction) {
				countIndividuals++
				continue OuterLoop
			}
		}
	}
	return countIndividuals
}

func (i *DOJInformation) countIndividualsFilteredByFullRelief(
	eligibilities map[int]*EligibilityInfo,
	convictionFilter func(conviction *DOJRow) bool,
	reliefFilter func(eligibility *EligibilityInfo) bool) int {
	countIndividuals := 0
	for _, subject := range i.Subjects {
		countConvictions := 0
		countRelief := 0
		for _, conviction := range subject.Convictions {
			if convictionFilter(conviction) {
				countConvictions++
				if eligibilities[conviction.Index] != nil {
					if reliefFilter(eligibilities[conviction.Index]) {
						countRelief++
					}
				}
			}
		}
		if countConvictions != 0 && (countConvictions == countRelief) {
			countIndividuals++
		}
	}
	return countIndividuals
}

func countByCodeSectionAndEligibilityDetermination(
	conviction *DOJRow,
	codeSection string,
	eligibilities map[int]*EligibilityInfo,
	convictionMap map[string]map[string]int) map[string]map[string]int {
	eligibilityDetermination := eligibilities[conviction.Index].EligibilityDetermination
	if convictionMap[eligibilityDetermination] == nil {
		convictionMap[eligibilityDetermination] = make(map[string]int)
	}
	convictionMap[eligibilityDetermination][codeSection]++
	return convictionMap
}

func countByEligibilityDeterminationAndReason(
	conviction *DOJRow,
	_ string,
	eligibilities map[int]*EligibilityInfo,
	convictionMap map[string]map[string]int) map[string]map[string]int {
	eligibilityDetermination := eligibilities[conviction.Index].EligibilityDetermination
	eligibilityReason := eligibilities[conviction.Index].EligibilityReason
	if convictionMap[eligibilityDetermination] == nil {
		convictionMap[eligibilityDetermination] = make(map[string]int)
	}
	convictionMap[eligibilityDetermination][eligibilityReason]++
	return convictionMap
}

func (i *DOJInformation) TotalConvictionsInCountyFiltered(county string, convictionFilter func(conviction *DOJRow) bool, matcher func(codeSection string) bool) int {
	result := 0
	for _, subject := range i.Subjects {
		for _, conviction := range subject.Convictions {
			if countyFilter(county, conviction) && convictionFilter(conviction) && matcher(conviction.CodeSection) {
				result++
			}
		}
	}
	return result
}

func countyFilter(county string, conviction *DOJRow) bool {
	return conviction.County == county
}

func emptyFilter(_ string, _ *DOJRow) bool {
	return true
}

func hasConvictionFilter(conviction *DOJRow) bool {
	return conviction != nil
}

func IsFelonyFilter(conviction *DOJRow) bool {
	return conviction.IsFelony
}

func IsNotFelonyFilter(conviction *DOJRow) bool {
	return !conviction.IsFelony
}

func occurredInLast7YearsFilter(conviction *DOJRow) bool {
	return conviction.OccurredInLast7Years()
}

func reducedOrDismissedFilter(eligibility *EligibilityInfo) bool {
	determination := eligibility.EligibilityDetermination
	return determination == "Eligible for Dismissal" || determination == "Eligible for Reduction"
}

func dismissedFilter(eligibility *EligibilityInfo) bool {
	determination := eligibility.EligibilityDetermination
	return determination == "Eligible for Dismissal"
}

func isHeaderRow(rowString string) bool {
	return strings.HasPrefix(rowString, "RECORD_ID")
}

func includesHeaders(reader *bufio.Reader) (bool, error) {
	firstRowBytes, err := reader.Peek(128)

	if err != nil {
		return false, err
	}

	firstRow := string(firstRowBytes)

	return isHeaderRow(firstRow), nil
}

type TimeSlice []time.Time

func (s TimeSlice) Len() int {
	return len(s)
}
func (s TimeSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TimeSlice) Less(i, j int) bool {
	return s[i].Before(s[j])
}
