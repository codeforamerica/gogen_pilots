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

		pathToDOJ, err = path.Abs(path.Join("test_fixtures", "cadoj.csv"))
		Expect(err).ToNot(HaveOccurred())
	})

	It("runs and has output", func() {
		pathToGogen, err := gexec.Build("gogen")
		Expect(err).ToNot(HaveOccurred())

		outputsFlag := fmt.Sprintf("--outputs=%s", outputDir)
		dojFlag := fmt.Sprintf("--input-doj=%s", pathToDOJ)
		countyFlag := fmt.Sprintf("--county=%s", "SAN FRANCISCO")
		command := exec.Command(pathToGogen, outputsFlag, dojFlag, countyFlag)
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		//Eventually(session.Out).Should(gbytes.Say(`Found 8 convictions in DOJ data \(8 felonies, 0 misdemeanors\)`))
		Eventually(session.Out).Should(gbytes.Say(`Found 24 Total Convictions in DOJ file`))
		Eventually(session.Out).Should(gbytes.Say(`Found 21 SAN FRANCISCO County Convictions in DOJ file`))
		Eventually(session.Out).Should(gbytes.Say(`Found 18 SAN FRANCISCO County Prop64 Convictions in DOJ file`))

		Eventually(session.Out).Should(gbytes.Say(`Found 7 Prop64 Convictions in this county that are eligible for dismissal in DOJ file`))
		Eventually(session.Out).Should(gbytes.Say(`Found 7 Prop64 Convictions in this county that are eligible for reduction in DOJ file`))
		Eventually(session.Out).Should(gbytes.Say(`Found 1 Prop64 Convictions in this county that are not eligible in DOJ file`))

		Eventually(session.Out).Should(gbytes.Say(`Found 2 Prop64 Convictions in this county that are eligible for dismissal in DOJ file because of Misdemeanor or Infraction`))
		Eventually(session.Out).Should(gbytes.Say(`Found 1 Prop64 Convictions in this county that are eligible for dismissal in DOJ file because of HS 11357b`))
		Eventually(session.Out).Should(gbytes.Say(`Found 3 Prop64 Convictions in this county that are eligible for dismissal in DOJ file because final conviction older than 10 years`))
		Eventually(session.Out).Should(gbytes.Say(`Found 5 Prop64 Convictions in this county that are eligible for reduction in DOJ file because there are later convictions`))
		Eventually(session.Out).Should(gbytes.Say(`Found 1 Prop64 Convictions in this county that are eligible for reduction in DOJ file because they did not complete their sentence`))
		Eventually(session.Out).Should(gbytes.Say(`Found 1 Prop64 Convictions in this county that are eligible for dismissal in DOJ file because they completed their sentence`))
		Eventually(session.Out).Should(gbytes.Say(`Found 1 Prop64 Convictions in this county that are not eligible because after November 9 2016`))

		Eventually(session.Out).Should(gbytes.Say(`Found 2 Prop64 Convictions in this county that need sentence data checked`))

		Eventually(session.Out).Should(gbytes.Say(`2 individuals will no longer have a felony on their record`))
		Eventually(session.Out).Should(gbytes.Say(`1 individuals will no longer have any convictions on their record`))
		Eventually(session.Out).Should(gbytes.Say(`6 individuals will no longer have any convictions on their record in the last 7 years`))

		Eventually(session).Should(gexec.Exit())
		Expect(session.Err).ToNot(gbytes.Say("required"))
	})
})
