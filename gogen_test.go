package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	path "path/filepath"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "gogen/test_fixtures"
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

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
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

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("Found 38 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 11 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 28 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 25 convictions in this county"))
	})

	It("can accept a compute-at option for determining eligibility", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "sacramento.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))
		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason HS 11357(b)")))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason No convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence Completed"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions with eligibility reason Has convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))
	})

	It("runs and has output for Sacramento", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "sacramento.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 32 Total rows in DOJ file"))
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
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 9 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions total that are Eligible for Reduction"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions total that are Not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason HS 11357(b)")))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason No convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence Completed"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions with eligibility reason Has convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions in this county")) //Sacramento did not include related charges

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("runs and has output for San Joaquin", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "san_joaquin.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		computeAtDate := "--compute-at=2019-05-01"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtDate)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 41 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 12 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 31 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 28 convictions in this county"))

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
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 9 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 6 convictions total that are Maybe Eligible - Flag for Review"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions total that are Not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason 11357 HS"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions with eligibility reason No convictions in past 5 years"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason Has convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county")) //until we start ignoring charges with no matching prop 64 charge
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 2 4149 BP convictions in this county"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions total that are Maybe Eligible - Flag for Review"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("10 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("12 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("7 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("1 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("5 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("runs and has output for Contra Costa", func() {
		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "contra_costa.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "CONTRA COSTA")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 39 Total rows in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 11 Total individuals in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 29 Total convictions in DOJ file"))
		Eventually(session).Should(gbytes.Say("Found 26 convictions in this county"))

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
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 3 11359 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 10 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions total that are Maybe Eligible - Flag for Review"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions total that are Not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason 11357 HS"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions with eligibility reason No convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence Completed"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions with eligibility reason Has convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("Not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county")) //until we start ignoring charges with no matching prop 64 charge
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 2 4149 BP convictions in this county"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are Maybe Eligible - Flag for Review"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions total that are Maybe Eligible - Flag for Review"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("11 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("6 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("1 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("5 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("runs and has output for Los Angeles", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "los_angeles.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "LOS ANGELES")
		computeAtFlag := "--compute-at=2019-11-11"

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 32 Total rows in DOJ file"))
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

		Eventually(session).Should(gbytes.Say("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("1 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("0 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("can accept path to eligibility options file", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToInputExcel := path.Join("test_fixtures", "configurable_flow.xlsx")
		inputCSV, _, _ := ExtractFullCSVFixtures(pathToInputExcel)

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToEligibilityOptions := path.Join("test_fixtures", "eligibility_options.json")

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", inputCSV)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		computeAtFlag := "--compute-at=2019-11-11"
		eligibilityOptionsFlag := fmt.Sprintf("--eligibility-options=%s", pathToEligibilityOptions)

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag, computeAtFlag, eligibilityOptionsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))

		Eventually(session).Should(gbytes.Say("----------- Overall summary of DOJ file --------------------"))
		Eventually(session).Should(gbytes.Say("Found 32 Total rows in DOJ file"))
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
		Eventually(session).Should(gbytes.Say("Found 2 11357 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 8 11358 convictions that are Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say("Found 10 convictions total that are Eligible for Dismissal"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 11357 convictions that are Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 4 11359 convictions that are Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions total that are Eligible for Reduction"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))

		Eventually(session).Should(gbytes.Say("Eligible for Dismissal"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 2 convictions with eligibility reason Dismiss all 11357(C)HS convictions")))
		Eventually(session).Should(gbytes.Say("Found 8 convictions with eligibility reason Dismiss all 11358 HS convictions"))

		Eventually(session).Should(gbytes.Say("Eligible for Reduction"))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason Reduce all 11357(B)HS convictions")))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 2 convictions with eligibility reason Reduce all 11359 HS convictions")))
		Eventually(session).Should(gbytes.Say(regexp.QuoteMeta("Found 1 convictions with eligibility reason Reduce all 11359(C) HS convictions")))
		Eventually(session).Should(gbytes.Say("Found 1 convictions with eligibility reason Reduce all 11359HS convictions"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals who had a felony will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals who had convictions will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))
	})

})
