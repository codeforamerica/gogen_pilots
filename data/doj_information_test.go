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
	county := "CONTRA COSTA"
	var (
		pathToDOJ      string
		comparisonTime time.Time
		err            error
		dojInformation *DOJInformation
	)

	BeforeEach(func() {
		_, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		inputPath := path.Join("..", "test_fixtures", "contra_costa.xlsx")
		pathToDOJ, _, err = ExtractFullCSVFixtures(inputPath)
		Expect(err).ToNot(HaveOccurred())

		comparisonTime = time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)
		dojInformation, _ = NewDOJInformation(pathToDOJ, comparisonTime, county)
	})

	It("Populates Eligibilities based on Index of Conviction", func() {
		Expect(dojInformation.Eligibilities[11].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
		Expect(dojInformation.Eligibilities[11].CaseNumber).To(Equal("998877; 34345"))
	})

	Context("Computing Aggregate Statistics", func() {
		It("Counts total number of rows in file", func(){
			Expect(dojInformation.TotalRows()).To(Equal(39))
		})
		It("Counts total number of individuals in file", func(){
			Expect(dojInformation.TotalIndividuals()).To(Equal(11))
		})
		It("Counts total convictions", func(){
			Expect(dojInformation.TotalConvictions).To(Equal(29))
		})

		It("Counts total convictions in this county", func(){
			Expect(dojInformation.TotalConvictionsInCounty).To(Equal(26))
		})

		It("Counts all Prop64 convictions sorted by code section", func(){
			Expect(dojInformation.OverallProp64ConvictionsByCodeSection()).To(Equal(map[string]int{"11357": 3, "11358": 10, "11359":5 }))
		})

		It("Counts Prop64 convictions in this county sorted by code section", func(){
			Expect(dojInformation.Prop64ConvictionsInThisCountyByCodeSection(county)).To(Equal(map[string]int{"11357": 3, "11358": 8, "11359":4 }))
		})

		It("Prop64 convictions in this county by code section and eligibility determination", func(){
			Expect(dojInformation.Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county)).To(Equal(map[string]int{"11357": 3, "11358": 8, "11359":4 }))
		})
	})
})
