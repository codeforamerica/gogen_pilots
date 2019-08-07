package exporter_test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gstruct"
	"gogen/data"
	. "gogen/exporter"
	. "gogen/test_fixtures"
	"gogen/utilities"
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

	Describe("Los Angeles", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, data.EligibilityFlows["LOS ANGELES"])
			dojEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["LOS ANGELES"])
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen.out")

			dataExporter = NewDataExporter(dojInformation,
				dojEligibilities, dismissAllProp64Eligibilities, dismissAllProp64AndRelatedEligibilities, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter, outputWriter, path.Join(outputDir, "gogen.json"))
		})

		It("runs and has output", func() {
			dataExporter.Export("LOS ANGELES", time.Now())
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
		COUNTY := "SACRAMENTO"
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "configurable_flow.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dismissCodeSections := []string{"11357(a)", "11357(c)", "11357(d)", "11357(no-sub-section)", "11358"}
			reduceCodeSections := []string{"11357(b)", "11359", "11360"}
			flow := createFlow(dismissCodeSections, reduceCodeSections, COUNTY)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, flow)
			dojEligibilities := dojInformation.DetermineEligibility(COUNTY, flow)
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility(COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility(COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojResultsPath := path.Join(outputDir, "results.csv")
			dojCondensedResultsPath := path.Join(outputDir, "condensed.csv")

			dojWriter := NewDOJWriter(dojResultsPath)
			dojCondensedWriter := NewCondensedDOJWriter(dojCondensedResultsPath)
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen.out")

			dataExporter = NewDataExporter(dojInformation,
				dojEligibilities, dismissAllProp64Eligibilities, dismissAllProp64AndRelatedEligibilities, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter, outputWriter, path.Join(outputDir, "gogen.json"))
		})

		It("runs and has condensed output", func() {
			dataExporter.Export(COUNTY, time.Now())
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "condensed.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			condensedInputPath := path.Join("..", "test_fixtures", "configurable_flow.xlsx")
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

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, data.EligibilityFlows["LOS ANGELES"])
			dojEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["LOS ANGELES"])
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen.out")

			dataExporter = NewDataExporter(dojInformation,
				dojEligibilities, dismissAllProp64Eligibilities, dismissAllProp64AndRelatedEligibilities, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter, outputWriter, path.Join(outputDir, "gogen.json"))
		})

		It("runs and has condensed output", func() {
			dataExporter.Export("LOS ANGELES", time.Now())
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

			dismissCodeSections := []string{"11357(a)", "11357(c)", "11357(d)", "11357(no-sub-section)", "11358"}
			reduceCodeSections := []string{"11357(b)", "11359", "11360"}
			flow := createFlow(dismissCodeSections, reduceCodeSections, COUNTY)

			dojInformation := data.NewDOJInformation(pathToDOJ, comparisonTime, flow)
			dojEligibilities := dojInformation.DetermineEligibility(COUNTY, flow)
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility(COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64"])
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility(COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"])

			dojWriter := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen.out")

			dataExporter = NewDataExporter(dojInformation,
				dojEligibilities, dismissAllProp64Eligibilities, dismissAllProp64AndRelatedEligibilities, dojWriter, dojCondensedWriter, dojProp64ConvictionsWriter, outputWriter, path.Join(outputDir, "gogen.json"))
		})

		It("runs and has output", func() {
			dataExporter.Export(COUNTY, time.Now())
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

			bytes, _ := ioutil.ReadFile(path.Join(outputDir, "gogen.json"))
			var summary Summary
			json.Unmarshal(bytes, &summary)
			Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
				"County":                  Equal("SACRAMENTO"),
				"LineCount":               Equal(36),
				"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
				"ProcessingTimeInSeconds": BeNumerically(">", 0),
				"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(4),
					"CountSubjectsNoConviction":           Equal(3),
					"CountSubjectsNoConvictionLast7Years": Equal(1),
				}),
				"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(4),
					"CountSubjectsNoConviction":           Equal(3),
					"CountSubjectsNoConvictionLast7Years": Equal(1),
				}),
				"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(4),
					"11358": Equal(6),
					"11359": Equal(7),
				}),
				"SubjectsWithProp64ConvictionCountInCounty": Equal(0),
				"Prop64FelonyConvictionsCountInCounty":      Equal(0),
				"Prop64MisdemeanorConvictionsCountInCounty": Equal(0),
				"SubjectsWithSomeReliefCount":               Equal(0),
				"ConvictionDismissalCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357(c)":              Equal(1),
					"11357(no sub-section)": Equal(1),
					"11358":                 Equal(5),
				}),
				"ConvictionReductionCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11359": Equal(1),
				}),
				"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
					"21 years or younger":                      Equal(1),
					"57 years or older":                        Equal(2),
					"Conviction occurred 10 or more years ago": Equal(1),
					"Individual is deceased":                   Equal(1),
					"Misdemeanor or Infraction":                Equal(3),
					"Only has 11357-60 charges":                Equal(1),
				}),
			}))

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

func createFlow(dismissCodeSections []string, reduceCodeSections []string, county string) data.EligibilityFlow {
	return data.NewConfigurableEligibilityFlow(data.EligibilityOptions{
		BaselineEligibility: data.BaselineEligibility{
			Dismiss: dismissCodeSections,
			Reduce:  reduceCodeSections,
		},
		AdditionalRelief: data.AdditionalRelief{
			SubjectUnder21AtConviction:    true,
			SubjectAgeThreshold:           57,
			YearsSinceConvictionThreshold: 10,
			SubjectHasOnlyProp64Charges:   true,
			SubjectIsDeceased:             true,
		},
	}, county)
}
