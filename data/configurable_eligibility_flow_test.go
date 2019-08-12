package data

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("ConfigurableEligibilityFlow", func() {
	const COUNTY = "ARBITRARY"

	var flow EligibilityFlow

	BeforeEach(func() {

		flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
			BaselineEligibility: BaselineEligibility{
				Dismiss: []string{"11357", "11358", "11359"},
				Reduce:  []string{"11360"},
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
				subject                 Subject
				conviction1             DOJRow
				conviction2             DOJRow
				conviction3             DOJRow
				conviction4             DOJRow
				conviction5             DOJRow
				other_county_conviction DOJRow
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
				other_county_conviction = DOJRow{
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

				rows := []DOJRow{conviction1, conviction2, conviction3, conviction4, conviction5, other_county_conviction}
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
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[2].EligibilityReason).To(Equal("Dismiss all HS 11358 convictions"))
				Expect(infos[3].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[3].EligibilityReason).To(Equal("Misdemeanor or Infraction"))
				Expect(infos[4].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[4].EligibilityReason).To(Equal("Reduce all HS 11360 convictions"))
			})
		})

		Context("Matches missing/wrong code section letters and 'attempted' code sections ", func() {
			var (
				subject                             Subject
				attemptedCodeSectionConviction      DOJRow
				missingCodeSectionLettersConviction DOJRow
				wrongCodeSectionLettersConviction   DOJRow
			)

			BeforeEach(func() {
				attemptedCodeSectionConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "66411357(A) HS",
					DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}
				missingCodeSectionLettersConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A)",
					DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}
				wrongCodeSectionLettersConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A)PC",
					DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           2,
					IsFelony:        true,
				}

				rows := []DOJRow{attemptedCodeSectionConviction, missingCodeSectionLettersConviction, wrongCodeSectionLettersConviction}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}
			})

			It("returns a map of eligibility infos", func() {
				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				_, ok := infos[0]
				Expect(ok).To(Equal(true))
				Expect(len(infos)).To(Equal(3))
			})

			It("returns the correct eligibility determination for each conviction", func() {
				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[2].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
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
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{},
						Reduce:  []string{"11357", "11358", "11359", "11360",},
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
				Expect(infos[0].EligibilityReason).To(Equal("Reduce all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("Misdemeanor or Infraction"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all HS 11358 convictions"))
				Expect(infos[3].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[3].EligibilityReason).To(Equal("Reduce all HS 11359 convictions"))
				Expect(infos[4].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[4].EligibilityReason).To(Equal("Reduce all HS 11360 convictions"))
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
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357",},
						Reduce:  []string{"11358", "11359", "11360",},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectUnder21AtConviction: true,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("21 years or younger"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all HS 11360 convictions"))
			})

			It("does not dismiss convictions under the age of 21 if subjectUnder21AtConviction option is not set", func() {
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357",},
						Reduce:  []string{"11358", "11359", "11360",},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectUnder21AtConviction: false,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Equal("Reduce all HS 11359 convictions"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all HS 11360 convictions"))
			})

		})

		Context("When additionalRelief -> yearsSinceConvictionThreshold is set", func() {
			var (
				subject                                    Subject
				conviction1                                DOJRow
				dismissableByYearsSinceConvictionThreshold DOJRow
				yearsSinceConvictionThreshold              int
			)

			BeforeEach(func() {
				yearsSinceConvictionThreshold = 13
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357",},
						Reduce:  []string{"11358", "11359", "11360",},
					},
					AdditionalRelief: AdditionalRelief{
						YearsSinceConvictionThreshold: yearsSinceConvictionThreshold,
					},
				}, COUNTY)
			})

			It("dismisses convictions over the yearsSinceConvictionThreshold setting", func() {
				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: comparisonTime.AddDate(-yearsSinceConvictionThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dismissableByYearsSinceConvictionThreshold = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: comparisonTime.AddDate(-yearsSinceConvictionThreshold-1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, dismissableByYearsSinceConvictionThreshold}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(MatchRegexp("Conviction occurred .* or more years ago"))
			})

			It("does not dismiss convictions under the yearsSinceConvictionThreshold setting", func() {
				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: comparisonTime.AddDate(-yearsSinceConvictionThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dismissableByYearsSinceConvictionThreshold = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: comparisonTime.AddDate(-yearsSinceConvictionThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, dismissableByYearsSinceConvictionThreshold}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Not(MatchRegexp("Conviction occurred .* years ago")))
			})

		})

		Context("When additionalRelief -> yearsSinceConvictionThreshold is not set", func() {
			var (
				subject                                    Subject
				conviction1                                DOJRow
				dismissableByYearsSinceConvictionThreshold DOJRow
				yearsSinceConvictionThreshold              int
			)

			BeforeEach(func() {
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357",},
						Reduce:  []string{"11358", "11359", "11360",},
					},
				}, COUNTY)
			})

			It("does not dismisses convictions for subjects with the reason of being over a certain age", func() {
				yearsSinceConvictionThreshold = 15
				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: comparisonTime.AddDate(-yearsSinceConvictionThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				dismissableByYearsSinceConvictionThreshold = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: comparisonTime.AddDate(-yearsSinceConvictionThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, dismissableByYearsSinceConvictionThreshold}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				for _, info := range infos {
					Expect(info.EligibilityReason).To(Not(MatchRegexp("years or older")))
				}
			})

		})

		Context("When additionalRelief -> yearsCrimeFreeThreshold is set", func() {
			var (
				subject                          Subject
				conviction1                      DOJRow
				potentiallyDismissableConviction DOJRow
				mostRecentConviction             DOJRow
				mostRecentConvictionDate         time.Time
				yearsCrimeFreeThreshold          int
			)

			BeforeEach(func() {
				yearsCrimeFreeThreshold = 6
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357",},
						Reduce:  []string{"11358", "11359", "11360",},
					},
					AdditionalRelief: AdditionalRelief{
						YearsCrimeFreeThreshold: yearsCrimeFreeThreshold,
					},
				}, COUNTY)
			})

			It("dismisses Prop 64 convictions if all convictions are older than the yearsCrimeFreeThreshold setting", func() {
				mostRecentConvictionDate = comparisonTime.AddDate(-(yearsCrimeFreeThreshold + 1), 0, 0)

				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: comparisonTime.AddDate(-10, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				potentiallyDismissableConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: comparisonTime.AddDate(-10, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				mostRecentConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "1234 HS",
					DispositionDate: mostRecentConvictionDate,
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, potentiallyDismissableConviction, mostRecentConviction}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(MatchRegexp("No convictions in the past .* years"))
			})

			It("does not dismiss additional Prop 64 convictions if any conviction is newer than the yearsCrimeFreeThreshold setting", func() {
				mostRecentConvictionDate = comparisonTime.AddDate(-(yearsCrimeFreeThreshold - 1), 0, 0)

				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: comparisonTime.AddDate(-yearsCrimeFreeThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				potentiallyDismissableConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: comparisonTime.AddDate(-yearsCrimeFreeThreshold+1, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				mostRecentConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "1234 HS",
					DispositionDate: mostRecentConvictionDate,
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, potentiallyDismissableConviction, mostRecentConviction}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Not(MatchRegexp("No convictions in the past .* years")))
			})

		})

		Context("When additionalRelief -> yearsCrimeFreeThreshold is not set", func() {
			var (
				subject                          Subject
				conviction1                      DOJRow
				potentiallyDismissableConviction DOJRow
				mostRecentConviction             DOJRow
				mostRecentConvictionDate         time.Time
				yearsCrimeFreeThreshold          int
			)

			BeforeEach(func() {
				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357",},
						Reduce:  []string{"11358", "11359", "11360",},
					},
				}, COUNTY)
			})

			It("does not dismiss or reduce convictions for subjects with the reason of age of last conviction", func() {
				yearsCrimeFreeThreshold = 2

				mostRecentConvictionDate = comparisonTime.AddDate(-(yearsCrimeFreeThreshold + 1), 0, 0)

				conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11357(A) HS",
					DispositionDate: comparisonTime.AddDate(-10, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				potentiallyDismissableConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: comparisonTime.AddDate(-10, 0, 0),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				mostRecentConviction = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "1234 HS",
					DispositionDate: mostRecentConvictionDate,
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}

				rows := []DOJRow{conviction1, potentiallyDismissableConviction, mostRecentConviction}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				for _, info := range infos {
					Expect(info.EligibilityReason).To(Not(MatchRegexp("in the past .* years")))
				}
			})

		})

		Context("When additionalRelief -> subjectHasOnlyProp64Charges", func() {
			var (
				subject              Subject
				prop64Conviction1    DOJRow
				prop64Conviction2    DOJRow
				prop64Conviction3    DOJRow
				nonProp64Conviction1 DOJRow
			)

			BeforeEach(func() {
				prop64Conviction1 = DOJRow{
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
				prop64Conviction2 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}
				prop64Conviction3 = DOJRow{
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
			})

			It("dismisses all convictions if they are all prop 64 convictions", func() {
				rows := []DOJRow{prop64Conviction1, prop64Conviction2, prop64Conviction3}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357"},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectHasOnlyProp64Charges: true,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("Only has 11357-60 charges"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[2].EligibilityReason).To(Equal("Only has 11357-60 charges"))
			})

			It("does not dismiss convictions if there are any non prop 64 convictions", func() {

				nonProp64Conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "5555 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				rows := []DOJRow{prop64Conviction1, prop64Conviction2, prop64Conviction3, nonProp64Conviction1}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357"},
						Reduce:  []string{"11358", "11359", "11360",},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectHasOnlyProp64Charges: true,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[1].EligibilityReason).To(Equal("Reduce all HS 11359 convictions"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[2].EligibilityReason).To(Equal("Reduce all HS 11360 convictions"))
			})

		})

		Context("When additionalRelief -> subjectIsDeceased", func() {
			var (
				subject           Subject
				prop64Conviction1 DOJRow
				prop64Conviction2 DOJRow
				prop64Conviction3 DOJRow
				randomConviction1 DOJRow
			)

			BeforeEach(func() {
				prop64Conviction1 = DOJRow{
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
				prop64Conviction2 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11359 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           1,
					IsFelony:        true,
				}
				prop64Conviction3 = DOJRow{
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
				randomConviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "57675 HS",
					DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1119999",
					County:          COUNTY,
					CountOrder:      "102001003000",
					Index:           2,
					IsFelony:        true,
				}
			})

			It("dismisses all p64 convictions if individual is deceased", func() {
				rows := []DOJRow{prop64Conviction1, prop64Conviction2, prop64Conviction3, randomConviction1}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				subject.IsDeceased = true

				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357"},
					},
					AdditionalRelief: AdditionalRelief{
						SubjectIsDeceased: true,
					},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[0].EligibilityReason).To(Equal("Dismiss all HS 11357 convictions"))
				Expect(infos[1].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[1].EligibilityReason).To(Equal("Individual is deceased"))
				Expect(infos[2].EligibilityDetermination).To(Equal("Eligible for Dismissal"))
				Expect(infos[2].EligibilityReason).To(Equal("Individual is deceased"))
			})

			It("does not dismiss convictions for a deceased person if additional relief is not selected", func() {

				prop64Conviction1 = DOJRow{
					DOB:             birthDate,
					WasConvicted:    true,
					CodeSection:     "11360 HS",
					DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC),
					OFN:             "1234",
					County:          COUNTY,
					CountOrder:      "101001001000",
					Index:           0,
					IsFelony:        true,
				}

				rows := []DOJRow{prop64Conviction1}
				subject = Subject{}
				for _, row := range rows {
					subject.PushRow(row, flow)
				}

				subject.IsDeceased = true

				flow, _ = NewConfigurableEligibilityFlow(EligibilityOptions{
					BaselineEligibility: BaselineEligibility{
						Dismiss: []string{"11357"},
						Reduce:  []string{"11358", "11359", "11360",},
					},
					AdditionalRelief: AdditionalRelief{},
				}, COUNTY)

				infos := flow.ProcessSubject(&subject, comparisonTime, COUNTY)
				Expect(infos[0].EligibilityDetermination).To(Equal("Eligible for Reduction"))
				Expect(infos[0].EligibilityReason).To(Equal("Reduce all HS 11360 convictions"))
			})

		})

	})
})
