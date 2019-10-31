package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gogen_pilots/data"
	"gogen_pilots/exporter"
	"gogen_pilots/test_fixtures"
	"gogen_pilots/utilities"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.2.9"

var defaultOpts struct{}

type runOpts struct {
	OutputFolder   string  `long:"outputs" description:"The folder in which to place result files"`
	DOJFiles       string  `long:"input-doj" description:"The files containing criminal histories from CA DOJ"`
	ComputeAt      string  `long:"compute-at" description:"The date for which eligibility will be evaluated, ex: 2020-10-31"`
	FileNameSuffix string  `long:"file-name-suffix" hidden:"true" description:"string to append to file names"`
	IndividualAge  int `long:"individual-age" hidden:"true" description:"minimum age of individual for record clearance"`
	YearsConvictionFree  int `long:"years-conviction-free" hidden:"true" description:"years (as a number) since last conviction"`
}

type exportTestCSVOpts struct {
	ExcelFixturePath string `long:"excel-fixture-path" short:"e" description:"Path to a county's excel fixture file to generate test CSVs"`
	OutputFolder     string `long:"outputs" short:"o" description:"The folder in which to place result files"`
}

type versionOpts struct{}

var opts struct {
	Version   versionOpts       `command:"version" description:"Print the version"`
	Run       runOpts           `command:"run" description:"Process an input DOJ file and produce an annotated DOJ data file"`
	ExportCSV exportTestCSVOpts `command:"export-test-csv" description:"Export example data files from excel fixtures"`
}

func (r runOpts) Execute(args []string) error {

	var processingStartTime time.Time
	processingStartTime = time.Now()

	utilities.SetErrorFileName(utilities.GenerateFileName(r.OutputFolder, "gogen_pilots%s.err", r.FileNameSuffix))

	if r.OutputFolder == "" || r.DOJFiles == "" {
		utilities.ExitWithError(errors.New("missing required field: Run gogen_pilots --help for more info"), utilities.INVALID_RUN_OPTION_ERROR)
	}

	inputFiles := strings.Split(r.DOJFiles, ",")

	computeAtDate := time.Now()

	if r.ComputeAt != "" {
		computeAtOption, err := time.Parse("2006-01-02", r.ComputeAt)
		if err != nil {
			utilities.ExitWithError(errors.New("invalid --compute-at date: Must be a valid date in the format YYYY-MM-DD"), utilities.INVALID_RUN_OPTION_ERROR)
		} else {
			computeAtDate = computeAtOption
		}
	}
	var age int

	if r.IndividualAge != 0 {
		age = r.IndividualAge
	} else {
		age = 50
	}
	var yearsConvictionFree int
	if r.YearsConvictionFree != 0 {
		yearsConvictionFree = r.YearsConvictionFree
	} else {
		yearsConvictionFree = 10
	}

	var countyEligibilityFlow data.EligibilityFlow

	countyEligibilityFlow = data.EligibilityFlows["LOS ANGELES"]

	var runErrors []error
	runSummary := exporter.Summary{
		County: "LOS ANGELES",
		IndividualDismissAge:age,
		YearsConvictionFree: yearsConvictionFree,
	}
	outputJsonFilePath := utilities.GenerateFileName(r.OutputFolder, "gogen_pilots%s.json", r.FileNameSuffix)

	for fileIndex, inputFile := range inputFiles {
		fileIndex = fileIndex + 1
		fileOutputFolder := utilities.GenerateIndexedOutputFolder(r.OutputFolder, fileIndex, r.FileNameSuffix)
		err := os.MkdirAll(fileOutputFolder, os.ModePerm)
		if err != nil {
			runErrors = append(runErrors, err)
			continue
		}
		dojInformation, err := data.NewDOJInformation(inputFile, computeAtDate, countyEligibilityFlow)
		if err != nil {
			runErrors = append(runErrors, err)
			continue
		}
		countyEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", countyEligibilityFlow, age, yearsConvictionFree)

		dismissAllProp64Eligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"], age, yearsConvictionFree)
		dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"], age, yearsConvictionFree)

		dojFilePath := utilities.GenerateIndexedFileName(fileOutputFolder, "doj_results_%d%s.csv", fileIndex, r.FileNameSuffix)
		condensedFilePath := utilities.GenerateIndexedFileName(fileOutputFolder, "doj_results_condensed_%d%s.csv", fileIndex, r.FileNameSuffix)
		prop64ConvictionsFilePath := utilities.GenerateIndexedFileName(fileOutputFolder, "doj_results_convictions_%d%s.csv", fileIndex, r.FileNameSuffix)
		outputFilePath := utilities.GenerateIndexedFileName(fileOutputFolder, "gogen_pilots_%d%s.out", fileIndex, r.FileNameSuffix)

		dojWriter, err := exporter.NewDOJWriter(dojFilePath)
		if err != nil {
			runErrors = append(runErrors, err)
			continue
		}
		condensedDojWriter, err := exporter.NewCondensedDOJWriter(condensedFilePath)
		if err != nil {
			runErrors = append(runErrors, err)
			continue
		}
		prop64ConvictionsDojWriter, err := exporter.NewDOJWriter(prop64ConvictionsFilePath)
		if err != nil {
			runErrors = append(runErrors, err)
			continue
		}
		aggregateFileStatsWriter := utilities.GetOutputWriter(outputFilePath)

		dataExporter := exporter.NewDataExporter(
			dojInformation,
			countyEligibilities,
			dismissAllProp64Eligibilities,
			dismissAllProp64AndRelatedEligibilities,
			dojWriter,
			condensedDojWriter,
			prop64ConvictionsDojWriter,
			aggregateFileStatsWriter)

		fileSummary := dataExporter.Export("LOS ANGELES", processingStartTime)
		runSummary = dataExporter.AccumulateSummaryData(runSummary, fileSummary)
	}

	if len(runErrors) > 0 {
		utilities.ExitWithErrors(runErrors, utilities.FILE_PROCESSING_ERROR)
	}

	ExportSummary(runSummary, processingStartTime, outputJsonFilePath)
	return nil
}

func ExportSummary(summary exporter.Summary, startTime time.Time, filePath string) {
	summary.ProcessingTimeInSeconds = time.Since(startTime).Seconds()

	summary.GitRef = "no_superstrikes"
	s, err := json.Marshal(summary)
	if err != nil {
		panic("Cannot marshal JSON") // TODO replace panic
	}
	err = ioutil.WriteFile(filePath, s, 0644)
	if err != nil {
		panic("Cannot write JSON") // TODO replace panic
	}
}

func (e exportTestCSVOpts) Execute(args []string) error {
	if e.ExcelFixturePath != "" {
		inputCSV, expectedResultsCSV, err := test_fixtures.ExportFullCSVFixtures(e.ExcelFixturePath, e.OutputFolder)
		if err != nil {
			fmt.Println("Extracting test CSVs failed")
			os.Exit(1)
		}

		fmt.Println("Wrote input CSV at: " + inputCSV)
		fmt.Println("Wrote expected results CSV at: " + expectedResultsCSV)

		openCommand := exec.Command("open", e.OutputFolder)
		err = openCommand.Run()
		if err != nil {
			panic(err)
		}
	} else {
		return errors.New("something went wrong")
	}

	return nil
}

func (v versionOpts) Execute(args []string) error {
	fmt.Println(VERSION)
	return nil
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
