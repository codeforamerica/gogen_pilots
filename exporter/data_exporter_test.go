package exporter_test

import (
	"encoding/csv"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"gogen/data"
	. "gogen/exporter"
	. "gogen/test_fixtures"
	"io/ioutil"
	"os"
	path "path/filepath"
	"time"
)

var _ = Describe("DataExporter", func() {
	var (
		outputDir                string
		dataExporter             DataExporter
		pathToDOJ                string
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
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "SACRAMENTO", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "SACRAMENTO", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataExporter.Export("SACRAMENTO")
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
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "SAN JOAQUIN", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "SAN JOAQUIN", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataExporter.Export("SAN JOAQUIN")
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
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataExporter.Export("CONTRA COSTA")
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
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataExporter.Export("LOS ANGELES")
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
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojResultsPath := path.Join(outputDir, "results.csv")
			dojCondensedResultsPath := path.Join(outputDir, "condensed.csv")

			dojWriter := NewDOJWriter(dojResultsPath)
			dojCondensedWriter := NewCondensedDOJWriter(dojCondensedResultsPath)
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has condensed output", func() {
			dataExporter.Export("CONTRA COSTA")
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

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["LOS ANGELES"])
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has condensed output", func() {
			dataExporter.Export("LOS ANGELES")
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

	Describe("Configurable eligibility flow", func() {
		var COUNTY = "SACRAMENTO"
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "configurable_flow.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			flow := data.NewConfigurableEligibilityFlow(data.EligibilityOptions{
				BaselineEligibility: data.BaselineEligibility{
					Dismiss: []string{"11357(A)", "11357(C)", "11357(D)", "11358"},
					Reduce:  []string{"11357(B)", "11359", "11360"},
				},
				AdditionalRelief: data.AdditionalRelief{
					Under21: true,
				},
			}, COUNTY)


			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, COUNTY, flow)
			dismissAllProp64DojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedDojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))

			dataExporter = NewDataExporter(dojInformation, dismissAllProp64DojInformation, dismissAllProp64AndRelatedDojInformation, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter)
		})

		It("runs and has output", func() {
			dataExporter.Export(COUNTY)
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
})

func expectCSVsToBeEqual(expectedCSV [][]string, actualCSV [][]string) {
	for i, row := range actualCSV {
		for j, item := range row {
			Expect(item).To(Equal(expectedCSV[i][j]), fmt.Sprintf("Failed on row %d, col %d\n", i+2, j+1))
		}
	}
	Expect(actualCSV).To(Equal(expectedCSV))
}
