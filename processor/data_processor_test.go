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

		pathToWeights, err := path.Abs(path.Join("..", "test_fixtures", "conviction_weights.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "cadoj.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToCMS, err := path.Abs(path.Join("..", "test_fixtures", "felonies_sf.csv"))
		Expect(err).ToNot(HaveOccurred())

		cmsFile, err := os.Open(pathToCMS)
		if err != nil {
			panic(err)
		}

		cmsCSV := csv.NewReader(cmsFile)

		weightsFile, err := os.Open(pathToWeights)
		if err != nil {
			panic(err)
		}

		weightsInformation, _ := data.NewWeightsInformation(csv.NewReader(weightsFile))

		dojFile, err := os.Open(pathToDOJ)
		if err != nil {
			panic(err)
		}

		dojInformation, _ := data.NewDOJInformation(csv.NewReader(dojFile))

		cmsWriter := NewCMSWriter(path.Join(outputDir, "results.csv"))

		comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

		dataProcessor = NewDataProcessor(cmsCSV, weightsInformation, dojInformation, cmsWriter, comparisonTime)
	})

	It("runs and has output", func() {
		dataProcessor.Process()

		pathToExpectedResults, err := path.Abs(path.Join("..", "test_fixtures", "felonies_sf_results.csv"))
		Expect(err).ToNot(HaveOccurred())
		expectedResultsBody, err := ioutil.ReadFile(pathToExpectedResults)
		Expect(err).ToNot(HaveOccurred())

		pathToOutput, err := path.Abs(path.Join(outputDir, "results.csv"))
		Expect(err).ToNot(HaveOccurred())
		outputBody, err := ioutil.ReadFile(pathToOutput)
		Expect(err).ToNot(HaveOccurred())

		format.TruncatedDiff = false
		Expect(string(outputBody)).To(Equal(string(expectedResultsBody)))
	})
})
