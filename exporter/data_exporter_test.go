package exporter_test

import (
	"encoding/csv"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gstruct"
	"gogen_pilots/data"
	. "gogen_pilots/exporter"
	. "gogen_pilots/test_fixtures"
	"gogen_pilots/utilities"
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
			outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			age := 50

			var yearsConvictionFree int
			yearsConvictionFree = 10

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, data.EligibilityFlows["LOS ANGELES"])
			dojEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["LOS ANGELES"], age, yearsConvictionFree)
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"], age, yearsConvictionFree)
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"], age, yearsConvictionFree)

			dojWriter, _ := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter, _ := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter, _ := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen_pilots.out")

			dataExporter = NewDataExporter(
				dojInformation,
				dojEligibilities,
				dismissAllProp64Eligibilities,
				dismissAllProp64AndRelatedEligibilities,
				dojWriter,
				dojCondensedWriter,
				dojProp64ConvictionsWriter,
				outputWriter)
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
		COUNTY := "LOS ANGELES"
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			flow := createFlow()
			age := 50

			var yearsConvictionFree int
			yearsConvictionFree = 10

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, flow)
			dojEligibilities := dojInformation.DetermineEligibility(COUNTY, flow, age, yearsConvictionFree)
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility(COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64"], age, yearsConvictionFree)
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility(COUNTY, data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"], age, yearsConvictionFree)

			dojResultsPath := path.Join(outputDir, "results.csv")
			dojCondensedResultsPath := path.Join(outputDir, "condensed.csv")

			dojWriter, _ := NewDOJWriter(dojResultsPath)
			dojCondensedWriter, _ := NewCondensedDOJWriter(dojCondensedResultsPath)
			dojProp64ConvictionsWriter, _ := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen_pilots.out")

			dataExporter = NewDataExporter(
				dojInformation,
				dojEligibilities,
				dismissAllProp64Eligibilities,
				dismissAllProp64AndRelatedEligibilities,
				dojWriter, dojCondensedWriter,
				dojProp64ConvictionsWriter,
				outputWriter)
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

			condensedInputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
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
			outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
			pathToDOJ, pathToExpectedDOJResults, err = ExtractFullCSVFixtures(inputPath)
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)
			age := 50

			var yearsConvictionFree int
			yearsConvictionFree = 10

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, data.EligibilityFlows["LOS ANGELES"])
			dojEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["LOS ANGELES"], age, yearsConvictionFree)
			dismissAllProp64Eligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64"], age, yearsConvictionFree)
			dismissAllProp64AndRelatedEligibilities := dojInformation.DetermineEligibility("LOS ANGELES", data.EligibilityFlows["DISMISS ALL PROP 64 AND RELATED"], age, yearsConvictionFree)

			dojWriter, _ := NewDOJWriter(path.Join(outputDir, "results.csv"))
			dojCondensedWriter, _ := NewDOJWriter(path.Join(outputDir, "condensed.csv"))
			dojProp64ConvictionsWriter, _ := NewDOJWriter(path.Join(outputDir, "convictions.csv"))
			outputWriter := utilities.GetOutputWriter("gogen_pilots.out")

			dataExporter = NewDataExporter(
				dojInformation,
				dojEligibilities,
				dismissAllProp64Eligibilities,
				dismissAllProp64AndRelatedEligibilities,
				dojWriter,
				dojCondensedWriter,
				dojProp64ConvictionsWriter,
				outputWriter)
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

	Describe("AccumulateSummaryData", func() {
		It("adds new stats to stats already accumulated", func() {
			existingStats := Summary{
				County:               "SANTA CARLA",
				IndividualDismissAge: 50,
				YearsConvictionFree:  10,
				LineCount:            21,
				EarliestConviction:   time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC),
				ReliefWithCurrentEligibilityChoices: map[string]int{
					"CountSubjectsNoFelony":               2,
					"CountSubjectsNoConvictionLast7Years": 3,
					"CountSubjectsNoConviction":           1,
				},
				ReliefWithDismissAllProp64: map[string]int{
					"CountSubjectsNoFelony":               5,
					"CountSubjectsNoConvictionLast7Years": 7,
					"CountSubjectsNoConviction":           4,
				},
				Prop64ConvictionsCountInCountyByCodeSection: map[string]int{
					"11357": 4,
					"11358": 6,
					"11359": 7,
				},
			}

			newStats := Summary{
				LineCount:          25,
				EarliestConviction: time.Date(1983, 6, 1, 0, 0, 0, 0, time.UTC),
				ReliefWithCurrentEligibilityChoices: map[string]int{
					"CountSubjectsNoFelony":               1,
					"CountSubjectsNoConvictionLast7Years": 5,
					"CountSubjectsNoConviction":           2,
				},
				ReliefWithDismissAllProp64: map[string]int{
					"CountSubjectsNoFelony":               4,
					"CountSubjectsNoConvictionLast7Years": 6,
					"CountSubjectsNoConviction":           3,
				},
				Prop64ConvictionsCountInCountyByCodeSection: map[string]int{
					"11357": 5,
					"11358": 7,
					"11359": 8,
				},
			}

			cumulativeStats := dataExporter.AccumulateSummaryData(existingStats, newStats)

			Expect(cumulativeStats).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"County":               Equal("SANTA CARLA"),
				"IndividualDismissAge": Equal(50),
				"YearsConvictionFree":  Equal(10),
				"LineCount":            Equal(46),
				"EarliestConviction":   Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
				"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(3),
					"CountSubjectsNoConvictionLast7Years": Equal(8),
					"CountSubjectsNoConviction":           Equal(3),
				}),
				"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(9),
					"CountSubjectsNoConvictionLast7Years": Equal(13),
					"CountSubjectsNoConviction":           Equal(7),
				}),
				"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(9),
					"11358": Equal(13),
					"11359": Equal(15),
				}),
			}))
		})

		It("does not use an empty date as the earliest date", func() {
			existingStats := Summary{}

			newStats := Summary{
				County:             "SANTA CARLA",
				LineCount:          25,
				EarliestConviction: time.Date(1983, 6, 1, 0, 0, 0, 0, time.UTC),
				ReliefWithCurrentEligibilityChoices: map[string]int{
					"CountSubjectsNoFelony":               1,
					"CountSubjectsNoConvictionLast7Years": 5,
					"CountSubjectsNoConviction":           2,
				},
				ReliefWithDismissAllProp64: map[string]int{
					"CountSubjectsNoFelony":               4,
					"CountSubjectsNoConvictionLast7Years": 6,
					"CountSubjectsNoConviction":           3,
				},
				Prop64ConvictionsCountInCountyByCodeSection: map[string]int{
					"11357": 5,
					"11358": 7,
					"11359": 8,
				},
			}

			cumulativeStats := dataExporter.AccumulateSummaryData(existingStats, newStats)

			Expect(cumulativeStats.EarliestConviction).To(Equal(time.Date(1983, 6, 1, 0, 0, 0, 0, time.UTC)))
		})
	})
})

var _ = Describe("PossibleP64ChargeOnlyInComment", func() {
	It("returns the comment text if the comment text has a Prop64 charge and the offense description doesn't", func() {
		Expect(PossibleP64ChargeOnlyInComment("912", "11357(A)")).To(Equal("11357(A)"))
		Expect(PossibleP64ChargeOnlyInComment("11350", "11357(A)")).To(Equal("11357(A)"))
	})

	It("returns empty string if the comment text is expected to override the offense description", func() {
		Expect(PossibleP64ChargeOnlyInComment("", "11357(A)")).To(Equal(""))
		Expect(PossibleP64ChargeOnlyInComment("SEE COMMENT FOR CHARGE", "11357(A)")).To(Equal(""))
	})
	It("returns the comment text if offense description and comment text have different Prop64 charges that are not both 11357 subsections", func() {
		Expect(PossibleP64ChargeOnlyInComment("11358", "11359")).To(Equal("11359"))
		Expect(PossibleP64ChargeOnlyInComment("11357(A)", "11358")).To(Equal("11358"))
		Expect(PossibleP64ChargeOnlyInComment("11358", "11357(C)")).To(Equal("11357(C)"))
	})
	It("returns the comment text if offense description and comment text have different 11357 subsections, one from the (A,B) group and the other from the (C,D) group", func() {
		Expect(PossibleP64ChargeOnlyInComment("11357(A)", "11357(C)")).To(Equal("11357(C)"))
		Expect(PossibleP64ChargeOnlyInComment("11357(D)", "11357(B)")).To(Equal("11357(B)"))
	})
	It("returns the empty string if offense description and comment text have different 11357 subsections, but from the same group (A,B) or (C,D)", func() {
		Expect(PossibleP64ChargeOnlyInComment("11357(A)", "11357(B)")).To(Equal(""))
		Expect(PossibleP64ChargeOnlyInComment("11357(D)", "11357(C)")).To(Equal(""))
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

func createFlow() data.EligibilityFlow {
	flow, _ := data.EligibilityFlows["LOS ANGELES"]
	return flow
}
