package exporter

import (
	"fmt"
	"gogen/data"
	"gogen/matchers"
	"gogen/utilities"
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
) DataExporter {

	return DataExporter{
		dojInformation:                          dojInformation,
		normalFlowEligibilities:                 countyEligibilities,
		dismissAllProp64Eligibilities:           dismissAllProp64Eligibilities,
		dismissAllProp64AndRelatedEligibilities: dismissAllProp64AndRelatedEligibilities,
		outputDOJWriter:                         outputDOJWriter,
		outputCondensedDOJWriter:                outputCondensedDOJWriter,
		outputProp64ConvictionsDOJWriter:        outputProp64ConvictionsDOJWriter,
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
	return d.NewSummary(county, configurableEligibilityFlow)
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
