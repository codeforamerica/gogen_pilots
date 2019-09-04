package main

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega/gstruct"
	"gogen_pilots/exporter"
	"io/ioutil"
	"os/exec"
	path "path/filepath"
	"regexp"
	"time"

	. "gogen_pilots/test_fixtures"

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

var _ = Describe("gogen_pilots", func() {
	var (
		outputDir string
		pathToDOJ string
		err       error
	)
	It("can handle a csv with extra comma at the end of headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).ToNot(gbytes.Say("required"))

		summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can handle an input file without headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "no_headers.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))
		Expect(summary.LineCount).To(Equal(38))
	})

	It("can accept a compute-at option for determining eligibility", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Err).ToNot(gbytes.Say("required"))
		summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))
		Expect(summary.LineCount).To(Equal(32))
	})

	It("can accept a suffix for the output file names", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())
		dateSuffix := "Feb_8_2019_3.32.43.PM"

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", dateSuffix)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, dateTimeFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))

		fileResultsOutputDir := path.Join(outputDir, fmt.Sprintf("DOJ_Input_File_1_Results_%s", dateSuffix))
		expectedDojResultsFileName := fmt.Sprintf("%v/doj_results_1_%s.csv", fileResultsOutputDir, dateSuffix)
		expectedCondensedFileName := fmt.Sprintf("%v/doj_results_condensed_1_%s.csv", fileResultsOutputDir, dateSuffix)
		expectedConvictionsFileName := fmt.Sprintf("%v/doj_results_convictions_1_%s.csv", fileResultsOutputDir, dateSuffix)
		expectedOutputFileName := fmt.Sprintf("%v/gogen_pilots_1_%s.out", fileResultsOutputDir, dateSuffix)
		expectedJsonOutputFileName := fmt.Sprintf("%v/gogen_pilots_%s.json", outputDir, dateSuffix)

		Ω(expectedDojResultsFileName).Should(BeAnExistingFile())
		Ω(expectedCondensedFileName).Should(BeAnExistingFile())
		Ω(expectedConvictionsFileName).Should(BeAnExistingFile())
		Ω(expectedOutputFileName).Should(BeAnExistingFile())
		Ω(expectedJsonOutputFileName).Should(BeAnExistingFile())
	})
	It("can accept an age for a participant to determine eligibility", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		var age float64
		age = 40
		ageFlag := fmt.Sprintf("--individual-age=%v", age)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, ageFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))
		Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
			"County":                  Equal("LOS ANGELES"),
			"LineCount":               Equal(32),
			"ProcessingTimeInSeconds": BeNumerically(">", 0),
			"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
			"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(1),
				"CountSubjectsNoConviction":           Equal(1),
				"CountSubjectsNoConvictionLast7Years": Equal(0),
			}),
			"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(2),
				"CountSubjectsNoConviction":           Equal(2),
				"CountSubjectsNoConvictionLast7Years": Equal(3),
			}),
			"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11357": Equal(3),
				"11358": Equal(8),
				"11359": Equal(4),
			}),
			"SubjectsWithProp64ConvictionCountInCounty": Equal(0),
			"Prop64FelonyConvictionsCountInCounty":      Equal(0),
			"Prop64MisdemeanorConvictionsCountInCounty": Equal(0),
			"SubjectsWithSomeReliefCount":               Equal(0),
			"ConvictionDismissalCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
			"ConvictionReductionCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
			"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
				"21 years or younger": Equal(1),
				"40 years or older":   Equal(4),
				"Only has 11357-60 charges and completed sentence": Equal(1),
				"11357(a) or 11357(b)":                             Equal(1),
			}),
		}))
	})

	It("can accept a number of years conviction free as a parameter to determine eligibility", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		var yearsConvictionFree int
		yearsConvictionFree = 2
		yearsConvictionFreeFlag := fmt.Sprintf("--years-conviction-free=%v", yearsConvictionFree)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, yearsConvictionFreeFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))
		Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
			"County":                  Equal("LOS ANGELES"),
			"LineCount":               Equal(32),
			"ProcessingTimeInSeconds": BeNumerically(">", 0),
			"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
			"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(2),
				"CountSubjectsNoConviction":           Equal(2),
				"CountSubjectsNoConvictionLast7Years": Equal(1),
			}),
			"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(2),
				"CountSubjectsNoConviction":           Equal(2),
				"CountSubjectsNoConvictionLast7Years": Equal(3),
			}),
			"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11357": Equal(3),
				"11358": Equal(8),
				"11359": Equal(4),
			}),
			"SubjectsWithProp64ConvictionCountInCounty": Equal(0),
			"Prop64FelonyConvictionsCountInCounty":      Equal(0),
			"Prop64MisdemeanorConvictionsCountInCounty": Equal(0),
			"SubjectsWithSomeReliefCount":               Equal(0),
			"ConvictionDismissalCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
			"ConvictionReductionCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
			"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
				"21 years or younger": Equal(1),
				"50 years or older":   Equal(4),
				"Only has 11357-60 charges and completed sentence": Equal(1),
				"11357(a) or 11357(b)":                             Equal(1),
				"No convictions in past 2 years":                   Equal(1),
			}),
		}))
	})

	It("validates required options", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(3))
		Eventually(session.Err).Should(gbytes.Say("missing required field: Run gogen_pilots --help for more info"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_pilots%s.err", outputDir, "")

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		data, _ := ioutil.ReadFile(expectedErrorFileName)
		Expect(string(data)).To(Equal("missing required field: Run gogen_pilots --help for more info\n"))
	})

	It("fails and reports errors for missing input files", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "missing.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())
		filenameSuffix := "a_suffix"

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, dateTimeFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(2))
		Eventually(session.Err).Should(gbytes.Say("open .*missing.csv: no such file or directory"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_pilots_%s.err", outputDir, filenameSuffix)

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		data, _ := ioutil.ReadFile(expectedErrorFileName)
		Expect(string(data)).To(MatchRegexp("open .*missing.csv: no such file or directory"))
	})

	It("fails and reports errors for invalid input files", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "bad.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())
		filenameSuffix := "a_suffix"

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		computeAtFlag := "--compute-at=2019-11-11"
		dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, dateTimeFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(2))
		Eventually(session.Err).Should(gbytes.Say("record on line 2: wrong number of fields"))

		expectedErrorFileName := fmt.Sprintf("%v/gogen_pilots_%s.err", outputDir, filenameSuffix)

		Ω(expectedErrorFileName).Should(BeAnExistingFile())
		data, _ := ioutil.ReadFile(expectedErrorFileName)
		Expect(string(data)).To(MatchRegexp("record on line 2: wrong number of fields"))
	})

	It("runs and has output for Los Angeles", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))
		Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
			"County":                  Equal("LOS ANGELES"),
			"LineCount":               Equal(35),
			"ProcessingTimeInSeconds": BeNumerically(">", 0),
			"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
			"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(1),
				"CountSubjectsNoConviction":           Equal(1),
				"CountSubjectsNoConvictionLast7Years": Equal(0),
			}),
			"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
				"CountSubjectsNoFelony":               Equal(2),
				"CountSubjectsNoConviction":           Equal(2),
				"CountSubjectsNoConvictionLast7Years": Equal(3),
			}),
			"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
				"11357": Equal(3),
				"11358": Equal(9),
				"11359": Equal(4),
			}),
			"SubjectsWithProp64ConvictionCountInCounty": Equal(0),
			"Prop64FelonyConvictionsCountInCounty":      Equal(0),
			"Prop64MisdemeanorConvictionsCountInCounty": Equal(0),
			"SubjectsWithSomeReliefCount":               Equal(0),
			"ConvictionDismissalCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
			"ConvictionReductionCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
			"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
				"21 years or younger": Equal(1),
				"50 years or older":   Equal(4),
				"Only has 11357-60 charges and completed sentence": Equal(1),
				"11357(a) or 11357(b)":                             Equal(1),
			}),
		}))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 35 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Based on your office’s eligibility choices, this application processed the data in .* seconds"))
		Eventually(session).Should(gbytes.Say("Found 10 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 28 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 25 convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions Overall--------------------"))
		Eventually(session).Should(gbytes.Say("Found 19 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 11 11358 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 5 11359 convictions total"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 16 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 9 11358 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 4 11359 convictions in this county"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 7 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 1 11357 convictions that are Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions total that are Hand Review"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions total that are Not eligible"))

		Eventually(session).Should(gbytes.Say("To be reviewed by City Attorneys"))
		Eventually(session).Should(gbytes.Say("Found 1 11357 convictions that are To be reviewed by City Attorney"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are To be reviewed by City Attorney"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are To be reviewed by City Attorney"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions total that are To be reviewed by City Attorney"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason 11357(a) or 11357(b)")))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason 21 years or younger"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions with eligibility reason 50 years or older"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Only has 11357-60 charges and completed sentence"))

		Eventually(session).Should(gbytes.Say("Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Currently serving sentence"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Other 11357"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason PC 290")))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 3 convictions with eligibility reason PC 667(e)(2)(c)(iv)")))

		Eventually(session).Should(gbytes.Say("To be reviewed by City Attorneys"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions with eligibility reason Misdemeanor or Infraction"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 1 466 PC convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("10 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("10 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("5 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility is run as specified for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("1 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("0 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))
	})

	Describe("Processing multiple input files", func() {
		It("nests and indexes the names of the results files for each input file", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
			inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

			pathToGogen, err := gexec.Build("gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV+","+inputCSV)
			computeAtFlag := "--compute-at=2019-11-11"

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag)
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
			expectedOutputFile1Name := fmt.Sprintf("%v/gogen_pilots_1.out", file1OutputDir)
			expectedOutputFile2Name := fmt.Sprintf("%v/gogen_pilots_2.out", file2OutputDir)
			expectedJsonOutputFileName := fmt.Sprintf("%v/gogen_pilots.json", outputDir)

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
			outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
			inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

			pathToGogen, err := gexec.Build("gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV+","+inputCSV)
			computeAtFlag := "--compute-at=2019-11-11"

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))

			summary := GetOutputSummary(path.Join(outputDir, "gogen_pilots.json"))

			Expect(summary).To(gstruct.MatchAllFields(gstruct.Fields{
				"County":                  Equal("LOS ANGELES"),
				"LineCount":               Equal(64),
				"EarliestConviction":      Equal(time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)),
				"ProcessingTimeInSeconds": BeNumerically(">", 0),
				"ReliefWithCurrentEligibilityChoices": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(2),
					"CountSubjectsNoConviction":           Equal(2),
					"CountSubjectsNoConvictionLast7Years": Equal(0),
				}),
				"ReliefWithDismissAllProp64": gstruct.MatchAllKeys(gstruct.Keys{
					"CountSubjectsNoFelony":               Equal(4),
					"CountSubjectsNoConviction":           Equal(4),
					"CountSubjectsNoConvictionLast7Years": Equal(6),
				}),
				"Prop64ConvictionsCountInCountyByCodeSection": gstruct.MatchAllKeys(gstruct.Keys{
					"11357": Equal(6),
					"11358": Equal(16),
					"11359": Equal(8),
				}),
				"SubjectsWithProp64ConvictionCountInCounty": Equal(0),
				"Prop64FelonyConvictionsCountInCounty":      Equal(0),
				"Prop64MisdemeanorConvictionsCountInCounty": Equal(0),
				"SubjectsWithSomeReliefCount":               Equal(0),
				"ConvictionDismissalCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
				"ConvictionReductionCountByCodeSection":     gstruct.MatchAllKeys(gstruct.Keys{}),
				"ConvictionDismissalCountByAdditionalRelief": gstruct.MatchAllKeys(gstruct.Keys{
					"21 years or younger": Equal(2),
					"50 years or older":   Equal(8),
					"Only has 11357-60 charges and completed sentence": Equal(2),
					"11357(a) or 11357(b)":                             Equal(2),
				}),
			}))
		})

		It("can return errors for multiple input files", func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen_pilots")
			Expect(err).ToNot(HaveOccurred())

			pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
			pathToValidDOJ, _, _ := ExtractFullCSVFixtures(pathToInputExcel)
			pathToBadDOJ, err := path.Abs(path.Join("test_fixtures", "bad.csv"))
			pathToMissingDOJ, err := path.Abs(path.Join("test_fixtures", "missing.csv"))

			pathToGogen, err := gexec.Build("gogen_pilots")
			Expect(err).ToNot(HaveOccurred())
			filenameSuffix := "a_suffix"

			runCommand := "run"
			outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
			dojFlag := fmt.Sprintf("--input-doj=%s", pathToValidDOJ+","+pathToBadDOJ+","+pathToMissingDOJ)

			computeAtFlag := "--compute-at=2019-11-11"
			dateTimeFlag := fmt.Sprintf("--file-name-suffix=%s", filenameSuffix)

			command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, computeAtFlag, dateTimeFlag)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(2))
			Eventually(session.Err).Should(gbytes.Say("record on line 2: wrong number of fields"))
			Eventually(session.Err).Should(gbytes.Say("open .*missing.csv: no such file or directory"))

			expectedErrorFileName := fmt.Sprintf("%v/gogen_pilots_%s.err", outputDir, filenameSuffix)

			Ω(expectedErrorFileName).Should(BeAnExistingFile())
			data, _ := ioutil.ReadFile(expectedErrorFileName)
			Expect(string(data)).To(MatchRegexp("open .*missing.csv: no such file or directory"))
			Expect(string(data)).To(MatchRegexp("record on line 2: wrong number of fields"))
		})
	})
})
