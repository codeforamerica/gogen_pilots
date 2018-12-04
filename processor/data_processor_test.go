package processor_test

import (
	"encoding/csv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gogen/data"
	. "gogen/processor"
	"io/ioutil"
	"os"
	path "path/filepath"
)

var _ = Describe("DataProcessor", func() {
	var (
		outputDir          string
		weightsInformation *data.WeightsInformation
		dojInformation     *data.DOJInformation
		cmsCSVReader       *csv.Reader
		dataProcessor      DataProcessor
		err                error
	)

	BeforeEach(func() {
		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		outputWriter := NewCMSWriter(path.Join(outputDir, "felonies_sf_results.csv"))

		pathToWeights, err := path.Abs(path.Join("..", "test_fixtures", "conviction_weights.csv"))
		Expect(err).ToNot(HaveOccurred())
		weightsFile, err := os.Open(pathToWeights)
		Expect(err).ToNot(HaveOccurred())
		weightsCSVReader := csv.NewReader(weightsFile)
		weightsInformation, err = data.NewWeightsInformation(weightsCSVReader)
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "cadoj.csv"))
		Expect(err).ToNot(HaveOccurred())
		dojFile, err := os.Open(pathToDOJ)
		Expect(err).ToNot(HaveOccurred())
		dojCSVReader := csv.NewReader(dojFile)
		dojInformation, err = data.NewDOJInformation(dojCSVReader)
		Expect(err).ToNot(HaveOccurred())

		pathToCMS, err := path.Abs(path.Join("..", "test_fixtures", "felonies_sf.csv"))
		Expect(err).ToNot(HaveOccurred())
		cmsFile, err := os.Open(pathToCMS)
		Expect(err).ToNot(HaveOccurred())
		cmsCSVReader = csv.NewReader(cmsFile)

		dataProcessor = NewDataProcessor(cmsCSVReader, weightsInformation, dojInformation, outputWriter)
	})

	PIt("runs and has output", func() {
		dataProcessor.Process()

		pathToExpectedResults, err := path.Abs(path.Join("..", "test_fixtures", "felonies_sf_results.csv"))
		Expect(err).ToNot(HaveOccurred())
		expectedResultsBody, err := ioutil.ReadFile(pathToExpectedResults)
		Expect(err).ToNot(HaveOccurred())

		pathToOutput, err := path.Abs(path.Join(outputDir, "felonies_sf_results.csv"))
		Expect(err).ToNot(HaveOccurred())
		outputBody, err := ioutil.ReadFile(pathToOutput)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(outputBody)).To(Equal(string(expectedResultsBody)))
	})
})
