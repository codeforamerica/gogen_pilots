package main_test

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	path "path/filepath"

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

	It("runs and has output for Sacramento", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "sacramento", "cadoj_sacramento.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SACRAMENTO")
		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
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

		Eventually(session).Should(gbytes.Say("Found 9 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 5 convictions that are eligible for reduction"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are eligible for reduction"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are eligible for reduction"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 1 convictions that are not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: HS 11357\\(b\\)"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county with eligibility reason: Has convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions in this county with eligibility reason: No convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Occurred after 11/09/2016"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence Completed"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions in this county")) //Sacramento did not include related charges

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record in the last 7 years"))

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

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "san_joaquin", "cadoj_san_joaquin.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
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

		Eventually(session).Should(gbytes.Say("Found 9 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are eligible for reduction"))

		Eventually(session).Should(gbytes.Say("Found 6 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 1 convictions that are not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: 11357 HS"))
		Eventually(session).Should(gbytes.Say("Found 5 convictions in this county with eligibility reason: Has convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county with eligibility reason: No convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Occurred after 11/09/2016"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence Completed"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county")) //until we start ignoring charges with no matching prop 64 charge
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 2 4149 BP convictions in this county"))

		Eventually(session).Should(gbytes.Say("Found 2 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 2 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("10 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("12 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("7 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("runs and has output for Contra Costa", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "contra_costa", "cadoj_contra_costa.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "CONTRA COSTA")
		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
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

		Eventually(session).Should(gbytes.Say("Found 9 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 3 11357 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 4 11358 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are eligible for reduction"))

		Eventually(session).Should(gbytes.Say("Found 5 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 1 convictions that are not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11358 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: 11357 HS"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county with eligibility reason: Has convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county with eligibility reason: No convictions in past 5 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Occurred after 11/09/2016"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence Completed"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))
		Eventually(session).Should(gbytes.Say("Found 4 convictions in this county")) //until we start ignoring charges with no matching prop 64 charge
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions in this county"))
		Eventually(session).Should(gbytes.Say("Found 2 4149 BP convictions in this county"))

		Eventually(session).Should(gbytes.Say("Found 2 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 148 PC convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 2 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 1 4060 BP convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 1 4149 BP convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("11 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("6 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("runs and has output for Los Angeles", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "los_angeles", "cadoj_los_angeles.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "LOS ANGELES")
		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
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

		Eventually(session).Should(gbytes.Say("Found 6 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 1 11357 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are eligible for dismissal"))
		Eventually(session).Should(gbytes.Say("Found 2 11359 convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 3 convictions that are eligible for reduction"))
		Eventually(session).Should(gbytes.Say("Found 2 11358 convictions that are eligible for reduction"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are eligible for reduction"))

		Eventually(session).Should(gbytes.Say("Found 2 convictions that are flagged for review"))
		Eventually(session).Should(gbytes.Say("Found 2 11357 convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 4 convictions that are not eligible"))
		Eventually(session).Should(gbytes.Say("Found 3 11358 convictions that are not eligible"))
		Eventually(session).Should(gbytes.Say("Found 1 11359 convictions that are not eligible"))


		Eventually(session).Should(gbytes.Say("----------- Eligibility Reasons --------------------"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: 11357\\(a\\) or 11357\\(b\\)"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: Has convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions in this county with eligibility reason: No convictions in past 10 years"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Occurred after 11/09/2016"))
		Eventually(session).Should(gbytes.Say("Found 2 convictions in this county with eligibility reason: Other 11357"))
		Eventually(session).Should(gbytes.Say("Found 3 convictions in this county with eligibility reason: PC 667\\(e\\)\\(2\\)\\(c\\)\\(iv\\)"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence Completed"))
		Eventually(session).Should(gbytes.Say("Found 1 convictions in this county with eligibility reason: Sentence not Completed"))

		Eventually(session).Should(gbytes.Say("----------- Prop64 Related Convictions In This County --------------------"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions in this county")) //Los Angeles did not include related charges

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are eligible for dismissal"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are flagged for review"))

		Eventually(session).Should(gbytes.Say("Found 0 convictions that are not eligible"))

		Eventually(session).Should(gbytes.Say("----------- Impact to individuals --------------------"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have a felony on their record"))
		Eventually(session).Should(gbytes.Say("9 individuals currently have convictions on their record"))
		Eventually(session).Should(gbytes.Say("4 individuals currently have convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have a felony on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("1 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record")) // VM - this changed from 2 to 3
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Eventually(session).Should(gbytes.Say("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have a felony on their record")) // VM - this changed from 2 to 3
		Eventually(session).Should(gbytes.Say("2 individuals will no longer have any convictions on their record"))
		Eventually(session).Should(gbytes.Say("3 individuals will no longer have any convictions on their record in the last 7 years"))
	})

	It("can handle a csv with extra comma at the end of headers", func() {

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "san_joaquin", "cadoj_san_joaquin_extra_comma.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN JOAQUIN")
		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
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
})
