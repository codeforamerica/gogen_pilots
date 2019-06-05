package data

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("losAngelesEligibilityFlow", func() {
	const COUNTY = "LOS ANGELES"

	var flow EligibilityFlow

	BeforeEach(func() {
		flow = EligibilityFlows[COUNTY]
	})

	Describe("MatchedCodeSection", func() {
		It("returns the matched substring for a given code section", func() {
			Expect(flow.MatchedCodeSection("11358(c) HS")).To(Equal("11358"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedCodeSection("12345(c) HS")).To(Equal(""))
		})

		It("recognizes attempted code sections for Prop 64", func() {
			Expect(flow.MatchedCodeSection("664.11357(c) HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("66411357(c) HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("664-11357(c) HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("664/11357(c) HS")).To(Equal("11357"))
		})
	})

	Describe("Processing a history", func() {

		birthDate := time.Date(1994, time.April, 10, 0, 0, 0, 0, time.UTC)
		days := time.Duration(24) * (time.Hour)
		comparisonTime := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)
		Context("Steps 1-4", func() {
			var (
				history          DOJHistory
				conviction1      DOJRow
				conviction2      DOJRow
				conviction3      DOJRow
				conviction4      DOJRow
				conviction5      DOJRow
				convictionPrison DOJRow
				nonConviction    DOJRow
			)

			nonConviction = DOJRow{
				DOB:             birthDate,
				WasConvicted:    false,
				CodeSection:     "11357 HS",
				DispositionDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC),
				OFN:             "1235",
				County:          COUNTY,
				CountOrder:      "101001002000",
				Index:           1,
			}
			registration := DOJRow{
				DOB:                 birthDate,
				WasConvicted:        false,
				CodeSection:         "290 PC",
				DispositionDate:     time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC),
				OFN:                 "1236 12345678-00",
				IsPC290Registration: true,
				County:              "",
				CountOrder:          "105001007000",
				Index:               7,
			}
			conviction1 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11357 HS",
				DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
				OFN:             "1234",
				County:          COUNTY,
				CountOrder:      "101001001000",
				Index:           0,
				IsFelony:        false,
			}
			conviction2 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "602 PC",
				DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
				OFN:             "1119999",
				County:          COUNTY,
				CountOrder:      "102001003000",
				Index:           2,
			}
			conviction3 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "187 PC",
				DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC),
				OFN:             "1118888",
				County:          "LOS ANGELES",
				CountOrder:      "103001004000",
				Index:           3,
			}
			conviction4 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11360 HS",
				DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 12345678-00",
				County:          COUNTY,
				CountOrder:      "104001005000",
				Index:           4,
				IsFelony:        true,
			}
			conviction5 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "266J PC",
				DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001006000",
				Index:           5,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			convictionPrison = DOJRow{
				DOB:                  birthDate,
				WasConvicted:         true,
				CodeSection:          "11360 HS",
				DispositionDate:      time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:                  "1236 334455-00",
				County:               COUNTY,
				CountOrder:           "104001006000",
				Index:                6,
				SentencePartDuration: time.Duration(30 * days),
				IsFelony:             true,
			}

			rows := []DOJRow{nonConviction, registration, conviction1, conviction2, conviction3, conviction4, conviction5, convictionPrison}
			history = DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
				_, ok = infos[4]
				Expect(ok).To(Equal(true))
				_, ok = infos[6]
				Expect(ok).To(Equal(true))
				Expect(len(infos)).To(Equal(3))
			})

			It("returns the correct eligibility determination for each conviction", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				//Expect(len(infos)).To(Equal(3))
				Expect(infos[0].EligibilityDetermination).To(Equal("Maybe Eligible - Flag for Review"))
				Expect(infos[0].EligibilityReason).To(Equal("Other 11357"))
				Expect(infos[4].EligibilityDetermination).To(Equal("Not eligible"))
				Expect(infos[4].EligibilityReason).To(Equal("PC 667(e)(2)(c)(iv)"))
				Expect(infos[6].EligibilityDetermination).To(Equal("Not eligible"))
				Expect(infos[6].EligibilityReason).To(Equal("PC 667(e)(2)(c)(iv)"))
			})
		})

		Context("Steps Two Priors and Age checks", func() {
			var (
				conviction6  DOJRow
				conviction7  DOJRow
				conviction8  DOJRow
				conviction9  DOJRow
				conviction10 DOJRow
				conviction11 DOJRow
			)
			conviction6 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(b) HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction7 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(a) HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction8 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(c) HS",
				DispositionDate: time.Date(2006, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}
			conviction9 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11355 HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction10 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11355 HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction11 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11355 HS",
				DispositionDate: time.Date(2006, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}

			rows := []DOJRow{conviction6, conviction7, conviction8, conviction9, conviction10, conviction11}
			history := DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
				_, ok = infos[2]
				Expect(ok).To(Equal(true))
				Expect(len(infos)).To(Equal(2))
			})

			It("Returns the correct eligbility determination", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(2))
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("No convictions in past 10 years"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Not eligible"))
				Expect(infos[2].EligibilityReason).To(Equal("Two priors"))
			})
		})

		Context("All Prop 64 convictions and completed sentences", func() {

			var (
				conviction6 DOJRow
				conviction7 DOJRow
				conviction8 DOJRow
			)
			conviction6 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(b) HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction7 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11358 HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction8 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11357 HS",
				DispositionDate: time.Date(2006, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}

			rows := []DOJRow{conviction6, conviction7, conviction8}
			history := DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(2))
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
				_, ok = infos[2]
				Expect(ok).To(Equal(true))
			})

			It("Returns the correct eligbility determination", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(2))
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Only has 11357-60 charges and completed sentence"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Maybe Eligible - Flag for Review"))
				Expect(infos[2].EligibilityReason).To(Equal("Other 11357"))
			})
		})

		Context("No convictions past 10 years", func() {

			var (
				conviction6 DOJRow
				conviction7 DOJRow
				conviction8 DOJRow
			)
			conviction6 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(b) HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction7 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11358 HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction8 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "666 HS",
				DispositionDate: time.Date(2006, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}

			rows := []DOJRow{conviction6, conviction7, conviction8}
			history := DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
			})

			It("Returns the correct eligbility determination", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("No convictions in past 10 years"))
			})
		})

		Context("Still serving a sentence", func() {

			var (
				conviction6 DOJRow
				conviction7 DOJRow
				conviction8 DOJRow
			)
			conviction6 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(b) HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction7 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11358 HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction8 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "666 HS",
				DispositionDate: time.Date(2016, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2025, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}

			rows := []DOJRow{conviction6, conviction7, conviction8}
			history := DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
			})

			It("Returns the correct eligbility determination", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				Expect(infos[0].EligibilityDetermination).To(Equal("Hand Review"))
				Expect(infos[0].EligibilityReason).To(Equal("Currently serving sentence"))
			})
		})

		Context("Still serving a sentence", func() {

			var (
				conviction6 DOJRow
				conviction7 DOJRow
				conviction8 DOJRow
			)
			conviction6 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(b) HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction7 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11358 HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction8 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "666 HS",
				DispositionDate: time.Date(2016, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2005, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}

			rows := []DOJRow{conviction6, conviction7, conviction8}
			history := DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}
			history.IsDeceased = true

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
			})

			It("Returns the correct eligbility determination", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Deceased"))
			})
		})

		Context("None of the above", func() {

			var (
				conviction6 DOJRow
				conviction7 DOJRow
				conviction8 DOJRow
			)
			conviction6 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11359(b) HS",
				DispositionDate: time.Date(2004, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001007000",
				Index:           0,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction7 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "11358 HS",
				DispositionDate: time.Date(2005, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          "SAN JOAQUIN",
				CountOrder:      "104001008000",
				Index:           1,
				SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
			}
			conviction8 = DOJRow{
				DOB:             birthDate,
				WasConvicted:    true,
				CodeSection:     "666 HS",
				DispositionDate: time.Date(2016, time.December, 5, 0, 0, 0, 0, time.UTC),
				OFN:             "1236 334455-00",
				County:          COUNTY,
				CountOrder:      "104001009000",
				Index:           2,
				SentenceEndDate: time.Date(2005, 03, 04, 0, 0, 0, 0, time.UTC),
				IsFelony:        true,
			}

			rows := []DOJRow{conviction6, conviction7, conviction8}
			history := DOJHistory{}
			for _, row := range rows {
				history.PushRow(row, COUNTY)
			}

			It("returns a map of eligibility infos", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
			})

			It("Returns the correct eligbility determination", func() {
				infos := EligibilityFlows[COUNTY].ProcessHistory(&history, comparisonTime, COUNTY)
				Expect(len(infos)).To(Equal(1))
				Expect(infos[0].EligibilityDetermination).To(Equal("Hand Review"))
				Expect(infos[0].EligibilityReason).To(Equal("????"))
			})
		})

	})
})
