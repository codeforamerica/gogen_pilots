package exporter

import (
	"fmt"
	"gogen/data"
	"gogen/matchers"
	"gogen/utilities"
	"io"
	"sort"
	"strings"
	"time"
)

type DataExporter struct {
	dojInformation                          *data.DOJInformation
	normalFlowEligibilities                 map[int]*data.EligibilityInfo
	dismissAllProp64Eligibilities           map[int]*data.EligibilityInfo
	dismissAllProp64AndRelatedEligibilities map[int]*data.EligibilityInfo
	outputDOJWriter                         DOJWriter
	outputCondensedDOJWriter                DOJWriter
	outputProp64ConvictionsDOJWriter        DOJWriter
	aggregateStatsWriter                    io.Writer
	outputJsonFilePath                      string
}

type Summary struct {
	County                                      string         `json:"county"`
	EarliestConviction                          time.Time      `json:"earliestConviction"`
	LineCount                                   int            `json:"lineCount"`
	ProcessingTimeInSeconds                     float64        `json:"processingTimeInSeconds"`
	ReliefWithCurrentEligibilityChoices         map[string]int `json:"reliefWithCurrentEligibilityChoices"`
	ReliefWithDismissAllProp64                  map[string]int `json:"reliefWithDismissAllProp64"`
	Prop64ConvictionsCountInCountyByCodeSection map[string]int `json:"prop64ConvictionsCountInCountyByCodeSection"`
	SubjectsWithProp64ConvictionCountInCounty   int            `json:"subjectsWithProp64ConvictionCountInCounty"`
	Prop64FelonyConvictionsCountInCounty        int            `json:"prop64FelonyConvictionsCountInCounty"`
	Prop64NonFelonyConvictionsCountInCounty     int            `json:"prop64NonFelonyConvictionsCountInCounty"`
	SubjectsWithSomeReliefCount                 int            `json:"subjectsWithSomeReliefCount"`
	ConvictionDismissalCountByCodeSection       map[string]int `json:"convictionDismissalCountByCodeSection"`
	ConvictionReductionCountByCodeSection       map[string]int `json:"convictionReductionCountByCodeSection"`
	ConvictionDismissalCountByAdditionalRelief  map[string]int `json:"convictionDismissalCountByAdditionalRelief"`
}

func NewDataExporter(
	dojInformation *data.DOJInformation,
	countyEligibilities map[int]*data.EligibilityInfo,
	dismissAllProp64Eligibilities map[int]*data.EligibilityInfo,
	dismissAllProp64AndRelatedEligibilities map[int]*data.EligibilityInfo,
	outputDOJWriter DOJWriter,
	outputCondensedDOJWriter DOJWriter,
	outputProp64ConvictionsDOJWriter DOJWriter,
	aggregateStatsWriter io.Writer,
) DataExporter {

	return DataExporter{
		dojInformation:                          dojInformation,
		normalFlowEligibilities:                 countyEligibilities,
		dismissAllProp64Eligibilities:           dismissAllProp64Eligibilities,
		dismissAllProp64AndRelatedEligibilities: dismissAllProp64AndRelatedEligibilities,
		outputDOJWriter:                         outputDOJWriter,
		outputCondensedDOJWriter:                outputCondensedDOJWriter,
		outputProp64ConvictionsDOJWriter:        outputProp64ConvictionsDOJWriter,
		aggregateStatsWriter:                    aggregateStatsWriter,
	}
}

func (d *DataExporter) Export(county string, configurableEligibilityFlow data.ConfigurableEligibilityFlow) Summary {
	for i, row := range d.dojInformation.Rows {
		d.outputDOJWriter.WriteEntryWithEligibilityInfo(row, d.normalFlowEligibilities[i])
		d.outputCondensedDOJWriter.WriteCondensedEntryWithEligibilityInfo(row, d.normalFlowEligibilities[i])
		if d.normalFlowEligibilities[i] != nil {
			d.outputProp64ConvictionsDOJWriter.WriteEntryWithEligibilityInfo(row, d.normalFlowEligibilities[i])
		}
	}

	d.outputDOJWriter.Flush()
	d.outputCondensedDOJWriter.Flush()
	d.outputProp64ConvictionsDOJWriter.Flush()
	d.PrintAggregateStatistics(county)
	return d.NewSummary(county, configurableEligibilityFlow)
}

func (d *DataExporter) PrintAggregateStatistics(county string) {
	fmt.Fprintf(d.aggregateStatsWriter, "----------- Overall summary of DOJ file --------------------\n")
	fmt.Fprintf(d.aggregateStatsWriter, "Found %d Total rows in DOJ file\n", d.dojInformation.TotalRows())
	fmt.Fprintf(d.aggregateStatsWriter, "Found %d Total individuals in DOJ file\n", d.dojInformation.TotalIndividuals())
	fmt.Fprintf(d.aggregateStatsWriter, "Found %d Total convictions in DOJ file\n", d.dojInformation.TotalConvictions())
	fmt.Fprintf(d.aggregateStatsWriter, "Found %d convictions in this county\n", d.dojInformation.TotalConvictionsInCounty(county))

	fmt.Fprintf(d.aggregateStatsWriter, "\n")
	fmt.Fprintf(d.aggregateStatsWriter, "----------- Prop64 Convictions Overall--------------------")
	d.printSummaryByCodeSection("total", d.dojInformation.OverallProp64ConvictionsByCodeSection())
	fmt.Fprintf(d.aggregateStatsWriter, "\n")

	fmt.Fprintf(d.aggregateStatsWriter, "----------- Prop64 Convictions In This County --------------------")
	d.printSummaryByCodeSection("in this county", d.dojInformation.Prop64ConvictionsInThisCountyByCodeSection(county))
	fmt.Fprintf(d.aggregateStatsWriter, "Date of earliest Prop 64 conviction: %s", d.dojInformation.EarliestProp64ConvictionDateInThisCounty(county).Format("January 2006"))
	d.printSummaryByCodeSectionByEligibility(d.dojInformation.Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county, d.normalFlowEligibilities))

	fmt.Fprintf(d.aggregateStatsWriter, "\n")
	fmt.Fprintf(d.aggregateStatsWriter, "----------- Eligibility Reasons --------------------\n")
	d.printSummaryByEligibilityByReason(d.dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county, d.normalFlowEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "\n\n")
	fmt.Fprintf(d.aggregateStatsWriter, "----------- Prop64 Related Convictions In This County --------------------")
	d.printSummaryByCodeSection("in this county", d.dojInformation.OverallRelatedConvictionsByCodeSection())
	fmt.Fprintf(d.aggregateStatsWriter, "\n")
	d.printSummaryByCodeSectionByEligibility(d.dojInformation.RelatedConvictionsInThisCountyByCodeSectionByEligibility(county, d.normalFlowEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "\n")

	fmt.Fprintf(d.aggregateStatsWriter, "----------- Impact to individuals --------------------\n")
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals currently have a felony on their record\n", d.dojInformation.CountIndividualsWithFelony())
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals currently have convictions on their record\n", d.dojInformation.CountIndividualsWithConviction())
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals currently have convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsWithConvictionInLast7Years())
	fmt.Fprintf(d.aggregateStatsWriter, "\n")

	fmt.Fprintf(d.aggregateStatsWriter, "----------- Eligibility is run as specified for Prop 64 and Related Charges --------------------\n")
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had a felony will no longer have a felony on their record\n", d.dojInformation.CountIndividualsNoLongerHaveFelony(d.normalFlowEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had convictions will no longer have any convictions on their record\n", d.dojInformation.CountIndividualsNoLongerHaveConviction(d.normalFlowEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years(d.normalFlowEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "\n")
	fmt.Fprintf(d.aggregateStatsWriter, "----------- If ALL Prop 64 convictions are dismissed and sealed --------------------\n")
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had a felony will no longer have a felony on their record\n", d.dojInformation.CountIndividualsNoLongerHaveFelony(d.dismissAllProp64Eligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had convictions will no longer have any convictions on their record\n", d.dojInformation.CountIndividualsNoLongerHaveConviction(d.dismissAllProp64Eligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years(d.dismissAllProp64Eligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "\n")
	fmt.Fprintf(d.aggregateStatsWriter, "----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------\n")
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had a felony will no longer have a felony on their record\n", d.dojInformation.CountIndividualsNoLongerHaveFelony(d.dismissAllProp64AndRelatedEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had convictions will no longer have any convictions on their record\n", d.dojInformation.CountIndividualsNoLongerHaveConviction(d.dismissAllProp64AndRelatedEligibilities))
	fmt.Fprintf(d.aggregateStatsWriter, "%d individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years\n", d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years(d.dismissAllProp64AndRelatedEligibilities))
}

func (d *DataExporter) printSummaryByCodeSection(description string, resultsByCodeSection map[string]int) {
	fmt.Fprintf(d.aggregateStatsWriter, "\nFound %d convictions %s\n", sumValues(resultsByCodeSection), description)
	formatString := fmt.Sprintf("Found %%d %%s convictions %s\n", description)
	d.printMap(formatString, resultsByCodeSection)
}

func (d *DataExporter) printSummaryByCodeSectionByEligibility(resultsByCodeSectionByEligibility map[string]map[string]int) {
	codeSections := make([]string, 0, len(resultsByCodeSectionByEligibility))
	for key := range resultsByCodeSectionByEligibility {
		codeSections = append(codeSections, key)
	}
	sort.Strings(codeSections)

	for _, codeSection := range codeSections {
		fmt.Fprintf(d.aggregateStatsWriter, "\n%v", codeSection)

		total := 0
		eligibilityMap := resultsByCodeSectionByEligibility[codeSection]
		eligibilities := make([]string, 0, len(eligibilityMap))
		for key := range eligibilityMap {
			eligibilities = append(eligibilities, key)
		}
		sort.Strings(eligibilities)

		for _, eligibility := range eligibilities {

			total += eligibilityMap[eligibility]
			fmt.Fprintf(d.aggregateStatsWriter, "\nFound %v %v convictions that are %v", eligibilityMap[eligibility], eligibility, codeSection)
		}
		fmt.Fprintf(d.aggregateStatsWriter, "\nFound %v convictions total that are %v\n", total, codeSection)
	}
}

func (d *DataExporter) printSummaryByEligibilityByReason(resultsByEligibilityByReason map[string]map[string]int) {
	determinations := make([]string, 0, len(resultsByEligibilityByReason))
	for key := range resultsByEligibilityByReason {
		determinations = append(determinations, key)
	}
	sort.Strings(determinations)

	for _, determination := range determinations {
		fmt.Fprintf(d.aggregateStatsWriter, "\n%v", determination)

		total := 0
		reasonMap := resultsByEligibilityByReason[determination]
		reasons := make([]string, 0, len(reasonMap))
		for key := range reasonMap {
			reasons = append(reasons, key)
		}
		sort.Strings(reasons)

		for _, reason := range reasons {
			total += reasonMap[reason]
			fmt.Fprintf(d.aggregateStatsWriter, "\nFound %v convictions with eligibility reason %v", reasonMap[reason], reason)
		}
		fmt.Fprintf(d.aggregateStatsWriter, "\n")
	}
}

func (d *DataExporter) printMap(formatString string, values map[string]int) {
	keys := getSortedKeys(values)

	for _, key := range keys {
		fmt.Fprintf(d.aggregateStatsWriter, formatString, values[key], key)
	}
}

func (d *DataExporter) AccumulateSummaryData(runSummary Summary, fileSummary Summary) Summary {
	return Summary{
		County:                              fileSummary.County,
		LineCount:                           runSummary.LineCount + fileSummary.LineCount,
		EarliestConviction:                  findEarliest(runSummary.EarliestConviction, fileSummary.EarliestConviction),
		ReliefWithCurrentEligibilityChoices: utilities.AddMaps(runSummary.ReliefWithCurrentEligibilityChoices, fileSummary.ReliefWithCurrentEligibilityChoices),
		ReliefWithDismissAllProp64:          utilities.AddMaps(runSummary.ReliefWithDismissAllProp64, fileSummary.ReliefWithDismissAllProp64),
		Prop64ConvictionsCountInCountyByCodeSection: utilities.AddMaps(runSummary.Prop64ConvictionsCountInCountyByCodeSection, fileSummary.Prop64ConvictionsCountInCountyByCodeSection),
		Prop64FelonyConvictionsCountInCounty:        runSummary.Prop64FelonyConvictionsCountInCounty + fileSummary.Prop64FelonyConvictionsCountInCounty,
		Prop64NonFelonyConvictionsCountInCounty:     runSummary.Prop64NonFelonyConvictionsCountInCounty + fileSummary.Prop64NonFelonyConvictionsCountInCounty,
		SubjectsWithSomeReliefCount:                 runSummary.SubjectsWithSomeReliefCount + fileSummary.SubjectsWithSomeReliefCount,
		ConvictionDismissalCountByAdditionalRelief:  utilities.AddMaps(runSummary.ConvictionDismissalCountByAdditionalRelief, fileSummary.ConvictionDismissalCountByAdditionalRelief),
		ConvictionDismissalCountByCodeSection:       utilities.AddMaps(runSummary.ConvictionDismissalCountByCodeSection, fileSummary.ConvictionDismissalCountByCodeSection),
		ConvictionReductionCountByCodeSection:       utilities.AddMaps(runSummary.ConvictionReductionCountByCodeSection, fileSummary.ConvictionReductionCountByCodeSection),
		SubjectsWithProp64ConvictionCountInCounty:   runSummary.SubjectsWithProp64ConvictionCountInCounty + fileSummary.SubjectsWithProp64ConvictionCountInCounty,
	}
}

func (d *DataExporter) NewSummary(county string, configurableEligibilityFlow data.ConfigurableEligibilityFlow) Summary {
	return Summary{
		County:             county,
		LineCount:          d.dojInformation.TotalRows(),
		EarliestConviction: d.dojInformation.EarliestProp64ConvictionDateInThisCounty(county),
		ReliefWithCurrentEligibilityChoices: map[string]int{
			"CountSubjectsNoFelony":               d.dojInformation.CountIndividualsNoLongerHaveFelony(d.normalFlowEligibilities),
			"CountSubjectsNoConvictionLast7Years": d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years(d.normalFlowEligibilities),
			"CountSubjectsNoConviction":           d.dojInformation.CountIndividualsNoLongerHaveConviction(d.normalFlowEligibilities),
		},
		ReliefWithDismissAllProp64: map[string]int{
			"CountSubjectsNoFelony":               d.dojInformation.CountIndividualsNoLongerHaveFelony(d.dismissAllProp64Eligibilities),
			"CountSubjectsNoConvictionLast7Years": d.dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years(d.dismissAllProp64Eligibilities),
			"CountSubjectsNoConviction":           d.dojInformation.CountIndividualsNoLongerHaveConviction(d.dismissAllProp64Eligibilities),
		},
		Prop64ConvictionsCountInCountyByCodeSection: d.dojInformation.Prop64ConvictionsInThisCountyByCodeSection(county),
		ConvictionDismissalCountByCodeSection:       d.getDismissalsByCodeSection(county, configurableEligibilityFlow),
		ConvictionReductionCountByCodeSection:       d.getReductionsByCodeSection(county, configurableEligibilityFlow),
		ConvictionDismissalCountByAdditionalRelief:  d.getDismissalsByAdditionalRelief(county, configurableEligibilityFlow),
		SubjectsWithSomeReliefCount:                 d.dojInformation.CountIndividualsWithSomeRelief(d.normalFlowEligibilities),
		Prop64FelonyConvictionsCountInCounty:        d.dojInformation.TotalConvictionsInCountyFiltered(county, data.IsFelonyFilter, matchers.IsProp64Charge),
		Prop64NonFelonyConvictionsCountInCounty:     d.dojInformation.TotalConvictionsInCountyFiltered(county, data.IsNotFelonyFilter, matchers.IsProp64Charge),
		SubjectsWithProp64ConvictionCountInCounty:   d.dojInformation.CountIndividualsWithProp64ConvictionInCounty(county),
	}
}

func findEarliest(time1 time.Time, time2 time.Time) time.Time {
	if !time1.IsZero() && time1.Before(time2) {
		return time1
	}
	return time2
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

func (d *DataExporter) getDismissalsByCodeSection(county string, configurableEligibilityFlow data.ConfigurableEligibilityFlow) map[string]int {
	result := make(map[string]int)
	var eligibilityReasonKey string
	for _, codeSection := range configurableEligibilityFlow.DismissSections {
		eligibilityReasonKey = fmt.Sprintf("Dismiss all HS %s convictions", codeSection)
		result[codeSection] = d.dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county, d.normalFlowEligibilities)["Eligible for Dismissal"][eligibilityReasonKey]
	}
	return result
}

func (d *DataExporter) getReductionsByCodeSection(county string, configurableEligibilityFlow data.ConfigurableEligibilityFlow) map[string]int {
	result := make(map[string]int)
	var eligibilityReasonKey string
	for _, codeSection := range configurableEligibilityFlow.ReduceSections {
		eligibilityReasonKey = fmt.Sprintf("Reduce all HS %s convictions", codeSection)
		result[codeSection] = d.dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county, d.normalFlowEligibilities)["Eligible for Reduction"][eligibilityReasonKey]
	}
	return result
}

func (d *DataExporter) getDismissalsByAdditionalRelief(county string, configurableEligibilityFlow data.ConfigurableEligibilityFlow) map[string]int {
	result := make(map[string]int)
	for key, value := range d.dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county, d.normalFlowEligibilities)["Eligible for Dismissal"] {
		if !strings.HasPrefix(key, "Dismiss all HS") && !strings.HasPrefix(key, "Misdemeanor or Infraction"){
			result[key] = value
		}
	}
	return result
}
