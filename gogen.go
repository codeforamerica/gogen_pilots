package main

import (
	"gogen/internal/data"
	"gogen/internal/processor"
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

	dojInformation, _ := data.NewDOJInformation(opts.DOJFile, time.Now(), opts.County)

	dojWriter := processor.NewDOJWriter(filepath.Join(opts.OutputFolder, "doj_results.csv"))
	condensedDojWriter := processor.NewCondensedDOJWriter(filepath.Join(opts.OutputFolder, "doj_results_condensed.csv"))

	dataProcessor := processor.NewDataProcessor(dojInformation, dojWriter, condensedDojWriter)

	dataProcessor.Process(opts.County)
}
