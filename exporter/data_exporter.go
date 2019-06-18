package exporter

import (
	"fmt"
	"gogen/data"
	"sort"
)

type DataExporter struct {
	dojInformation                           *data.DOJInformation
	dismissAllProp64DojInformation           *data.DOJInformation
	dismissAllProp64AndRelatedDojInformation *data.DOJInformation
	outputDOJWriter                          DOJWriter
	outputCondensedDOJWriter                 DOJWriter
	outputProp64ConvictionsDOJWriter         DOJWriter
}

func NewDataExporter(
	dojInformation *data.DOJInformation,
	dismissAllProp64DojInformation *data.DOJInformation,
	dismissAllProp64AndRelatedDojInformation *data.DOJInformation,
	outputDOJWriter DOJWriter,
	outputCondensedDOJWriter DOJWriter,
	outputProp64ConvictionsDOJWriter DOJWriter,
) DataExporter {
	return DataExporter{
		dojInformation:                           dojInformation,
		dismissAllProp64DojInformation:           dismissAllProp64DojInformation,
		dismissAllProp64AndRelatedDojInformation: dismissAllProp64AndRelatedDojInformation,
		outputDOJWriter:                          outputDOJWriter,
		outputCondensedDOJWriter:                 outputCondensedDOJWriter,
		outputProp64ConvictionsDOJWriter:         outputProp64ConvictionsDOJWriter,
	}
}

func (d *DataExporter) Export(county string) {
	fmt.Printf("Processing Subjects\n")

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

func (d *DataExporter) PrintAggregateStatistics(county string) {
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
	fmt.Printf("%d individuals who had a felony will no longer have a felony on their record\n", d.dojInformation.CountIndividualsNoLongerHaveFelony())
	fmt.Printf("%d individuals who had convictions will no longer have any convictions on their record\n", d.dojInformation.CountIndividualsNoLongerHaveConviction())
	fmt.Printf("%d individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years())
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
	codeSections := make([]string, 0, len(resultsByCodeSectionByEligibility))
	for key := range resultsByCodeSectionByEligibility {
		codeSections = append(codeSections, key)
	}
	sort.Strings(codeSections)

	for _, codeSection := range codeSections {
		fmt.Printf("\n%v", codeSection)

		total := 0
		eligibilityMap := resultsByCodeSectionByEligibility[codeSection]
		eligibilities := make([]string, 0, len(eligibilityMap))
		for key := range eligibilityMap {
			eligibilities = append(eligibilities, key)
		}
		sort.Strings(eligibilities)

		for _, eligibility := range eligibilities {

			total += eligibilityMap[eligibility]
			fmt.Printf("\nFound %v %v convictions that are %v", eligibilityMap[eligibility], eligibility, codeSection)
		}
		fmt.Printf("\nFound %v convictions total that are %v\n", total, codeSection)
	}
}

func printSummaryByEligibilityByReason(resultsByEligibilityByReason map[string]map[string]int) {
	determinations := make([]string, 0, len(resultsByEligibilityByReason))
	for key := range resultsByEligibilityByReason {
		determinations = append(determinations, key)
	}
	sort.Strings(determinations)

	for _, determination := range determinations {
		fmt.Printf("\n%v", determination)

		total := 0
		reasonMap := resultsByEligibilityByReason[determination]
		reasons := make([]string, 0, len(reasonMap))
		for key := range reasonMap {
			reasons = append(reasons, key)
		}
		sort.Strings(reasons)

		for _, reason := range reasons {
			total += reasonMap[reason]
			fmt.Printf("\nFound %v convictions with eligibility reason %v", reasonMap[reason], reason)
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
