package main_test

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"io/ioutil"
	"os/exec"

	path "path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ bool = Describe("gogen", func() {
	var (
		outputDir     string
		pathToWeights string
		pathToDOJ     string
		pathToCMS     string
		err           error
	)

	BeforeEach(func() {
		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		pathToWeights, err = path.Abs(path.Join("test_fixtures", "conviction_weights.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "cadoj.csv"))
		Expect(err).ToNot(HaveOccurred())

		pathToCMS, err = path.Abs(path.Join("test_fixtures", "felonies_sf.csv"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("runs and has output", func() {
		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		weightsFlag := fmt.Sprintf("--conviction-weights=%s", pathToWeights)
		cmsFlag := fmt.Sprintf("--input-csv=%s", pathToCMS)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		command := exec.Command(pathToGogen, outputsFlag, weightsFlag, dojFlag, cmsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session.Out).Should(gbytes.Say(`Found 9 charges in CMS data \(8 felonies, 1 misdemeanors\)`))
		//Eventually(session.Out).Should(gbytes.Say(`Found 11 charges in DOJ data \(11 felony and 0 misdemeanors\)`))
		//Eventually(session.Out).Should(gbytes.Say(`Failed to match 2 out of 9 charges in CMS data`))
		//Eventually(session.Out).Should(gbytes.Say(`Failed to match # out of # charges in DOJ data`))
		//Eventually(session.Out).Should(gbytes.Say(`Failed to match #  out of # unique subjects in DOJ data`))
		//Eventually(session.Out).Should(gbytes.Say(`Match details: ...`))

		pathToExpectedResults, err := path.Abs(path.Join("test_fixtures", "felonies_sf_results.csv"))
		Expect(err).ToNot(HaveOccurred())
		//expectedResultsBody, err := ioutil.ReadFile(pathToExpectedResults)
		_, err = ioutil.ReadFile(pathToExpectedResults)
		Expect(err).ToNot(HaveOccurred())

		pathToOutput, err := path.Abs(path.Join(outputDir, "results.csv"))
		Expect(err).ToNot(HaveOccurred())
		//outputBody, err := ioutil.ReadFile(pathToOutput)
		_, err = ioutil.ReadFile(pathToOutput)
		Expect(err).ToNot(HaveOccurred())

		format.TruncatedDiff = false
		//Expect(string(outputBody)).To(Equal(string(expectedResultsBody)))

		Eventually(session).Should(gexec.Exit())
	})
})
