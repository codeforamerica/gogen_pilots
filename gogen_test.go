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

	PIt("runs and has output", func() {
		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		weightsFlag := fmt.Sprintf("--conviction-weights=%s", pathToWeights)
		cmsFlag := fmt.Sprintf("--input-csv=%s", pathToCMS)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		command := exec.Command(pathToGogen, outputsFlag, weightsFlag, dojFlag, cmsFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session.Out).Should(gbytes.Say("Processed 8 entries"))
		Eventually(session.Out).Should(gbytes.Say("Found 7 felony charges"))

		//TODO expect an outputs csv file to exist
		Expect(true).To(Equal(false))

		Eventually(session).Should(gexec.Exit())
	})
})
