package main

import (
	"encoding/csv"

	. "gogen/data"
	. "gogen/processor"

	"os"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	OutputFolder      string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	ConvictionWeights string `long:"conviction-weights" description:"The file containing conviction weights from SFDA" required:"true"`
	DOJFile           string `long:"input-doj" description:"The file containing criminal histories from CA DOJ" required:"true"`
	CMSFile           string `long:"input-csv" description:"The file containing criminal histories from SF's cms" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	cmsFile, err := os.Open(opts.CMSFile)
	if err != nil {
		panic(err)
	}

	cmsCSV := csv.NewReader(cmsFile)

	weightsFile, err := os.Open(opts.ConvictionWeights)
	if err != nil {
		panic(err)
	}

	weightsInformation, _ := NewWeightsInformation(csv.NewReader(weightsFile))

	dojFile, err := os.Open(opts.DOJFile)
	if err != nil {
		panic(err)
	}

	dojInformation, _ := NewDOJInformation(csv.NewReader(dojFile))

	dataProcessor := NewDataProcessor(cmsCSV, weightsInformation, dojInformation, opts.OutputFolder)

	dataProcessor.Process()
}
