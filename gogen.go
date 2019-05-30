package main

import (
	"fmt"
	"gogen/data"
	"gogen/processor"
	"os"
	"path/filepath"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.0.2"

var defaultOpts struct{}

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files"`
	DOJFile      string `long:"input-doj" description:"The file containing criminal histories from CA DOJ"`
	County       string `long:"county" short:"c" description:"The county for which eligibility will be computed"`
	Version      bool   `long:"version" short:"v" description:"Print the version"`
	ComputeAt    string `long:"compute-at" description:"The date for which eligibility will be evaluated, ex: 2020-10-31"`
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

	dojInformation, _ := data.NewDOJInformation(opts.DOJFile, computeAtDate, opts.County)

	dojWriter := processor.NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results.csv"))
	condensedDojWriter := processor.NewCondensedDOJWriter(filepath.Join(opts.OutputFolder, "doj_results_condensed.csv"))
	prop64ConvictionsDojWriter := processor.NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results_convictions.csv"))

	dataProcessor := processor.NewDataProcessor(dojInformation, dojWriter, condensedDojWriter, prop64ConvictionsDojWriter)

	dataProcessor.Process(opts.County)
}
