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
	var conviction6 data.DOJRow
	var conviction7 data.DOJRow
	var conviction8 data.DOJRow
	var conviction9 data.DOJRow
	var conviction10 data.DOJRow
	var conviction11 data.DOJRow
	var conviction12 data.DOJRow
	var conviction13 data.DOJRow
	var conviction14 data.DOJRow
	var conviction15 data.DOJRow
	var conviction5Prison data.DOJRow
	var nonConviction data.DOJRow
	var birthDate time.Time

	days := time.Duration(24) * (time.Hour)
	BeforeEach(func() {
		birthDate = time.Date(1994, time.April, 10, 0, 0, 0, 0, time.UTC)
		conviction1 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1234", DOB: birthDate, CodeSection: "11357 HS", Convicted: true, CycleDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "101001001000", DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", NumCrtCase: "777CRTCASE"}
		nonConviction = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1235", DOB: birthDate, CodeSection: "11357 HS", Convicted: false, CycleDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC), CountOrder: "101001002000", DispositionDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC)}
		conviction2 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "102001003000", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
		conviction3 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "286(Q)(1) PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "103001004000", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
		conviction4 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), CountOrder: "104001005000", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
		conviction5 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006000", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
		conviction5Prison = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "11360 HS", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006000", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentencePartDuration: time.Duration(30 * days)}
		registration := data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A05555555", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "290 PC", Convicted: false, CycleDate: time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC), CountOrder: "105001007000", DispositionDate: time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC), PC290Registration: true}

		rows := []data.DOJRow{conviction1, nonConviction, conviction2, registration, conviction3, conviction4, conviction5, conviction5Prison}
		history = data.DOJHistory{}
		for _, row := range rows {
			history.PushRow(row, "SACRAMENTO")
		}
	})

	Describe("PushRow", func() {
		It("Sets the correct values on the history", func() {
			Expect(history.SubjectID).To(Equal("subj_id"))
			Expect(history.Name).To(Equal("SOUP,ZAK E"))
			Expect(history.WeakName).To(Equal("SOUP,ZAK"))
			Expect(history.CII).To(Equal("A012345678"))
			Expect(history.SSN).To(Equal("345678125"))
			Expect(history.DOB).To(Equal(birthDate))
			Expect(history.CDL).To(Equal("testcdl"))

			expectedConviction1 := conviction1
			expectedConviction2 := conviction2
			expectedConviction3 := conviction3
			expectedConviction4 := conviction4
			expectedConviction5 := conviction5

			expectedConviction1.HasProp64ChargeInCycle = true
			expectedConviction2.HasProp64ChargeInCycle = true
			expectedConviction3.HasProp64ChargeInCycle = false
			expectedConviction4.HasProp64ChargeInCycle = true
			expectedConviction5.HasProp64ChargeInCycle = true
			expectedConviction5.SentenceEndDate = time.Date(2012, 04, 03, 0, 0, 0, 0, time.UTC)

			Expect(history.Convictions).To(ConsistOf(
				&expectedConviction1,
				&expectedConviction2,
				&expectedConviction3,
				&expectedConviction4,
				&expectedConviction5,
			))

			Expect(history.Convictions).ToNot(ConsistOf(&conviction5Prison))
			Expect(history.Convictions[4].SentenceEndDate).To(Equal(time.Date(2012, 04, 03, 0, 0, 0, 0, time.UTC)))
		})
	})

	Describe("MostRecentConvictionDate", func() {
		It("returns the most recent conviction date", func() {
			Expect(history.MostRecentConvictionDate()).To(Equal(time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC)))
		})
	})

	Describe("SuperstrikeCodeSections", func() {
		Describe("codes between 187 and 220", func() {
			BeforeEach(func() {
				conviction6 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "187 PC", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "102001003300", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction7 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "191.5 PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "103001004300", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction8 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "187-664 PC", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), CountOrder: "104001005700", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
				conviction9 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "191.5-664 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction10 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "209 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006900", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction11 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "220 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006100", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction12 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "189 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006300", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction13 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "187a PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007600", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction14 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "191.55 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007700", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "219 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}

				extraRows := []data.DOJRow{conviction6, conviction7, conviction8, conviction9, conviction10, conviction11, conviction12, conviction13, conviction14, conviction15}

				for _, row := range extraRows {
					history.PushRow(row, "SACRAMENTO")
				}

				Expect(len(history.Convictions)).To(Equal(15))
			})

			It("returns code sections for matched codes", func() {
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("187 PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("191.5 PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("187-664 PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("191.5-664 PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("209 PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("220 PC"))
			})

			It("DOES NOT return code sections for unmatched codes", func() {
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("189 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("187a PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("191.55 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("219 PC"))
			})
		})

		Describe("codes between 245 and 269", func() {
			BeforeEach(func() {
				conviction6 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "245 PC", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "102001003300", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction7 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "245(D)(3) PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "103001004300", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction8 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "261 PC", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), CountOrder: "104001005700", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
				conviction9 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "261(A)(2) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction10 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "264.1 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006900", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction11 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "264.11 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006100", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction12 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "269 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006300", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction13 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "555 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007600", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction14 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "55 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007700", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "269.1 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}

				extraRows := []data.DOJRow{conviction6, conviction7, conviction8, conviction9, conviction10, conviction11, conviction12, conviction13, conviction14, conviction15}

				for _, row := range extraRows {
					history.PushRow(row, "SACRAMENTO")
				}

				Expect(len(history.Convictions)).To(Equal(15))
			})

			It("returns code sections for matched codes", func() {
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("245(D)(3) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("261(A)(2) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("264.1 PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("269 PC"))
			})

			It("DOES NOT return code sections for unmatched codes", func() {
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("245 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("261 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("264.11 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("555 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("55 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("269.1 PC"))
			})
		})

		Describe("286 codes", func() {
			BeforeEach(func() {
				conviction6 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "286(C)(2)(A) PC", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "102001003300", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction7 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "286(C)(1) PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "103001004300", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction8 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "286(C)(2)(B) PC", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), CountOrder: "104001005700", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
				conviction9 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(C)(2)(C) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction10 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(C)(3) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006900", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction11 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(D)(1) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006100", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction12 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(D)(2) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006300", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction13 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(D)(3) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007600", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction14 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(D)(1)(2) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007700", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "286(C)(2) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}

				extraRows := []data.DOJRow{conviction6, conviction7, conviction8, conviction9, conviction10, conviction11, conviction12, conviction13, conviction14, conviction15}

				for _, row := range extraRows {
					history.PushRow(row, "SACRAMENTO")
				}

				Expect(len(history.Convictions)).To(Equal(15))
			})

			It("returns code sections for matched codes", func() {
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(C)(2)(A) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(C)(1) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(C)(2)(B) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(C)(2)(C) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(C)(3) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(D)(1) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(D)(2) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("286(D)(3) PC"))
			})

			It("DOES NOT return code sections for unmatched codes", func() {
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("286(C)(2) PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("286(D)(1)(2) PC"))
			})
		})

		Describe("288 codes", func() {
			BeforeEach(func() {
				conviction6 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "288(A) PC", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "102001003300", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction7 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "288(B)(1) PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "103001004300", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction8 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "288(B)(2) PC", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), CountOrder: "104001005700", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
				conviction9 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288A(C)(1) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction10 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288A(C)(2)(A) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006900", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction11 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288A(C)(2)(B) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006100", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction12 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288A(C)(2)(C) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006300", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction13 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288A(D) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007600", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction14 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288.5(A) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007700", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288(B) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007900", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288(A)(C)(1) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001008000", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "288B PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001008100", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}

				extraRows := []data.DOJRow{conviction6, conviction7, conviction8, conviction9, conviction10, conviction11, conviction12, conviction13, conviction14, conviction15}

				for _, row := range extraRows {
					history.PushRow(row, "SACRAMENTO")
				}

				Expect(len(history.Convictions)).To(Equal(15))
			})

			It("returns code sections for matched codes", func() {
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288(A) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288(B)(1) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288(B)(2) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288A(C)(1) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288A(C)(2)(A) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288A(C)(2)(B) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288A(C)(2)(C) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288A(D) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("288.5(A) PC"))
			})

			It("DOES NOT return code sections for unmatched codes", func() {
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("288 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("288(B) PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("288(A)(C)(1) PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("288B PC"))
			})
		})

		Describe("289 codes", func() {
			BeforeEach(func() {
				conviction6 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1119999", DOB: birthDate, CodeSection: "289(A)(1)(A) PC", Convicted: true, CycleDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "102001003300", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction7 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1118888", DOB: birthDate, CodeSection: "289(A)(1)(B) PC", Convicted: true, CycleDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), CountOrder: "103001004300", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction8 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "289(A)(1)(C) PC", Convicted: true, CycleDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), CountOrder: "104001005700", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
				conviction9 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "289(A)(2)(C) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction10 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "289(J) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006900", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction11 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "653F PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006100", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction12 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "11418(A)(1) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001006300", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction13 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "653(F) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007600", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction14 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "289 PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007700", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
				conviction15 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", CDL: "testcdl", CII: "A012345678", SSN: "345678125", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "289(A)(2)(A) PC", Convicted: true, CycleDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), CountOrder: "104001007800", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}

				extraRows := []data.DOJRow{conviction6, conviction7, conviction8, conviction9, conviction10, conviction11, conviction12, conviction13, conviction14, conviction15}

				for _, row := range extraRows {
					history.PushRow(row, "SACRAMENTO")
				}

				Expect(len(history.Convictions)).To(Equal(15))
			})

			It("returns code sections for matched codes", func() {
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("289(A)(1)(A) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("289(A)(1)(B) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("289(A)(1)(C) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("289(A)(2)(C) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("289(J) PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("653F PC"))
				Expect(history.SuperstrikeCodeSections()).To(ContainElement("11418(A)(1) PC"))
			})

			It("DOES NOT return code sections for unmatched codes", func() {
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("653(F) PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("289 PC"))
				Expect(history.SuperstrikeCodeSections()).NotTo(ContainElement("289(A)(2)(A) PC"))
			})
		})
	})
})