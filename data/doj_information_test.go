package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gogen/data"
	. "gogen/test_fixtures"
	"io/ioutil"
	"path"
	"time"
)

var _ = Describe("DojInformation", func() {
	var (
		pathToDOJ      string
		comparisonTime time.Time
		err            error
		dojInformation *DOJInformation
	)

	BeforeEach(func() {
		_, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		inputPath := path.Join("..", "test_fixtures", "contra_costa", "cadoj_contra_costa_source.xlsx")
		pathToDOJ, _, err = ExtractFullCSVFixtures(inputPath)
		Expect(err).ToNot(HaveOccurred())

		comparisonTime = time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)
		dojInformation, _ = NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA")
	})

	It("Populates Eligibilities based on Index of Conviction", func() {
		Expect(dojInformation.Eligibilities[11].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
		Expect(dojInformation.Eligibilities[11].CaseNumber).To(Equal("998877; 34345"))
	})

	Context("Computing Aggregate Statistics", func() {
		It("Counts total convictions", func(){
			Expect(dojInformation.TotalConvictions).To(Equal(29))
		})

		It("Counts total convictions in this county", func(){
			Expect(dojInformation.TotalConvictionsInCounty).To(Equal(26))
		})

	})
})
