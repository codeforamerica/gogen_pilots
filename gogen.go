package main

import (
	"encoding/csv"
	. "gogen/data"
	. "gogen/processor"
	"os"
	"path/filepath"
	"time"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	DOJFile      string `long:"input-doj" description:"The file containing criminal histories from CA DOJ" required:"true"`
	Delimiter    string `long:"delimiter" short:"d" default:"," hidden:"true"`
	County       string `long:"county" short:"c" description:"The county for which eligibility will be computed" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	dojFile, err := os.Open(opts.DOJFile)
	if err != nil {
		panic(err)
	}

	dojInformation, _ := NewDOJInformation(csv.NewReader(dojFile), time.Now(), opts.County)

	dojWriter := NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results.csv"))

	dataProcessor := NewDataProcessor(dojInformation, dojWriter)

	dataProcessor.Process(opts.County)
}
