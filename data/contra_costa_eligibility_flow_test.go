package data_test

import (
	. "gogen/data"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("contraCostaEligibilityFlow", func() {
	const COUNTY = "CONTRA COSTA"

	var flow EligibilityFlow

	BeforeEach(func() {
		flow = EligibilityFlows[COUNTY]
	})

	Describe("MatchedCodeSection", func() {
		It("returns the matched substring for a given Prop 64 code section", func() {
			Expect(flow.MatchedCodeSection("11358(c) HS")).To(Equal("11358"))
			Expect(flow.MatchedCodeSection("/11357 HS")).To(Equal("11357"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedCodeSection("12345(c) HS")).To(Equal(""))
			Expect(flow.MatchedCodeSection("647(f) HS")).To(Equal(""))
			Expect(flow.MatchedCodeSection("4050.6 BP")).To(Equal(""))
			Expect(flow.MatchedCodeSection("14859 PC")).To(Equal(""))
		})

		It("returns empty string if the code section is for a related charge", func() {
			Expect(flow.MatchedCodeSection("647(f) PC")).To(Equal(""))
			Expect(flow.MatchedCodeSection("148.9 PC")).To(Equal(""))
			Expect(flow.MatchedCodeSection("4060    BP")).To(Equal(""))
			Expect(flow.MatchedCodeSection("--40508 VC--")).To(Equal(""))
			Expect(flow.MatchedCodeSection("1320(a) PC")).To(Equal(""))
		})
	})

	Describe("MatchedRelatedCodeSection", func() {
		It("returns the matched substring for a given related charge code section", func() {
			Expect(flow.MatchedRelatedCodeSection("647(f) PC")).To(Equal("647(f) PC"))
			Expect(flow.MatchedRelatedCodeSection("148.9 PC")).To(Equal("148.9 PC"))
			Expect(flow.MatchedRelatedCodeSection("4060    BP")).To(Equal("4060    BP"))
			Expect(flow.MatchedRelatedCodeSection("--40508 VC--")).To(Equal("40508 VC"))
			Expect(flow.MatchedRelatedCodeSection("1320(a) PC")).To(Equal("1320(a) PC"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedRelatedCodeSection("12345(c) HS")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("647(f) HS")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("4050.6 BP")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("14859 PC")).To(Equal(""))
		})

		It("returns empty string if the code section is for a Prop 64 charge", func() {
			Expect(flow.MatchedRelatedCodeSection("11358(c) HS")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("/11357 HS")).To(Equal(""))
		})
	})

	Describe("Processing a history", func() {
		var (
			history           DOJHistory
			conviction1       DOJRow
			conviction2       DOJRow
			conviction3       DOJRow
			conviction4       DOJRow
			conviction5       DOJRow
			conviction5Prison DOJRow
			nonConviction     DOJRow
			birthDate         time.Time
			comparisonTime    time.Time
		)

		BeforeEach(func() {
			days := time.Duration(24) * (time.Hour)
			birthDate = time.Date(1994, time.April, 10, 0, 0, 0, 0, time.UTC)
			conviction1 = DOJRow{
				DOB:             birthDate,
				WasConvicted:       true,
				CodeSection:     "11357 HS",
				DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
				OFN:             "1234",
				County:          COUNTY,
				CountOrder:      "101001001000",
				Index:           0,
				IsFelony:          false,
			}
			nonConviction = DOJRow{
				DOB:             birthDate,
				WasConvicted:       false,
				CodeSection:     "11357 HS",
				DispositionDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC),
				OFN:             "1235",
				County:          COUNTY,
				CountOrder:      "101001002000",
				Index:           1,
			}
			conviction2 = DOJRow{
				DOB:             birthDate,
				WasConvicted:       true,
				CodeSection:     "602 PC",
				DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
				OFN:             "1119999",
				County:          COUNTY,
				CountOrder:      "102001003000",
				Index:           2,
			}
			conviction3 = DOJRow{
				DOB:             birthDate,
				WasConvicted:       true,
				CodeSection:     "187 PC",
				DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC),
				OFN:             "1118888",
				County:          "LOS ANGELES",
				CountOrder:      "103001004000",
				Index:           3,
			}
			conviction4 = DOJRow{
				DOB:             birthDate,
				WasConvicted:       true,
				CodeSection:     "11360 HS",
				DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 12345678-00",
				County:          COUNTY,
				CountOrder:      "104001005000",
				Index:           4,
				IsFelony:          true,
			}
			conviction5 = DOJRow{
				DOB:             birthDate,
				WasConvicted:       true,
				CodeSection:     "266J PC",
				DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001006000",
				Index:           5,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction5Prison = DOJRow{
				DOB:                  birthDate,
				WasConvicted:            true,
				CodeSection:          "11360 HS",
				DispositionDate:      time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:                  "1236 334455-00",
				County:               COUNTY,
				CountOrder:           "104001006000",
				Index:                6,
				SentencePartDuration: time.Duration(30 * days),
				IsFelony:               true,
			}
			registration := DOJRow{
				DOB:               birthDate,
				WasConvicted:         false,
				CodeSection:       "290 PC",
				DispositionDate:   time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC),
				OFN:               "1236 12345678-00",
				IsPC290Registration: true,
				County:            "",
				CountOrder:        "105001007000",
				Index:             7,
			}

			comparisonTime = time.Date(2019, 4, 10, 0, 0, 0, 0, time.UTC)

			rows := []DOJRow{conviction1, nonConviction, conviction2, registration, conviction3, conviction4, conviction5, conviction5Prison}
			history = DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}
		})

		It("returns a map of eligibility infos", func() {
			infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime)
			Expect(len(infos)).To(Equal(4))
			_, ok := infos[0]
			Expect(ok).To(Equal(true))
			_, ok = infos[2]
			Expect(ok).To(Equal(true))
			_, ok = infos[4]
			Expect(ok).To(Equal(true))
			_, ok = infos[6]
			Expect(ok).To(Equal(true))
		})

		It("returns the correct eligibility determination for each conviction", func() {
			infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime)
			Expect(len(infos)).To(Equal(4))
			Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
			Expect(infos[0].EligibilityReason).To(Equal("Misdemeanor or Infraction"))
			Expect(infos[2].EligibilityDetermination).To(Equal("Maybe Eligible - Flag for Review"))
			Expect(infos[2].EligibilityReason).To(Equal("No Related Prop64 Charges"))
			Expect(infos[4].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
			Expect(infos[4].EligibilityReason).To(Equal("No convictions in past 5 years"))
			Expect(infos[6].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
			Expect(infos[6].EligibilityReason).To(Equal("No convictions in past 5 years"))
		})
	})
})
