package main

import (
	"fmt"
	. "gogen/data"
	. "gogen/processor"
	"os"
	"path/filepath"
	"time"

	"github.com/jessevdk/go-flags"
)

const VERSION = "0.0.1"

var defaultOpts struct {}

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files"`
	DOJFile      string `long:"input-doj" description:"The file containing criminal histories from CA DOJ"`
	Delimiter    string `long:"delimiter" short:"d" default:"," hidden:"true"`
	County       string `long:"county" short:"c" description:"The county for which eligibility will be computed"`
	Version      bool `long:"version" short:"v" description:"Print the version"`
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

	dojInformation, _ := NewDOJInformation(opts.DOJFile, time.Now(), opts.County)

	dojWriter := NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results.csv"))
	condensedDojWriter := NewCondensedDOJWriter(filepath.Join(opts.OutputFolder, "doj_results_condensed.csv"))

	dataProcessor := NewDataProcessor(dojInformation, dojWriter, condensedDojWriter)

	dataProcessor.Process(opts.County)
}
