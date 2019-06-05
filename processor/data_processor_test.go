package processor_test

import (
	"encoding/csv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"gogen/data"
	. "gogen/processor"
	. "gogen/test_fixtures"
	"io/ioutil"
	"os"
	path "path/filepath"
	"strconv"
	"time"
)

var _ = Describe("DataProcessor", func() {
	var (
		outputDir                string
		dataProcessor            DataProcessor
		pathToDOJ 				 string
		pathToExpectedDOJResults string
		err                      error
	)

	Describe("Sacramento", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "sacramento.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "SACRAMENTO", data.EligibilityFlows["SACRAMENTO"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("SACRAMENTO")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			expectCSVsToBeEqual(expectedDOJResultsCSV, outputDOJCSV)
		})
	})

	Describe("San Joaquin", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "san_joaquin.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "SAN JOAQUIN", data.EligibilityFlows["SAN JOAQUIN"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("SAN JOAQUIN")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			expectCSVsToBeEqual(expectedDOJResultsCSV, outputDOJCSV)
		})
	})

	Describe("Contra Costa", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "contra_costa.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA", data.EligibilityFlows["CONTRA COSTA"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("CONTRA COSTA")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			expectCSVsToBeEqual(expectedDOJResultsCSV, outputDOJCSV)
		})
	})

	Describe("Los Angeles", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["LOS ANGELES"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("LOS ANGELES")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			expectCSVsToBeEqual(expectedDOJResultsCSV, outputDOJCSV)
		})
	})

	Describe("Condensed columns output file", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "contra_costa.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA", data.EligibilityFlows["CONTRA COSTA"])

			dojResultsPath := path.Join(outputDir, "results.csv")
			dojCondensedResultsPath := path.Join(outputDir, "condensed.csv")

			dojWriter := NewDOJWriter(dojResultsPath)
			dojCondensedWriter := NewCondensedDOJWriter(dojCondensedResultsPath)
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has condensed output", func() {
			dataProcessor.Process("CONTRA COSTA")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "condensed.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			condensedInputPath := path.Join("..", "test_fixtures", "contra_costa.xlsx")
			expectedCondensedCSVResult, err := ExtractCondensedCSVFixture(condensedInputPath)
			ExpectedDOJResultsFile, err := os.Open(expectedCondensedCSVResult)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			expectCSVsToBeEqual(expectedDOJResultsCSV, outputDOJCSV)
		})
	})

	Describe("Prop 64 convictions output file", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation:= data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["LOS ANGELES"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has condensed output", func() {
			dataProcessor.Process("LOS ANGELES")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "convictions.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			condensedInputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			expectedProp64CSVResult, err := ExtractProp64ConvictionsCSVFixture(condensedInputPath)
			ExpectedDOJResultsFile, err := os.Open(expectedProp64CSVResult)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			expectCSVsToBeEqual(expectedDOJResultsCSV, outputDOJCSV)
		})
	})
})

func expectCSVsToBeEqual(expectedCSV [][]string, actualCSV [][]string) {
	for i, row := range actualCSV {
		for j, item := range row {
			Expect(item).To(Equal(expectedCSV[i][j]), "failed on row "+strconv.Itoa(i+1))
		}
	}
	Expect(actualCSV).To(Equal(expectedCSV))
}