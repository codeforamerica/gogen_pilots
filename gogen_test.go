package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	path "path/filepath"

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

		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("Found 38 Total rows in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 11 Total individuals in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 28 Total convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 25 convictions in this county"))
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

		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("Found 38 Total rows in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 11 Total individuals in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 28 Total convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 25 convictions in this county"))
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
		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("----------- Eligibility Reasons --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason HS 11357(b)"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions with eligibility reason Has convictions in past 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions with eligibility reason No convictions in past 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence not Completed"))
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
		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("----------- Overall summary of DOJ file --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 32 Total rows in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 9 Total individuals in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 25 Total convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 22 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions Overall--------------------"))
		Expect(sessionString).To(ContainSubstring("Found 18 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 10 11358 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 5 11359 convictions total"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 15 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 8 11358 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 4 11359 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 9 convictions total that are Eligible for Dismissal"))

		Expect(sessionString).To(ContainSubstring("Eligible for Reduction"))
		Expect(sessionString).To(ContainSubstring("Found 3 11358 convictions that are Eligible for Reduction"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are Eligible for Reduction"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions total that are Eligible for Reduction"))

		Expect(sessionString).To(ContainSubstring("Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 11358 convictions that are Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions total that are Not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Eligibility Reasons --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason HS 11357(b)"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions with eligibility reason Has convictions in past 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions with eligibility reason No convictions in past 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence not Completed"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Related Convictions In This County --------------------"))

		Expect(sessionString).To(ContainSubstring("Found 0 convictions in this county")) //Sacramento did not include related charges

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("4 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individual(s) who had a felony will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("1 individual(s) who had convictions will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("1 individual(s) who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))
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

		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))
		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("----------- Overall summary of DOJ file --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 41 Total rows in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 12 Total individuals in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 31 Total convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 28 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions Overall--------------------"))
		Expect(sessionString).To(ContainSubstring("Found 19 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 11 11358 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 5 11359 convictions total"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 16 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 9 11358 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 4 11359 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 9 convictions total that are Eligible for Dismissal"))

		Expect(sessionString).To(ContainSubstring("Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 6 convictions total that are Maybe Eligible - Flag for Review"))

		Expect(sessionString).To(ContainSubstring("Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 11358 convictions that are Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions total that are Not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Eligibility Reasons --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason 11357 HS"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions with eligibility reason Has convictions in past 5 years"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions with eligibility reason No convictions in past 5 years"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence not Completed"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Related Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county")) //until we start ignoring charges with no matching prop 64 charge
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 2 4149 BP convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions total that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are Eligible for Dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions total that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are Maybe Eligible - Flag for Review"))

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("10 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("12 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("7 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Expect(sessionString).To(ContainSubstring("1 individual(s) who had a felony will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("3 individual(s) who had convictions will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("2 individual(s) who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))
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
		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("----------- Overall summary of DOJ file --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 39 Total rows in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 11 Total individuals in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 29 Total convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 26 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions Overall--------------------"))
		Expect(sessionString).To(ContainSubstring("Found 18 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 10 11358 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 5 11359 convictions total"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 15 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 8 11358 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 4 11359 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11359 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 10 convictions total that are Eligible for Dismissal"))

		Expect(sessionString).To(ContainSubstring("Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 3 11358 convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 11359 convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions total that are Maybe Eligible - Flag for Review"))

		Expect(sessionString).To(ContainSubstring("Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 11358 convictions that are Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions total that are Not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Eligibility Reasons --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason 11357 HS"))
		Expect(sessionString).To(ContainSubstring("Found 3 convictions with eligibility reason Has convictions in past 5 years"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions with eligibility reason No convictions in past 5 years"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Sentence not Completed"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Related Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county")) //until we start ignoring charges with no matching prop 64 charge
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 2 4149 BP convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions total that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are Eligible for Dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions total that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are Maybe Eligible - Flag for Review"))

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("11 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("6 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Expect(sessionString).To(ContainSubstring("1 individual(s) who had a felony will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("3 individual(s) who had convictions will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("2 individual(s) who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))
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
		sessionString := string(session.Out.Contents())

		Expect(sessionString).To(ContainSubstring("----------- Overall summary of DOJ file --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 32 Total rows in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 9 Total individuals in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 25 Total convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 22 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions Overall--------------------"))
		Expect(sessionString).To(ContainSubstring("Found 18 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 10 11358 convictions total"))
		Expect(sessionString).To(ContainSubstring("Found 5 11359 convictions total"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 15 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 8 11358 convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 4 11359 convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 11357 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 5 11358 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 8 convictions total that are Eligible for Dismissal"))

		Expect(sessionString).To(ContainSubstring("Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 2 11357 convictions that are Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions total that are Maybe Eligible - Flag for Review"))

		Expect(sessionString).To(ContainSubstring("Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 2 11358 convictions that are Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions total that are Not eligible"))

		Expect(sessionString).To(ContainSubstring("Hand Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 11358 convictions that are Hand Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions total that are Hand Review"))

		Expect(sessionString).To(ContainSubstring("----------- Eligibility Reasons --------------------"))

		Expect(sessionString).To(ContainSubstring("Hand Review"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Currently serving sentence"))

		Expect(sessionString).To(ContainSubstring("Not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason PC 667(e)(2)(c)(iv)"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Two priors"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Occurred after 11/09/2016"))

		Expect(sessionString).To(ContainSubstring("Maybe Eligible - Flag for Review"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions with eligibility reason Other 11357"))

		Expect(sessionString).To(ContainSubstring("Eligible for Dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason Only has 11357-60 charges and completed sentence"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason 21 years or younger"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions with eligibility reason 50 years or older"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions with eligibility reason 11357(a) or 11357(b)"))

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("4 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If eligibility is ran as is for Prop 64 and Related Charges --------------------"))
		Expect(sessionString).To(ContainSubstring("1 individual(s) who had a felony will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("1 individual(s) who had convictions will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("0 individual(s) who had convictions in the last 7 years will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have a felony on their record")) // VM - this changed from 2 to 3
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))
	})

})
