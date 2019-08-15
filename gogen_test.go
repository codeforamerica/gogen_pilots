package main

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega/gstruct"
	"gogen/exporter"
	"io/ioutil"
	"os/exec"
	path "path/filepath"
	"time"

	. "gogen/test_fixtures"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func GetOutputSummary(filePath string) exporter.Summary {
	bytes, _ := ioutil.ReadFile(filePath)
	var summary exporter.Summary
	json.Unmarshal(bytes, &summary)
	return summary
}

var _ = Describe("gogen", func() {
	var (
		outputDir string
		pathToDOJ string
		err       error
	)
	It("can handle a csv with extra comma at the end of headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).ToNot(gbytes.Say("required"))

		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can handle an input file without headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "no_headers.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can accept a compute-at option for determining eligibility", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).ToNot(gbytes.Say("required"))
		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary.LineCount).To(Equal(37))
	})

	It("can accept a suffix for the output file names", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())
		dateSuffix := "Feb_8_2019_3.32.43.PM"

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", dateSuffix)
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		fileResultsOutputDir := path.Join(outputDir, fmt.Sprintf("DOJ_Input_File_1_Results_%s", dateSuffix))
		expectedDojResultsFileName := fmt.Sprintf("%v/doj_results_1_%s.csv", fileResultsOutputDir, dateSuffix)
		expectedCondensedFileName := fmt.Sprintf("%v/doj_results_condensed_1_%s.csv", fileResultsOutputDir, dateSuffix)
		expectedConvictionsFileName := fmt.Sprintf("%v/doj_results_convictions_1_%s.csv", fileResultsOutputDir, dateSuffix)
		expectedOutputFileName := fmt.Sprintf("%v/gogen_1_%s.out", fileResultsOutputDir, dateSuffix)
		expectedJsonOutputFileName := fmt.Sprintf("%v/gogen_%s.json", outputDir, dateSuffix)

		Ω(expectedDojResultsFileName).Should(BeAnExistingFile())
		Ω(expectedCondensedFileName).Should(BeAnExistingFile())
		Ω(expectedConvictionsFileName).Should(BeAnExistingFile())
		Ω(expectedOutputFileName).Should(BeAnExistingFile())
		Ω(expectedJsonOutputFileName).Should(BeAnExistingFile())
	})

	It("validates required options", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(3))
		Eventually(session.Err).Should(gbytes.Say("missing required field: Run gogen --help for more info"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen%s.err", outputDir, "")

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		data, _ := ioutil.ReadFile(expectedErrorFileName)
		Expect(string(data)).To(Equal("missing required field: Run gogen --help for more info\n"))
	})

	It("fails and reports errors for missing input files", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "missing.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())
		filenameSuffix := "a_suffix"

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(2))
		Eventually(session.Err).Should(gbytes.Say("open .*missing.csv: no such file or directory"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_%s.err", outputDir, filenameSuffix)

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		data, _ := ioutil.ReadFile(expectedErrorFileName)
		Expect(string(data)).To(MatchRegexp("open .*missing.csv: no such file or directory"))
	})

	It("fails and reports errors for invalid input files", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "bad.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())
		filenameSuffix := "a_suffix"

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(2))
		Eventually(session.Err).Should(gbytes.Say("record on line 2: wrong number of fields"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_%s.err", outputDir, filenameSuffix)

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		data, _ := ioutil.ReadFile(expectedErrorFileName)
		Expect(string(data)).To(MatchRegexp("record on line 2: wrong number of fields"))
	})

	It("can accept path to eligibility options file", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
		Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
			"County":                  Equal("SACRAMENTO"),
			"LineCount":               Equal(37),
			"ProcessingTimeInSeconds": BeNumerically(">", 0),
			"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
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
				"11357": Equal(3),
				"11358": Equal(6),
				"11359": Equal(8),
			}),
			"SubjectsWithProp64ConvictionCountInCounty": Equal(11),
			"Prop64FelonyConvictionsCountInCounty":      Equal(14),
			"Prop64NonFelonyConvictionsCountInCounty":   Equal(3),
			"SubjectsWithSomeReliefCount":               Equal(11),
			"ConvictionDismissalCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11357": Equal(2),
				"11358": Equal(5),
			}),
			"ConvictionReductionCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11359": Equal(1),
				"11360": Equal(0),
			}),
			"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
				"21 years or younger":                      Equal(1),
				"57 years or older":                        Equal(2),
				"Conviction occurred 10 or more years ago": Equal(1),
				"Individual is deceased":                   Equal(1),
				"Only has 11357-60 charges":                Equal(1),
			}),
		}))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 37 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 12 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 29 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 26 convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions Overall--------------------"))
		Eventually(session).Should(gbytes.Say("Found 20 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 8 11358 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 9 11359 convictions total"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 17 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 6 11358 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 8 11359 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Date of earliest Prop 64 conviction: June 1979"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 6 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 7 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 16 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions total that are Eligible for Reduction"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason 21 years or younger"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason 57 years or older"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Conviction occurred 10 or more years ago"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason Dismiss all HS 11357 convictions"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason Dismiss all HS 11358 convictions"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Individual is deceased"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions with eligibility reason Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Only has 11357-60 charges"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Reduce all HS 11359 convictions"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 0 convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("12 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("12 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility is run as specified for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("4 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("4 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("4 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))
	})

	Describe("Processing multiple input files", func() {
		It("nests and indexes the names of the results files for each input file", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
			inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

			pathToGogen, err := gexec.Build("gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV+","+inputCSV)
			countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
			computeAtFlag := "--compute-at=2019-11-11"
			eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			file1OutputDir := path.Join(outputDir, "DOJ_Input_File_1_Results")
			file2OutputDir := path.Join(outputDir, "DOJ_Input_File_2_Results")
			expectedDojResultsFile1Name := fmt.Sprintf("%v/doj_results_1.csv", file1OutputDir)
			expectedDojResultsFile2Name := fmt.Sprintf("%v/doj_results_2.csv", file2OutputDir)
			expectedCondensedFile1Name := fmt.Sprintf("%v/doj_results_condensed_1.csv", file1OutputDir)
			expectedCondensedFile2Name := fmt.Sprintf("%v/doj_results_condensed_2.csv", file2OutputDir)
			expectedConvictionsFile1Name := fmt.Sprintf("%v/doj_results_convictions_1.csv", file1OutputDir)
			expectedConvictionsFile2Name := fmt.Sprintf("%v/doj_results_convictions_2.csv", file2OutputDir)
			expectedOutputFile1Name := fmt.Sprintf("%v/gogen_1.out", file1OutputDir)
			expectedOutputFile2Name := fmt.Sprintf("%v/gogen_2.out", file2OutputDir)
			expectedJsonOutputFileName := fmt.Sprintf("%v/gogen.json", outputDir)

			Ω(expectedDojResultsFile1Name).Should(BeAnExistingFile())
			Ω(expectedDojResultsFile2Name).Should(BeAnExistingFile())
			Ω(expectedCondensedFile1Name).Should(BeAnExistingFile())
			Ω(expectedCondensedFile2Name).Should(BeAnExistingFile())
			Ω(expectedConvictionsFile1Name).Should(BeAnExistingFile())
			Ω(expectedConvictionsFile2Name).Should(BeAnExistingFile())
			Ω(expectedOutputFile1Name).Should(BeAnExistingFile())
			Ω(expectedOutputFile2Name).Should(BeAnExistingFile())
			Ω(expectedJsonOutputFileName).Should(BeAnExistingFile())
		})

		It("can aggregate statistics for multiple input files", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
			inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

			pathToGogen, err := gexec.Build("gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV+","+inputCSV)
			countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
			computeAtFlag := "--compute-at=2019-11-11"
			eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			summary := GetOutputSummary(path.Join(outputDir, "gogen.json"))
			Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
				"County":                  Equal("SACRAMENTO"),
				"LineCount":               Equal(74),
				"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
				"ProcessingTimeInSeconds": BeNumerically(">", 0),
				"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(8),
					"CountSubjectsNoConviction":           Equal(6),
					"CountSubjectsNoConvictionLast7Years": Equal(2),
				}),
				"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(8),
					"CountSubjectsNoConviction":           Equal(6),
					"CountSubjectsNoConvictionLast7Years": Equal(2),
				}),
				"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(6),
					"11358": Equal(12),
					"11359": Equal(16),
				}),
				"SubjectsWithProp64ConvictionCountInCounty": Equal(22),
				"Prop64FelonyConvictionsCountInCounty":      Equal(28),
				"Prop64NonFelonyConvictionsCountInCounty":   Equal(6),
				"SubjectsWithSomeReliefCount":               Equal(22),
				"ConvictionDismissalCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(4),
					"11358": Equal(10),
				}),
				"ConvictionReductionCountByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11359": Equal(2),
					"11360": Equal(0),
				}),
				"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
					"21 years or younger":                      Equal(2),
					"57 years or older":                        Equal(4),
					"Conviction occurred 10 or more years ago": Equal(2),
					"Individual is deceased":                   Equal(2),
					"Only has 11357-60 charges":                Equal(2),
				}),
			}))
		})

		It("can return errors for multiple input files", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
			pathToValidDOJ, _, _ := ExtractFullCSVFixtures(pathToInputExcel)
			pathToBadDOJ, err := path.Abs(path.Join("test_fixtures", "bad.csv"))
			pathToMissingDOJ, err := path.Abs(path.Join("test_fixtures", "missing.csv"))

			pathToGogen, err := gexec.Build("gogen")
			Expect(err).ToNot(HaveOccurred())
			filenameSuffix := "a_suffix"

			pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", pathToValidDOJ+","+pathToBadDOJ+","+pathToMissingDOJ)

			countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
			computeAtFlag := "--compute-at=2019-11-11"
			dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)
			eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag, dateTimeFlag, eligibilityOptionsFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(2))
			Eventually(session.Err).Should(gbytes.Say("record on line 2: wrong number of fields"))
			Eventually(session.Err).Should(gbytes.Say("open .*missing.csv: no such file or directory"))

			expectedErrorFileName := fmt.Sprintf("%v/gogen_%s.err", outputDir, filenameSuffix)

			Ω(expectedErrorFileName).Should(BeAnExistingFile())
			data, _ := ioutil.ReadFile(expectedErrorFileName)
			Expect(string(data)).To(MatchRegexp("open .*missing.csv: no such file or directory"))
			Expect(string(data)).To(MatchRegexp("record on line 2: wrong number of fields"))
		})
	})
})
