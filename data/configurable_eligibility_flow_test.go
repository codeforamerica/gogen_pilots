package data

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("configurableEligibilityFlow", func() {
	const COUNTY = "ARBITRARY"

	var flow EligibilityFlow

	BeforeEach(func() {
		flow = NewConfigurableEligibilityFlow(EligibilityOptions{
			BaselineEligibility: BaselineEligibility{
				Dismiss: []string{"11357(A)", "11357(B)", "11357(C)", "11358", "11359"},
			},
		}, COUNTY)

	})

	Describe("Processing a subject", func() {
		birthDate := time.Date(1978, time.April, 10, 0, 0, 0, 0, time.UTC)
		comparisonTime := time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC)

		Context("Filtering relevant convictions", func() {
			var (
				subject               Subject
				conviction1           DOJRow
				conviction2           DOJRow
				nonProp64conviction   DOJRow
				otherCountyConviction DOJRow
				registration          DOJRow
				nonConviction         DOJRow
			)

			BeforeEach(func() {
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
				registration = DOJRow{
					DOB:                 birthDate,
					WasConvicted:        false,
					CodeSection:         "290 PC",
					DispositionDate:     time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC),
					OFN:                 "1236 12345678-00",
					IsPC290Registration: true,
					County:              COUNTY,
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
					CodeSection:     "11358 PC",
					DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1119999",
					County:          COUNTY,
					CountOrder:      "102001003000",
					Index:           2,
				}
				nonProp64conviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "187 PC",
					DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1118888",
					County:          COUNTY,
					CountOrder:      "103001004000",
					Index:           3,
				}
				otherCountyConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11360 HS",
					DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC),
					OFN:             "1236 12345678-00",
					County:          "OTHER COUNTY",
					CountOrder:      "104001005000",
					Index:           4,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, conviction2, nonProp64conviction, otherCountyConviction, registration, nonConviction}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}
			})

			It("only returns eligibility infos for Prop 64 convictions in the given county", func() {
				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
				_, ok = infos[2]
				Expect(ok).To(Equal(true))
				Expect(len(infos)).To(Equal(2))
			})
		})

		Context("Dismissing and reducing by code section", func() {
			var (
				subject     Subject
				conviction1 DOJRow
				conviction2 DOJRow
				conviction3 DOJRow
				conviction4 DOJRow
				conviction5 DOJRow
				conviction6 DOJRow
			)

			BeforeEach(func() {
				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
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
					CodeSection:     "11357(D) HS",
					DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}
				conviction3 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11358 HS",
					DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1119999",
					County:          COUNTY,
					CountOrder:      "102001003000",
					Index:           2,
					IsFelony:        true,
				}
				conviction4 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1118888",
					County:          COUNTY,
					CountOrder:      "103001004000",
					Index:           3,
				}
				conviction5 = DOJRow{
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
				conviction6 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11360 PC",
					DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC),
					OFN:             "1236 334455-00",
					County:          "OTHER COUNTY",
					CountOrder:      "104001006000",
					Index:           5,
					SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
				}

				rows := []DOJRow{conviction1, conviction2, conviction3, conviction4, conviction5, conviction6}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}
			})

			It("returns a map of eligibility infos", func() {
				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
				Expect(len(infos)).To(Equal(5))
			})

			It("returns the correct eligibility determination for each conviction", func() {
				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Misdemeanor or Infraction"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Equal("Reduce all 11357(D) HS convictions"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[2].EligibilityReason).To(Equal("Dismiss all 11358 HS convictions"))
				Expect(infos[3].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[3].EligibilityReason).To(Equal("Misdemeanor or Infraction"))
				Expect(infos[4].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[4].EligibilityReason).To(Equal("Reduce all 11360 HS convictions"))
			})
		})

		Context("When a matcher is empty", func() {
			var (
				subject     Subject
				conviction1 DOJRow
				conviction2 DOJRow
				conviction3 DOJRow
				conviction4 DOJRow
				conviction5 DOJRow
				conviction6 DOJRow
			)

			BeforeEach(func() {
				flow = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{},
					},
				}, COUNTY)

				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}
				conviction2 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(D) HS",
					DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        false,
				}
				conviction3 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11358 HS",
					DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1119999",
					County:          COUNTY,
					CountOrder:      "102001003000",
					Index:           2,
					IsFelony:        true,
				}
				conviction4 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1118888",
					County:          COUNTY,
					CountOrder:      "103001004000",
					Index:           3,
					IsFelony:        true,
				}
				conviction5 = DOJRow{
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
				conviction6 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11360 PC",
					DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC),
					OFN:             "1236 334455-00",
					County:          "OTHER COUNTY",
					CountOrder:      "104001006000",
					Index:           5,
					IsFelony:        true,
					SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC),
				}

				rows := []DOJRow{conviction1, conviction2, conviction3, conviction4, conviction5, conviction6}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}
			})

			It("returns the correct eligibility determination for each conviction", func() {
				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[0].EligibilityReason).To(Equal("Reduce all 11357(A) HS convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("Misdemeanor or Infraction"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all 11358 HS convictions"))
				Expect(infos[3].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[3].EligibilityReason).To(Equal("Reduce all 11359 HS convictions"))
				Expect(infos[4].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[4].EligibilityReason).To(Equal("Reduce all 11360 HS convictions"))
			})
		})

		Context("When additionalRelief -> subjectUnder21AtConviction", func() {
			var (
				subject           Subject
				conviction1       DOJRow
				under21Conviction DOJRow
				conviction3       DOJRow
			)

			BeforeEach(func() {
				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dateWhenSubjectWas16 := birthDate.AddDate(16, 0, 0)
				under21Conviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: dateWhenSubjectWas16,
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}
				conviction3 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11360 HS",
					DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1119999",
					County:          COUNTY,
					CountOrder:      "102001003000",
					Index:           2,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, under21Conviction, conviction3}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}
			})

			It("dismisses convictions under the age of 21 if subjectUnder21AtConviction option is set", func() {
				flow = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357(A)", "11357(B)", "11357(C)", "11357(D)",},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectUnder21AtConviction: true,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all 11357(A) HS convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("21 years or younger"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all 11360 HS convictions"))
			})

			It("does not dismiss convictions under the age of 21 if subjectUnder21AtConviction option is not set", func() {
				flow = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357(A)", "11357(B)", "11357(C)", "11357(D)",},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectUnder21AtConviction: false,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all 11357(A) HS convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Equal("Reduce all 11359 HS convictions"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all 11360 HS convictions"))
			})

		})

		Context("When additionalRelief -> dismissByAge is set", func() {
			var (
				conviction1                DOJRow
				dismissableByAgeConviction DOJRow
				DismissByAgeSetting        int
			)

			BeforeEach(func() {
				DismissByAgeSetting = 45
				flow = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357(A)", "11357(B)", "11357(C)", "11357(D)",},
					},
					AdditionalRelief: AdditionalRelief{
						DismissByAge: DismissByAgeSetting,
					},
				}, COUNTY)
			})

			It("dismisses convictions for subjects over the dismissbyAge setting", func() {
				olderSubject := Subject{DOB: comparisonTime.AddDate(-DismissByAgeSetting-1, 0, 0)}
				conviction1 = DOJRow{
					DOB:             olderSubject.DOB,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dismissableByAgeConviction = DOJRow{
					DOB:             olderSubject.DOB,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, dismissableByAgeConviction}
				for _, row := range rows {
					olderSubject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&olderSubject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all 11357(A) HS convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("45 years or older"))
			})

			It("does not dismiss convictions for subjects under the dismissbyAge setting", func() {
				youngerSubject := Subject{DOB: comparisonTime.AddDate(-DismissByAgeSetting+1, 0, 0)}
				conviction1 = DOJRow{
					DOB:             youngerSubject.DOB,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dismissableByAgeConviction = DOJRow{
					DOB:             youngerSubject.DOB,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, dismissableByAgeConviction}
				for _, row := range rows {
					youngerSubject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&youngerSubject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all 11357(A) HS convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Equal("Reduce all 11359 HS convictions"))
			})

		})

		Context("When additionalRelief -> dismissByAge is not set", func() {
			var (
				conviction1                DOJRow
				dismissableByAgeConviction DOJRow
				DismissByAgeSetting        int
			)

			BeforeEach(func() {
				DismissByAgeSetting = 45
				flow = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357(A)", "11357(B)", "11357(C)", "11357(D)",},
					},
				}, COUNTY)
			})

			It("does not dismisses convictions for subjects with the reason of being over a certain age", func() {
				olderSubject := Subject{DOB: comparisonTime.AddDate(-DismissByAgeSetting-1, 0, 0)}
				conviction1 = DOJRow{
					DOB:             olderSubject.DOB,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dismissableByAgeConviction = DOJRow{
					DOB:             olderSubject.DOB,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, dismissableByAgeConviction}
				for _, row := range rows {
					olderSubject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&olderSubject, comparisonTime, COUNTY)
				for _, info := range infos {
					Expect(info.EligibilityReason).To(Not(MatchRegexp("years or older")))
				}
			})

		})
	})
})
