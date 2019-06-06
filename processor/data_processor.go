package processor

import (
	"fmt"
	"gogen/data"
	"regexp"
	"sort"
)

type DataProcessor struct {
	dojInformation                   *data.DOJInformation
	dismissAllProp64DojInformation   *data.DOJInformation
	dismissAllProp64AndRelatedDojInformation   *data.DOJInformation
	outputDOJWriter                  DOJWriter
	outputCondensedDOJWriter         DOJWriter
	outputProp64ConvictionsDOJWriter DOJWriter
	prop64Matcher                    *regexp.Regexp
}

func NewDataProcessor(
	dojInformation *data.DOJInformation,
	dismissAllProp64DojInformation *data.DOJInformation,
	dismissAllProp64AndRelatedDojInformation *data.DOJInformation,
	outputDOJWriter DOJWriter,
	outputCondensedDOJWriter DOJWriter,
	outputProp64ConvictionsDOJWriter DOJWriter,
) DataProcessor {
	return DataProcessor{
		dojInformation:                   dojInformation,
		dismissAllProp64DojInformation:   dismissAllProp64DojInformation,
		dismissAllProp64AndRelatedDojInformation:   dismissAllProp64AndRelatedDojInformation,
		outputDOJWriter:                  outputDOJWriter,
		outputCondensedDOJWriter:         outputCondensedDOJWriter,
		outputProp64ConvictionsDOJWriter: outputProp64ConvictionsDOJWriter,
	}
}


func (d *DataProcessor) Process(county string) {
	fmt.Printf("Processing Histories\n")

	for i, row := range d.dojInformation.Rows {
		d.outputDOJWriter.WriteEntryWithEligibilityInfo(row, d.dojInformation.Eligibilities[i])
		d.outputCondensedDOJWriter.WriteCondensedEntryWithEligibilityInfo(row, d.dojInformation.Eligibilities[i])
		if d.dojInformation.Eligibilities[i] != nil {
			d.outputProp64ConvictionsDOJWriter.WriteEntryWithEligibilityInfo(row, d.dojInformation.Eligibilities[i])
		}
	}

	d.outputDOJWriter.Flush()
	d.outputCondensedDOJWriter.Flush()
	d.outputProp64ConvictionsDOJWriter.Flush()

	d.PrintAggregateStatistics(county)
}

func (d *DataProcessor) PrintAggregateStatistics(county string) {
	fmt.Println()
	fmt.Println("----------- Overall summary of DOJ file --------------------")
	fmt.Printf("Found %d Total rows in DOJ file\n", d.dojInformation.TotalRows())
	fmt.Printf("Found %d Total individuals in DOJ file\n", d.dojInformation.TotalIndividuals())
	fmt.Printf("Found %d Total convictions in DOJ file\n", d.dojInformation.TotalConvictions)
	fmt.Printf("Found %d convictions in this county\n", d.dojInformation.TotalConvictionsInCounty)

	fmt.Println()
	fmt.Printf("----------- Prop64 Convictions Overall--------------------")
	printSummaryByCodeSection("total", d.dojInformation.OverallProp64ConvictionsByCodeSection())
	fmt.Println()

	fmt.Printf("----------- Prop64 Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", d.dojInformation.Prop64ConvictionsInThisCountyByCodeSection(county))

	printSummaryByCodeSectionByEligibility(d.dojInformation.Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county))

	fmt.Println()
	fmt.Println("----------- Eligibility Reasons --------------------")
	printSummaryByEligibilityByReason(d.dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county))
	fmt.Println()
	fmt.Println()
	fmt.Printf("----------- Prop64 Related Convictions In This County --------------------")
	printSummaryByCodeSection("in this county", d.dojInformation.OverallRelatedConvictionsByCodeSection())
	fmt.Println()
	printSummaryByCodeSectionByEligibility(d.dojInformation.RelatedConvictionsInThisCountyByCodeSectionByEligibility(county))
	fmt.Println()

	fmt.Println("----------- Impact to individuals --------------------")
	fmt.Printf("%d individuals currently have a felony on their record\n", d.dojInformation.CountIndividualsWithFelony())
	fmt.Printf("%d individuals currently have convictions on their record\n", d.dojInformation.CountIndividualsWithConviction())
	fmt.Printf("%d individuals currently have convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsWithConvictionInLast7Years())
	fmt.Println()

	fmt.Println("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------")
	fmt.Printf("%d individual(s) who had a felony will no longer have a felony on their record\n", d.dojInformation.CountIndividualsNoLongerHaveFelony())
	fmt.Printf("%d individual(s) who had convictions will no longer have any convictions on their record\n", d.dojInformation.CountIndividualsNoLongerHaveConviction())
	fmt.Printf("%d individual(s) who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years())
	fmt.Println()
	fmt.Println("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.dismissAllProp64DojInformation.CountIndividualsNoLongerHaveFelony())
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.dismissAllProp64DojInformation.CountIndividualsNoLongerHaveConviction())
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.dismissAllProp64DojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years())
	fmt.Println()
	fmt.Println("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------")
	fmt.Printf("%d individuals will no longer have a felony on their record\n", d.dismissAllProp64AndRelatedDojInformation.CountIndividualsNoLongerHaveFelony())
	fmt.Printf("%d individuals will no longer have any convictions on their record\n", d.dismissAllProp64AndRelatedDojInformation.CountIndividualsNoLongerHaveConviction())
	fmt.Printf("%d individuals will no longer have any convictions on their record in the last 7 years\n", d.dismissAllProp64AndRelatedDojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years())
}

func printSummaryByCodeSection(description string, resultsByCodeSection map[string]int) {
	fmt.Printf("\nFound %d convictions %s\n", sumValues(resultsByCodeSection), description)
	formatString := fmt.Sprintf("Found %%d %%s convictions %s\n", description)
	printMap(formatString, resultsByCodeSection)
}

func printSummaryByCodeSectionByEligibility(resultsByCodeSectionByEligibility map[string]map[string]int) {

	for determination, codeSectionMap := range resultsByCodeSectionByEligibility {
		fmt.Printf("\n%v", determination)

		total := 0
		for codeSection, value := range codeSectionMap {

			total += value
			fmt.Printf("\nFound %v %v convictions that are %v", value, codeSection, determination)
		}
		fmt.Printf("\nFound %v convictions total that are %v\n", total, determination)
	}
}

func printSummaryByEligibilityByReason(resultsByCodeSectionByEligibility map[string]map[string]int) {

	for determination, eligibilityMap := range resultsByCodeSectionByEligibility {
		fmt.Printf("\n%v", determination)

		total := 0
		for eligibilityReason, value := range eligibilityMap {

			total += value
			fmt.Printf("\nFound %v convictions with eligibility reason %v", value, eligibilityReason)
		}
		fmt.Println()
	}
}

func printMap(formatString string, values map[string]int) {
	keys := getSortedKeys(values)

	for _, key := range keys {
		fmt.Printf(formatString, values[key], key)
	}
}

func getSortedKeys(mapWithStringKeys map[string]int) []string {
	keys := make([]string, 0, len(mapWithStringKeys))
	for key := range mapWithStringKeys {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sumValues(mapOfInts map[string]int) int {
	total := 0
	for _, value := range mapOfInts {
		total += value
	}
	return total
}
