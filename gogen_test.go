package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	path "path/filepath"
	"regexp"

	. "gogen/test_fixtures"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

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

		Eventually(session).Should(gbytes.Say("Found 38 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 11 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 28 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 25 convictions in this county"))
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
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("Found 38 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 11 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 28 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 25 convictions in this county"))
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
		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason 21 years or younger"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason 57 years or older"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Conviction occurred 10 or more years ago"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason Dismiss all HS 11357(c) convictions")))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason Dismiss all HS 11358 convictions"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions with eligibility reason Misdemeanor or Infraction"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Reduce all HS 11359 convictions"))
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

		expectedDojResultsFileName := fmt.Sprintf("%v/doj_results_%s.csv", outputDir, dateSuffix)
		expectedCondensedFileName := fmt.Sprintf("%v/doj_results_condensed_%s.csv", outputDir, dateSuffix)
		expectedConvictionsFileName := fmt.Sprintf("%v/doj_results_convictions_%s.csv", outputDir, dateSuffix)
		expectedOutputFileName := fmt.Sprintf("%v/gogen_%s.out", outputDir, dateSuffix)

		Ω(expectedDojResultsFileName).Should(BeAnExistingFile())
		Ω(expectedCondensedFileName).Should(BeAnExistingFile())
		Ω(expectedConvictionsFileName).Should(BeAnExistingFile())
		Ω(expectedOutputFileName).Should(BeAnExistingFile())
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

		Eventually(session).Should(gexec.Exit(1))
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

	It("runs and has output for Los Angeles", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		runCommand := "run"
		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "LOS ANGELES")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, runCommand, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 32 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Based on your office’s eligibility choices, this application processed the data in .* seconds"))
		Eventually(session).Should(gbytes.Say("Found 9 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 25 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 22 convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions Overall--------------------"))
		Eventually(session).Should(gbytes.Say("Found 18 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 10 11358 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 5 11359 convictions total"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 15 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 8 11358 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 4 11359 convictions in this county"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 5 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 8 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions total that are Hand Review"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 2 11357 convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions total that are Maybe Eligible - Flag for Review"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 2 11358 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions total that are Not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason 11357(a) or 11357(b)")))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason 21 years or younger"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason 50 years or older"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Only has 11357-60 charges and completed sentence"))

		Eventually(session).Should(gbytes.Say("Hand Review"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Currently serving sentence"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason Other 11357"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 2 convictions with eligibility reason PC 667(e)(2)(c)(iv)")))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Two priors"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 1 466 PC convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals currently have convictions on their record in the last 7 years"))

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
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("&&&&&&"))
		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 36 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 11 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 28 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 25 convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions Overall--------------------"))
		Eventually(session).Should(gbytes.Say("Found 20 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 4 11357 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 8 11358 convictions total"))
		Eventually(session).Should(gbytes.Say("Found 8 11359 convictions total"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 17 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 4 11357 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 6 11358 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 7 11359 convictions in this county"))
		Eventually(session).Should(gbytes.Say("Date of earliest Prop 64 conviction: June 1979"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 6 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 6 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 16 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions total that are Eligible for Reduction"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason 21 years or younger"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason 57 years or older"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Conviction occurred 10 or more years ago"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason Dismiss all HS 11357 convictions (when no sub-section is specified)")))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason Dismiss all HS 11357(c) convictions")))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason Dismiss all HS 11358 convictions"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Individual is deceased"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions with eligibility reason Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Only has 11357-60 charges"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Reduce all HS 11359 convictions"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 0 convictions in this county"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("11 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("11 individuals currently have convictions on their record"))
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
})
