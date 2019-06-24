package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gogen/data"
	"gogen/exporter"
	"gogen/test_fixtures"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.0.3"

var defaultOpts struct{}

type runOpts struct {
	OutputFolder       string `long:"outputs" description:"The folder in which to place result files"`
	DOJFile            string `long:"input-doj" description:"The file containing criminal histories from CA DOJ"`
	County             string `long:"county" short:"c" description:"The county for which eligibility will be computed"`
	ComputeAt          string `long:"compute-at" description:"The date for which eligibility will be evaluated, ex: 2020-10-31"`
	EligibilityOptions string `long:"eligibility-options" description:"File containing options for which eligibility logic to apply"`
	DateForFileName    string `long:"date-for-file-name" short:"date" hidden:"true" description:"Date to append to file names for clarity"`
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
	if r.OutputFolder == "" || r.DOJFile == "" || r.County == "" {
		panic("Missing required field! Run gogen --help for more info.")
	}

	computeAtDate := time.Now()

	if r.ComputeAt != "" {
		computeAtOption, err := time.Parse("2006-01-02", r.ComputeAt)
		if err != nil {
			panic("Invalid --compute-at date. Must be a valid date of the format YYYY-MM-DD.")
		} else {
			computeAtDate = computeAtOption
		}
	}

	var countyEligibilityFlow data.EligibilityFlow

	if r.EligibilityOptions != "" {
		var options data.EligibilityOptions
		optionsFile, err := os.Open(r.EligibilityOptions)
		if err != nil {
			panic(err)
		}
		defer optionsFile.Close()

		optionsBytes, err := ioutil.ReadAll(optionsFile)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(optionsBytes, &options)
		if err != nil {
			panic(err)
		}
		countyEligibilityFlow = data.NewConfigurableEligibilityFlow(options, r.County)
	} else {
		countyEligibilityFlow = data.EligibilityFlows[r.County]
	}

	countyDojInformation := data.NewDOJInformation(r.DOJFile, computeAtDate, r.County, countyEligibilityFlow)

	dismissAllProp64EligibilityFlow := data.EligibilityFlows["DISMISS ALL PROP 64"]
	dismissAllProp64AndRelatedEligibilityFlow := data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"]

	dismissAllProp64DojInformation := data.NewDOJInformation(r.DOJFile, computeAtDate, r.County, dismissAllProp64EligibilityFlow)
	dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(r.DOJFile, computeAtDate, r.County, dismissAllProp64AndRelatedEligibilityFlow)

	resultsStrings := []string{"doj", "results", r.DateForFileName}
	dojResultsFileName := strings.Join(resultsStrings, "_")

	condensedStrings := []string{"doj", "results", "condensed", r.DateForFileName}
	dojResultsCondensedFileName := strings.Join(condensedStrings, "_")

	convictionsStrings := []string{"doj", "results", "convictions", r.DateForFileName}
	dojResultsConvictionsFileName := strings.Join(convictionsStrings, "_")

	dojWriter := exporter.NewDOJWriter(filepath.Join(r.OutputFolder, dojResultsFileName))
	condensedDojWriter := exporter.NewCondensedDOJWriter(filepath.Join(r.OutputFolder, dojResultsCondensedFileName))
	prop64ConvictionsDojWriter := exporter.NewDOJWriter(filepath.Join(r.OutputFolder, dojResultsConvictionsFileName))

	dataExporter := exporter.NewDataExporter(countyDojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, condensedDojWriter, prop64ConvictionsDojWriter)

	dataExporter.Export(r.County)

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
