package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gogen/data"
	"gogen/exporter"
	"gogen/test_fixtures"
	"gogen/utilities"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.2.2"

var defaultOpts struct{}

type runOpts struct {
	OutputFolder       string `long:"outputs" description:"The folder in which to place result files"`
	DOJFile            string `long:"input-doj" description:"The file containing criminal histories from CA DOJ"`
	County             string `long:"county" short:"c" description:"The county for which eligibility will be computed"`
	ComputeAt          string `long:"compute-at" description:"The date for which eligibility will be evaluated, ex: 2020-10-31"`
	EligibilityOptions string `long:"eligibility-options" description:"File containing options for which eligibility logic to apply"`
	FileNameSuffix     string `long:"file-name-suffix" hidden:"true" description:"string to append to file names"`
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

	utilities.SetErrorFileName(utilities.GenerateFileName(r.OutputFolder, "gogen%s.err", r.FileNameSuffix))

	if r.OutputFolder == "" || r.DOJFile == "" || r.County == "" {
		utilities.ExitWithError(errors.New("missing required field: Run gogen --help for more info"), utilities.INVALID_OPTION_ERROR)
	}

	computeAtDate := time.Now()

	if r.ComputeAt != "" {
		computeAtOption, err := time.Parse("2006-01-02", r.ComputeAt)
		if err != nil {
			utilities.ExitWithError(errors.New("invalid --compute-at date: Must be a valid date in the format YYYY-MM-DD"), utilities.INVALID_OPTION_ERROR)
		} else {
			computeAtDate = computeAtOption
		}
	}

	var countyEligibilityFlow data.EligibilityFlow

	if r.EligibilityOptions != "" {
		var options data.EligibilityOptions
		optionsFile, err := os.Open(r.EligibilityOptions)
		if err != nil {
			utilities.ExitWithError(err, utilities.INVALID_OPTION_ERROR)
		}
		defer optionsFile.Close()

		optionsBytes, err := ioutil.ReadAll(optionsFile)
		if err != nil {
			utilities.ExitWithError(err, utilities.INVALID_OPTION_ERROR)
		}

		err = json.Unmarshal(optionsBytes, &options)
		if err != nil {
			utilities.ExitWithError(err, utilities.INVALID_OPTION_ERROR)
		}
		countyEligibilityFlow = data.NewConfigurableEligibilityFlow(options, r.County)
	} else {
		countyEligibilityFlow = data.EligibilityFlows[r.County]
	}

	dojInformation := data.NewDOJInformation(r.DOJFile, computeAtDate, countyEligibilityFlow)
	countyEligibilities := dojInformation.DetermineEligibility(r.County, countyEligibilityFlow)

	dismissAllProp64Eligibilities := dojInformation.DetermineEligibility(r.County, data.EligibilityFlows["DISMISS ALL PROP 64"])
	dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility(r.County, data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

	dojFilePath := utilities.GenerateFileName(r.OutputFolder, "doj_results%s.csv", r.FileNameSuffix)
	condensedFilePath := utilities.GenerateFileName(r.OutputFolder, "doj_results_condensed%s.csv", r.FileNameSuffix)
	prop64ConvictionsFilePath := utilities.GenerateFileName(r.OutputFolder, "doj_results_convictions%s.csv", r.FileNameSuffix)
	outputFilePath := utilities.GenerateFileName(r.OutputFolder, "gogen%s.out", r.FileNameSuffix)

	dojWriter := exporter.NewDOJWriter(dojFilePath)
	condensedDojWriter := exporter.NewCondensedDOJWriter(condensedFilePath)
	prop64ConvictionsDojWriter := exporter.NewDOJWriter(prop64ConvictionsFilePath)
	outputWriter := utilities.GetOutputWriter(outputFilePath)

	dataExporter := exporter.NewDataExporter(dojInformation, countyEligibilities, dismissAllProp64Eligibilities, dismissAllProp64AndRelatedEligibilities, dojWriter, condensedDojWriter, prop64ConvictionsDojWriter, outputWriter)

	dataExporter.Export(r.County, processingStartTime)

	return nil
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
