package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gogen/data"
	"gogen/matchers"
	. "gogen/test_fixtures"

	"io/ioutil"
	"path"
	"regexp"
	"time"
)

type testEligibilityFlow struct {
	prop64Matcher *regexp.Regexp
}

func (ef testEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	ef.EligibleDismissal(info, "Because")
}

func (ef testEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Test is Eligible"
	info.EligibilityReason = reason
}

func (ef testEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, "SACRAMENTO")
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef testEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return ef.IsProp64Charge(codeSection)
}

func (ef testEligibilityFlow) IsProp64Charge(codeSection string) bool {
	return matchers.IsProp64Charge(codeSection)
}

func (ef testEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

var _ = Describe("DojInformation", func() {
	county := "CONTRA COSTA"
	var (
		pathToDOJ          string
		comparisonTime     time.Time
		err                error
		testDojInformation *DOJInformation
		dojInformation     *DOJInformation
	)

	BeforeEach(func() {
		_, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())

		inputPath := path.Join("..", "test_fixtures", "contra_costa.xlsx")
		pathToDOJ, _, err = ExtractFullCSVFixtures(inputPath)
		Expect(err).ToNot(HaveOccurred())

		comparisonTime = time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)
		testFlow := testEligibilityFlow{prop64Matcher: regexp.MustCompile(`(11357|11358|11359|11360)`)}
		contraCostaFlow := EligibilityFlows["CONTRA COSTA"]
		testDojInformation = NewDOJInformation(pathToDOJ, comparisonTime, county, testFlow)
		dojInformation = NewDOJInformation(pathToDOJ, comparisonTime, county, contraCostaFlow)
	})

	It("Uses the provided eligibility flow", func() {
		Expect(testDojInformation.Eligibilities[11].EligibilityDetermination).To(Equal("Test is Eligible"))
		Expect(testDojInformation.Eligibilities[11].CaseNumber).To(Equal("998877; 34345"))
	})

	It("Populates Eligibilities based on Index of Conviction", func() {
		Expect(dojInformation.Eligibilities[11].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
		Expect(dojInformation.Eligibilities[11].CaseNumber).To(Equal("998877; 34345"))
	})

	Context("Computing Aggregate Statistics for convictions", func() {
		It("Counts total number of rows in file", func() {
			Expect(dojInformation.TotalRows()).To(Equal(39))
		})

		It("Counts total convictions", func() {
			Expect(dojInformation.TotalConvictions).To(Equal(29))
		})

		It("Counts total convictions in this county", func() {
			Expect(dojInformation.TotalConvictionsInCounty).To(Equal(26))
		})

		It("Counts all Prop64 convictions sorted by code section", func() {
			Expect(dojInformation.OverallProp64ConvictionsByCodeSection()).To(Equal(map[string]int{"11357": 3, "11358": 10, "11359": 5}))
		})

		It("Counts Prop64 convictions in this county sorted by code section", func() {
			Expect(dojInformation.Prop64ConvictionsInThisCountyByCodeSection(county)).To(Equal(map[string]int{"11357": 3, "11358": 8, "11359": 4}))
		})

		It("Prop64 convictions in this county by code section and eligibility determination", func() {
			Expect(dojInformation.Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county)).To(Equal(
				map[string]map[string]int{
					"Eligible for Dismissal":           {"11357": 3, "11358": 4, "11359": 3},
					"Maybe Eligible - Flag for Review": {"11358": 3, "11359": 1},
					"Not eligible":                     {"11358": 1}}))
		})

		It("Prop64 convictions in this county by eligibility determination and reason", func() {
			Expect(dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county)).To(Equal(
				map[string]map[string]int{
					"Eligible for Dismissal": {
						"No convictions in past 5 years": 5,
						"11357 HS":                       2,
						"Misdemeanor or Infraction":      2,
						"Sentence Completed":             1,
					},
					"Maybe Eligible - Flag for Review": {
						"Has convictions in past 5 years": 3,
						"Sentence not Completed":          1},
					"Not eligible": {"Occurred after 11/09/2016": 1}}))
		})

		It("Related convictions in this county by code section and eligibility determination", func() {
			Expect(dojInformation.RelatedConvictionsInThisCountyByCodeSectionByEligibility(county)).To(Equal(
				map[string]map[string]int{
					"Eligible for Dismissal":           {"4149 BP": 1, "148 PC": 1},
					"Maybe Eligible - Flag for Review": {"4149 BP": 1, "4060 BP": 1}}))
		})

		Context("Computing aggregate statistics for individuals", func() {
			It("Counts total number of individuals in file", func() {
				Expect(dojInformation.TotalIndividuals()).To(Equal(11))
			})

			Context("Before eligibility is run", func() {
				It("Calculates individuals with a felony", func() {
					Expect(dojInformation.CountIndividualsWithFelony()).To(Equal(9))
				})

				It("Calculates individuals with any conviction", func() {
					Expect(dojInformation.CountIndividualsWithConviction()).To(Equal(11))
				})

				It("Calculates individuals with any conviction in the last 7 years", func() {
					Expect(dojInformation.CountIndividualsWithConvictionInLast7Years()).To(Equal(6))
				})
			})

			Context("After eligibility is run", func() {
				It("Calculates individuals who will no longer have a felony ", func() {
					Expect(dojInformation.CountIndividualsNoLongerHaveFelony()).To(Equal(1))
				})

				It("Calculates individuals who no longer have any conviction", func() {
					Expect(dojInformation.CountIndividualsNoLongerHaveConviction()).To(Equal(3))
				})

				It("Calculates individuals who no longer have any conviction in the last 7 years", func() {
					Expect(dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years()).To(Equal(2))
				})
			})

		})
	})
})
