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

		sessionString := string(session.Out.Contents())

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

		Expect(sessionString).To(ContainSubstring("Found 9 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are eligible for dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 5 convictions that are eligible for reduction"))
		Expect(sessionString).To(ContainSubstring("Found 3 11358 convictions that are eligible for reduction"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are eligible for reduction"))

		Expect(sessionString).To(ContainSubstring("Found 1 convictions that are not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 11358 convictions that are not eligible"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: HS 11357(b)"))
		Expect(sessionString).To(ContainSubstring("Found 5 convictions in this county with eligibility reason: No convictions in past 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county with eligibility reason: Has convictions in past 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: Sentence not Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: Occurred after 11/09/2016"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Related Convictions In This County --------------------"))

		Expect(sessionString).To(ContainSubstring("Found 0 convictions in this county"))//Sacramento did not include related charges

		Expect(sessionString).To(ContainSubstring("Found 0 convictions that are eligible for dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 0 convictions that are flagged for review"))

		Expect(sessionString).To(ContainSubstring("Found 0 convictions that are not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("4 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record in the last 7 years"))

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

		sessionString := string(session.Out.Contents())

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

		Expect(sessionString).To(ContainSubstring("Found 9 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are eligible for dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 6 convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are flagged for review"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Related Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county"))//until we start ignoring charges with no matching prop 64 charge
		Expect(sessionString).To(ContainSubstring("Found 2 4149 BP convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are eligible for dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions that are flagged for review"))

		Expect(sessionString).To(ContainSubstring("Found 0 convictions that are not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("10 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("12 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("7 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("4 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("4 individuals will no longer have any convictions on their record in the last 7 years"))
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

		Expect(sessionString).To(ContainSubstring("Found 9 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 3 11357 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 4 11358 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are eligible for dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 5 convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 3 11358 convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 2 11359 convictions that are flagged for review"))

		Expect(sessionString).To(ContainSubstring("Found 1 convictions that are not eligible"))
		Expect(sessionString).To(ContainSubstring("Found 1 11358 convictions that are not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Eligibility Reasons --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county with eligibility reason: No convictions in past 5 years"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county with eligibility reason: Has convictions in past 5 years"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: Sentence not Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 2 convictions in this county with eligibility reason: 11357 HS"))
		Expect(sessionString).To(ContainSubstring("Found 1 convictions in this county with eligibility reason: Occurred after 11/09/2016"))

		Expect(sessionString).To(ContainSubstring("----------- Prop64 Related Convictions In This County --------------------"))
		Expect(sessionString).To(ContainSubstring("Found 4 convictions in this county"))//until we start ignoring charges with no matching prop 64 charge
		Expect(sessionString).To(ContainSubstring("Found 2 4149 BP convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions in this county"))
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions in this county"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 148 PC convictions that are eligible for dismissal"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are eligible for dismissal"))

		Expect(sessionString).To(ContainSubstring("Found 2 convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4149 BP convictions that are flagged for review"))
		Expect(sessionString).To(ContainSubstring("Found 1 4060 BP convictions that are flagged for review"))

		Expect(sessionString).To(ContainSubstring("Found 0 convictions that are not eligible"))

		Expect(sessionString).To(ContainSubstring("----------- Impact to individuals --------------------"))
		Expect(sessionString).To(ContainSubstring("9 individuals currently have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("11 individuals currently have convictions on their record"))
		Expect(sessionString).To(ContainSubstring("6 individuals currently have convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ELIGIBLE Prop 64 convictions are dismissed or reduced --------------------"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If ALL Prop 64 convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have any convictions on their record in the last 7 years"))

		Expect(sessionString).To(ContainSubstring("----------- If all Prop 64 AND related convictions are dismissed and sealed --------------------"))
		Expect(sessionString).To(ContainSubstring("3 individuals will no longer have a felony on their record\n4 individuals"))//to differentiate from the same content in Prop64 only stats
		Expect(sessionString).To(ContainSubstring("4 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("4 individuals will no longer have any convictions on their record in the last 7 years"))
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
