package data_test

import (
	"gogen/data"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DOJHistory", func() {
	var history data.DOJHistory
	var conviction1 data.DOJRow
	var conviction2 data.DOJRow
	var conviction3 data.DOJRow
	var conviction4 data.DOJRow
	var conviction5 data.DOJRow
	var nonConviction data.DOJRow
	var birthDate time.Time

	BeforeEach(func() {
		birthDate = time.Date(1994, time.April, 10, 0, 0, 0, 0, time.UTC)
		conviction1 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1234", DOB: birthDate, CodeSection: "11357 HS", Convicted: true, CycleDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", NumCrtCase: "777CRTCASE"}
		nonConviction = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1235", DOB: birthDate, CodeSection: "11357 HS", Convicted: false, CycleDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC)}
		conviction2 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
		conviction3 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "286(Q)(1) PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
		conviction4 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
		conviction5 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
		registration := data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "290 PC", Convicted: false, CycleDate: time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC), DispositionDate: time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC), PC290Registration: true}

		rows := []data.DOJRow{conviction1, nonConviction, conviction2, registration, conviction3, conviction4, conviction5}
		history = data.DOJHistory{}
		for _, row := range rows {
			history.PushRow(row)
		}
	})

	Describe("PushRow", func() {
		It("Sets the correct values on the history", func() {
			Expect(history.SubjectID).To(Equal("subj_id"))
			Expect(history.Name).To(Equal("SOUP,ZAK E"))
			Expect(history.WeakName).To(Equal("SOUP,ZAK"))
			Expect(history.CII).To(Equal("12345678"))
			Expect(history.SSN).To(Equal("345678125"))
			Expect(history.DOB).To(Equal(birthDate))
			Expect(history.CDL).To(Equal("testcdl"))
			Expect(history.PC290Registration).To(BeTrue())
			Expect(history.Convictions).To(ConsistOf(&conviction1, &conviction2, &conviction3, &conviction4, &conviction5))
		})
	})

	Describe("Match", func() {
		var cmsEntry data.CMSEntry

		BeforeEach(func() {
			cmsEntry = data.CMSEntry{
				CourtNumber:     "no match",
				Level:           "foo",
				SSN:             "bar",
				CII:             "bla",
				Charge:          "bla",
				IncidentNumber:  "bla",
				Name:            "BLAH/BLAH",
				FormattedName:   "BLAH,BLAH",
				WeakName:        "BLAH,BLAH",
				CDL:             "bla",
				DateOfBirth:     time.Date(1985, time.May, 12, 0, 0, 0, 0, time.UTC),
				DispositionDate: time.Time{},
				RawRow:          nil,
			}
		})

		It("Does not match if none of the provided identifiers match", func() {
			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": false,
					"courtno":        false,
					"cdl":            false,
				},
				MatchStrength: 0,
			}))
		})

		It("matches name and dob", func() {
			cmsEntry.Name = "SOUP/ZAK/E"
			cmsEntry.FormattedName = "SOUP,ZAK E"
			cmsEntry.WeakName = "SOUP,ZAK"
			cmsEntry.DateOfBirth = birthDate

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     true,
					"weakNameAndDob": true,
					"courtno":        false,
					"cdl":            false,
				},
				MatchStrength: 2,
			}))
		})

		It("matches weak name and dob (excluding middle name)", func() {
			cmsEntry.Name = "SOUP/ZAK/F"
			cmsEntry.FormattedName = "SOUP,ZAK F"
			cmsEntry.WeakName = "SOUP,ZAK"
			cmsEntry.DateOfBirth = birthDate

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": true,
					"courtno":        false,
					"cdl":            false,
				},
				MatchStrength: 1,
			}))
		})

		It("doesn't consider zero values to be matching", func() {
			history.DOB = time.Time{}

			cmsEntry.Name = "SOUP/ZAK/E"
			cmsEntry.FormattedName = "SOUP,ZAK E"
			cmsEntry.DateOfBirth = time.Time{}

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":     false,
					"ssn":     false,
					"courtno": false,
					"cdl":     false,
				},
				MatchStrength: 0,
			}))
		})

		It("matches on CII", func() {
			cmsEntry.CII = "12345678"

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            true,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": false,
					"courtno":        false,
					"cdl":            false,
				},
				MatchStrength: 1,
			}))
		})

		It("matches if the courtno matches the OFN for any conviction", func() {
			cmsEntry.CourtNumber = "1236 334455-00"

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": false,
					"courtno":        true,
					"cdl":            false,
				},
				MatchStrength: 1,
			}))
		})

		It("matches if the courtno matches partial OFN for any conviction", func() {
			cmsEntry.CourtNumber = "334455"

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": false,
					"courtno":        true,
					"cdl":            false,
				},
				MatchStrength: 1,
			}))
		})

		It("matches if the courtno matches FE_NUM_CRT_CASE for any conviction", func() {
			cmsEntry.CourtNumber = "777CRTCASE"

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": false,
					"courtno":        true,
					"cdl":            false,
				},
				MatchStrength: 1,
			}))
		})

		It("only matches if the for courtno if the conviction is in the right county", func() {
			cmsEntry.CourtNumber = "1119999"

			Expect(history.Match(cmsEntry)).To(Equal(data.MatchData{
				History: &history,
				MatchResults: map[string]bool{
					"cii":            false,
					"ssn":            false,
					"nameAndDob":     false,
					"weakNameAndDob": false,
					"courtno":        false,
					"cdl":            false,
				},
				MatchStrength: 0,
			}))
		})
	})

	Describe("pc290CodeSections", func() {
		It("Returns a list of convicted code sections that match PC290 registerable offenses", func() {
			Expect(history.PC290CodeSections()).To(ConsistOf("286(Q)(1) PC"))
		})
	})

	Describe("ThreeConvictionsSameCode", func() {
		It("returns true if there are three convictions of same code section", func() {
			Expect(history.ThreeConvictionsSameCode("11360HS")).To(BeTrue())
			Expect(history.ThreeConvictionsSameCode("11357HS")).To(BeFalse())
		})

		It("does not consider convictions to be separate if they have the same cycle date", func() {
			conviction5.CycleDate = conviction4.CycleDate
			history = data.DOJHistory{}
			history.PushRow(conviction2)
			history.PushRow(conviction4)
			history.PushRow(conviction5)

			Expect(history.ThreeConvictionsSameCode("11360HS")).To(BeFalse())
		})
	})

	Describe("MostRecentConvictionDate", func() {
		It("returns the most recent conviction date", func() {
			Expect(history.MostRecentConvictionDate()).To(Equal(time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC)))
		})
	})

})
