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

	BeforeEach(func() {
		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "sacramento", "cadoj_sacramento.csv"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("runs and has output", func() {
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

		Expect(sessionString).To(ContainSubstring("Found 24 Total Convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 18 Total Prop64 Convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 3 HS 11357 Convictions total in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 10 HS 11358 Convictions total in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 5 HS 11359 Convictions total in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 0 HS 11360 Convictions total in DOJ file"))

		Expect(sessionString).To(ContainSubstring("Found 21 County Convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 15 County Prop64 Convictions in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 3 HS 11357 Convictions in this county in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 8 HS 11358 Convictions in this county in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 4 HS 11359 Convictions in this county in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 0 HS 11360 Convictions in this county in DOJ file"))

		Expect(sessionString).To(ContainSubstring("Found 7 Prop64 Convictions in this county that are eligible for dismissal in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 2 HS 11357 Convictions in this county that are eligible for dismissal in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 4 HS 11358 Convictions in this county that are eligible for dismissal in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 1 HS 11359 Convictions in this county that are eligible for dismissal in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 0 HS 11360 Convictions in this county that are eligible for dismissal in DOJ file"))

		Expect(sessionString).To(ContainSubstring("Found 7 Prop64 Convictions in this county that are eligible for reduction in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 1 HS 11357 Convictions in this county that are eligible for reduction in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 3 HS 11358 Convictions in this county that are eligible for reduction in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 3 HS 11359 Convictions in this county that are eligible for reduction in DOJ file"))
		Expect(sessionString).To(ContainSubstring("Found 0 HS 11360 Convictions in this county that are eligible for reduction in DOJ file"))

		Expect(sessionString).To(ContainSubstring("Found 1 Prop64 Convictions in this county that are not eligible in DOJ file"))

		Expect(sessionString).To(ContainSubstring("Found 2 Prop64 Convictions in this county with eligibility reason: Misdemeanor or Infraction"))
		Expect(sessionString).To(ContainSubstring("Found 1 Prop64 Convictions in this county with eligibility reason: HS 11357(b)"))
		Expect(sessionString).To(ContainSubstring("Found 3 Prop64 Convictions in this county with eligibility reason: Final Conviction older than 10 years"))
		Expect(sessionString).To(ContainSubstring("Found 5 Prop64 Convictions in this county with eligibility reason: Later Convictions"))
		Expect(sessionString).To(ContainSubstring("Found 1 Prop64 Convictions in this county with eligibility reason: Sentence not Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 Prop64 Convictions in this county with eligibility reason: Sentence Completed"))
		Expect(sessionString).To(ContainSubstring("Found 1 Prop64 Convictions in this county with eligibility reason: Occurred after 11/09/2016"))

		Expect(sessionString).To(ContainSubstring("2 individuals will no longer have a felony on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record"))
		Expect(sessionString).To(ContainSubstring("1 individuals will no longer have any convictions on their record in the last 7 years"))

	})
})
