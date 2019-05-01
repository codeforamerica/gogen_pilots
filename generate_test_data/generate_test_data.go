package main

import (
	"github.com/jessevdk/go-flags"
	. "gogen/processor"
	"path/filepath"
)

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	TargetSize   int    `long:"target-size" description:"Desired number of lines in the output file" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	testWriter := NewWriter(filepath.Join(opts.OutputFolder, "generated_test_data.csv"), DojFullHeaders)
	line := []string{"2", "3", "5", "7", "11", "13"}

	for i := 0; i < opts.TargetSize; i++ {
		testWriter.Write(line)
	}

	testWriter.Flush()
}
