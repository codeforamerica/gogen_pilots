package main

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	OutputFolder      string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	ConvictionWeights string `long:"conviction-weights" description:"The file containing conviction weights from SFDA" required:"true"`
	DOJFile           string `long:"input-doj" description:"The file containing criminal histories from CA DOJ" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	fmt.Println(opts.OutputFolder)
}
