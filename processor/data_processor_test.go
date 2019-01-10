package processor_test

import (
	"encoding/csv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"gogen/data"
	. "gogen/processor"
	"io/ioutil"
	"os"
	path "path/filepath"
	"time"
)

var _ = Describe("DataProcessor", func() {
	var (
		outputDir     string
		dataProcessor DataProcessor
		err           error
	)

	BeforeEach(func() {
		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "cadoj.csv"))
		Expect(err).ToNot(HaveOccurred())

		dojFile, err := os.Open(pathToDOJ)
		if err != nil {
			panic(err)
		}

		comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

		dojInformation, _ := data.NewDOJInformation(csv.NewReader(dojFile), comparisonTime)

		dojWriter := NewDOJWriter(path.Join(outputDir, "doj_results.csv"))

		dataProcessor = NewDataProcessor(dojInformation, dojWriter)
	})

	It("runs and has output", func() {
		dataProcessor.Process()
		format.TruncatedDiff = false

		pathToDOJOutput, err := path.Abs(path.Join(outputDir, "doj_results.csv"))
		Expect(err).ToNot(HaveOccurred())
		outputDOJBody, err := ioutil.ReadFile(pathToDOJOutput)
		Expect(err).ToNot(HaveOccurred())

		pathToExpectedDOJResults, err := path.Abs(path.Join("..", "test_fixtures", "doj_results.csv"))
		Expect(err).ToNot(HaveOccurred())
		expectedDOJResultsBody, err := ioutil.ReadFile(pathToExpectedDOJResults)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(outputDOJBody)).To(Equal(string(expectedDOJResultsBody)))
	})
})
