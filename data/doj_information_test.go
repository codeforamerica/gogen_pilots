package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gogen_pilots/data"
	. "gogen_pilots/data"
	"gogen_pilots/matchers"
	. "gogen_pilots/test_fixtures"

	"io/ioutil"
	"path"
	"time"
)

type testEligibilityFlow struct {
}

func (ef testEligibilityFlow) BeginEligibilityFlow(info *EligibilityInfo, row *DOJRow, subject *Subject) {
	ef.EligibleDismissal(info, "Because")
}

func (ef testEligibilityFlow) EligibleDismissal(info *EligibilityInfo, reason string) {
	info.EligibilityDetermination = "Test is Eligible"
	info.EligibilityReason = reason
}

func (ef testEligibilityFlow) ProcessSubject(subject *Subject, comparisonTime time.Time, county string, age int, yearsConvictionFree int) map[int]*EligibilityInfo {
	infos := make(map[int]*EligibilityInfo)
	for _, conviction := range subject.Convictions {
		if ef.checkRelevancy(conviction.CodeSection, conviction.County) {
			info := NewEligibilityInfo(conviction, subject, comparisonTime, "LOS ANGELES")
			ef.BeginEligibilityFlow(info, conviction, subject)
			infos[conviction.Index] = info
		}
	}
	return infos
}

func (ef testEligibilityFlow) checkRelevancy(codeSection string, county string) bool {
	return matchers.IsProp64Charge(codeSection)
}

func (ef testEligibilityFlow) ChecksRelatedCharges() bool {
	return false
}

var _ = Describe("DojInformation", func() {
	county := "LOS ANGELES"
	var (
		pathToDOJ         string
		comparisonTime    time.Time
		err               error
		testEligibilities map[int]*EligibilityInfo
		dojInformation    *DOJInformation
		dojEligibilities  map[int]*EligibilityInfo
	)

	BeforeEach(func() {
		_, err = ioutil.TempDir("/tmp", "gogen_pilots")
		Expect(err).ToNot(HaveOccurred())

		inputPath := path.Join("..", "test_fixtures", "los_angeles.xlsx")
		pathToDOJ, _, err = ExtractFullCSVFixtures(inputPath)
		Expect(err).ToNot(HaveOccurred())

		comparisonTime = time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)
		testFlow := testEligibilityFlow{}
		losAngelesFlow := data.EligibilityFlows["LOS ANGELES"]
		dojInformation, _ = NewDOJInformation(pathToDOJ, comparisonTime, losAngelesFlow)
		age := 50
		yearsConvictionFree := 10

		testEligibilities = dojInformation.DetermineEligibility(county, testFlow, age, yearsConvictionFree)
		dojEligibilities = dojInformation.DetermineEligibility(county, losAngelesFlow, age, yearsConvictionFree)

	})

	It("Uses the provided eligibility flow", func() {
		Expect(testEligibilities[11].EligibilityDetermination).To(Equal("Test is Eligible"))
		Expect(testEligibilities[11].CaseNumber).To(Equal("140194; 140195"))
	})

	It("Populates ToBeRemovedEligibilities based on Index of Conviction", func() {
		Expect(dojEligibilities[11].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
		Expect(dojEligibilities[11].CaseNumber).To(Equal("140194; 140195"))
	})

	Context("Computing Aggregate Statistics for convictions", func() {
		It("Counts total number of rows in file", func() {
			Expect(dojInformation.TotalRows()).To(Equal(35))
		})

		It("Counts total convictions", func() {
			Expect(dojInformation.TotalConvictions()).To(Equal(28))
		})

		It("Counts total convictions in this county", func() {
			Expect(dojInformation.TotalConvictionsInCounty(county)).To(Equal(25))
		})

		It("Counts all Prop64 convictions sorted by code section", func() {
			Expect(dojInformation.OverallProp64ConvictionsByCodeSection()).To(Equal(map[string]int{"11357": 3, "11358": 11, "11359": 5}))
		})

		It("Counts Prop64 convictions in this county sorted by code section", func() {
			Expect(dojInformation.Prop64ConvictionsInThisCountyByCodeSection(county)).To(Equal(map[string]int{"11357": 3, "11358": 9, "11359": 4}))
		})

		It("Finds the date of the earliest Prop64 conviction in the county", func() {
			expectedDate := time.Date(1979, 6, 1, 0, 0, 0, 0, time.UTC)
			Expect(dojInformation.EarliestProp64ConvictionDateInThisCounty(county)).To(Equal(expectedDate))
		})

		It("Prop64 convictions in this county by code section and eligibility determination", func() {
			Expect(dojInformation.Prop64ConvictionsInThisCountyByCodeSectionByEligibility(county, dojEligibilities)).To(Equal(
				map[string]map[string]int{
					"Eligible for Dismissal": {"11357": 2, "11359": 2, "11358": 5},
					"Hand Review": {"11358": 2},
					"Not eligible": {"11358": 1, "11359": 1},
					"To be reviewed by City Attorneys": {"11358": 1, "11357": 1, "11359": 1}}))
		})

		It("Prop64 convictions in this county by eligibility determination and reason", func() {
			Expect(dojInformation.Prop64ConvictionsInThisCountyByEligibilityByReason(county, dojEligibilities)).To(Equal(
				map[string]map[string]int{
					"Not eligible": {
						"PC 667(e)(2)(c)(iv)": 2,
					},
					"To be reviewed by City Attorneys": {
						"Misdemeanor or Infraction": 3,
					},
					"Hand Review": {
						"Currently serving sentence": 1,
						"No applicable eligibility criteria": 1,
					},
					"Eligible for Dismissal": {
						"11357": 2,
						"50 years or older": 5,
						"21 years or younger": 1,
						"Only has 11357-60 charges and completed sentence": 1,
					},
				}))
		})

		It("Related convictions in this county by code section and eligibility determination", func() {
			Expect(dojInformation.RelatedConvictionsInThisCountyByCodeSectionByEligibility(county, dojEligibilities)).To(Equal(
				map[string]map[string]int{}))
		})

		Context("Computing aggregate statistics for individuals", func() {
			It("Counts total number of individuals in file", func() {
				Expect(dojInformation.TotalIndividuals()).To(Equal(10))
			})

			Context("Before eligibility is run", func() {
				It("Calculates individuals with a felony", func() {
					Expect(dojInformation.CountIndividualsWithFelony()).To(Equal(10))
				})

				It("Calculates individuals with any conviction", func() {
					Expect(dojInformation.CountIndividualsWithConviction()).To(Equal(10))
				})

				It("Calculates individuals with any conviction in the last 7 years", func() {
					Expect(dojInformation.CountIndividualsWithConvictionInLast7Years()).To(Equal(5))
				})
			})

			Context("After eligibility is run", func() {
				It("Calculates individuals who will no longer have a felony ", func() {
					Expect(dojInformation.CountIndividualsNoLongerHaveFelony(dojEligibilities)).To(Equal(1))
				})

				It("Calculates individuals who no longer have any conviction", func() {
					Expect(dojInformation.CountIndividualsNoLongerHaveConviction(dojEligibilities)).To(Equal(1))
				})

				It("Calculates individuals who no longer have any conviction in the last 7 years", func() {
					Expect(dojInformation.CountIndividualsNoLongerHaveConvictionInLast7Years(dojEligibilities)).To(Equal(0))
				})
			})

		})
	})
})
