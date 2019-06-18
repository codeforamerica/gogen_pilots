package main

import (
	"encoding/json"
	"fmt"
	"gogen/data"
	"gogen/exporter"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.0.2"

var defaultOpts struct{}

var opts struct {
	OutputFolder       string `long:"outputs" description:"The folder in which to place result files"`
	DOJFile            string `long:"input-doj" description:"The file containing criminal histories from CA DOJ"`
	County             string `long:"county" short:"c" description:"The county for which eligibility will be computed"`
	Version            bool   `long:"version" short:"v" description:"Print the version"`
	ComputeAt          string `long:"compute-at" description:"The date for which eligibility will be evaluated, ex: 2020-10-31"`
	EligibilityOptions string `long:"eligibility-options" description:"File containing options for which eligibility logic to apply"`
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

	if opts.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if opts.OutputFolder == "" || opts.DOJFile == "" || opts.County == "" {
		panic("Missing required field! Run gogen --help for more info.")
	}

	computeAtDate := time.Now()

	if opts.ComputeAt != "" {
		computeAtOption, err := time.Parse("2006-01-02", opts.ComputeAt)
		if err != nil {
			panic("Invalid --compute-at date. Must be a valid date of the format YYYY-MM-DD.")
		} else {
			computeAtDate = computeAtOption
		}
	}

	var countyEligibilityFlow data.EligibilityFlow

	if opts.EligibilityOptions != "" {
		var options data.EligibilityOptions
		optionsFile, err := os.Open(opts.EligibilityOptions)
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
		countyEligibilityFlow = data.NewConfigurableEligibilityFlow(options, opts.County)
	} else {
		countyEligibilityFlow = data.EligibilityFlows[opts.County]
	}

	countyDojInformation := data.NewDOJInformation(opts.DOJFile, computeAtDate, opts.County, countyEligibilityFlow)

	dismissAllProp64EligibilityFlow := data.EligibilityFlows["DISMISS ALL PROP 64"]
	dismissAllProp64AndRelatedEligibilityFlow := data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"]

	dismissAllProp64DojInformation := data.NewDOJInformation(opts.DOJFile, computeAtDate, opts.County, dismissAllProp64EligibilityFlow)
	dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(opts.DOJFile, computeAtDate, opts.County, dismissAllProp64AndRelatedEligibilityFlow)

	dojWriter := exporter.NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results.csv"))
	condensedDojWriter := exporter.NewCondensedDOJWriter(filepath.Join(opts.OutputFolder, "doj_results_condensed.csv"))
	prop64ConvictionsDojWriter := exporter.NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results_convictions.csv"))

	dataExporter := exporter.NewDataExporter(countyDojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, condensedDojWriter, prop64ConvictionsDojWriter)

	dataExporter.Export(opts.County)
}
