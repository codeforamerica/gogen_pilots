package data_test

import (
	"gogen/data"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Subject", func() {
	var (
		subject           data.Subject
		conviction1       data.DOJRow
		conviction2       data.DOJRow
		conviction3       data.DOJRow
		conviction4       data.DOJRow
		conviction5       data.DOJRow
		conviction6       data.DOJRow
		conviction7       data.DOJRow
		conviction5Prison data.DOJRow
		nonConviction     data.DOJRow
		birthDate         time.Time
	)

	days := time.Duration(24) * (time.Hour)
	sacramentoEligibilityFlow := data.EligibilityFlows["SACRAMENTO"]

	BeforeEach(func() {
		birthDate = time.Date(1994, time.April, 10, 0, 0, 0, 0, time.UTC)
		conviction1 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1234", DOB: birthDate, CodeSection: "11357 HS", WasConvicted: true, CountOrder: "101001001000", DispositionDate: time.Date(1999, time.May, 4, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", NumCrtCase: "777CRTCASE"}
		nonConviction = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1235", DOB: birthDate, CodeSection: "11357 HS", WasConvicted: false,CountOrder: "101001002000", DispositionDate: time.Date(2008, time.April, 14, 0, 0, 0, 0, time.UTC)}
		conviction2 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1119999", DOB: birthDate, CodeSection: "286(D)(1) PC", WasConvicted: true, CountOrder: "102001003000", DispositionDate: time.Date(2009, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
		conviction3 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1118888", DOB: birthDate, CodeSection: "187 PC", WasConvicted: true, CountOrder: "103001004000", DispositionDate: time.Date(2001, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
		conviction4 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "11360 HS", WasConvicted: true,CountOrder: "104001005000", DispositionDate: time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
		conviction5 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "266J PC", WasConvicted: true, CountOrder: "104001006000", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentenceEndDate: time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC)}
		conviction5Prison = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1236 334455-00", DOB: birthDate, CodeSection: "11360 HS", WasConvicted: true, CountOrder: "104001006000", DispositionDate: time.Date(2009, time.December, 5, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO", SentencePartDuration: time.Duration(30 * days)}
		registration := data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1236 12345678-00", DOB: birthDate, CodeSection: "290 PC", WasConvicted: false, CountOrder: "105001007000", DispositionDate: time.Date(2008, time.June, 19, 0, 0, 0, 0, time.UTC), IsPC290Registration: true}

		rows := []data.DOJRow{conviction1, nonConviction, conviction2, registration, conviction3, conviction4, conviction5, conviction5Prison}
		subject = data.Subject{}
		for _, row := range rows {
			subject.PushRow(row, sacramentoEligibilityFlow)
		}
	})

	Describe("PushRow", func() {
		It("Sets the correct values on the subject", func() {
			Expect(subject.ID).To(Equal("subj_id"))
			Expect(subject.Name).To(Equal("SOUP,ZAK E"))
			Expect(subject.DOB).To(Equal(birthDate))

			expectedConviction1 := conviction1
			expectedConviction2 := conviction2
			expectedConviction3 := conviction3
			expectedConviction4 := conviction4
			expectedConviction5 := conviction5

			expectedConviction1.HasProp64ChargeInCycle = true
			expectedConviction2.HasProp64ChargeInCycle = false
			expectedConviction3.HasProp64ChargeInCycle = false
			expectedConviction4.HasProp64ChargeInCycle = true
			expectedConviction5.HasProp64ChargeInCycle = true
			expectedConviction5.SentenceEndDate = time.Date(2012, 04, 03, 0, 0, 0, 0, time.UTC)

			Expect(subject.Convictions).To(ConsistOf(
				&expectedConviction1,
				&expectedConviction2,
				&expectedConviction3,
				&expectedConviction4,
				&expectedConviction5,
			))

			Expect(subject.Convictions).ToNot(ConsistOf(&conviction5Prison))
			Expect(subject.Convictions[4].SentenceEndDate).To(Equal(time.Date(2012, 04, 03, 0, 0, 0, 0, time.UTC)))
		})
	})

	Describe("MostRecentConvictionDate", func() {
		It("returns the most recent conviction date", func() {
			Expect(subject.MostRecentConvictionDate()).To(Equal(time.Date(2011, time.May, 12, 0, 0, 0, 0, time.UTC)))
		})
	})

	Describe("NumberOfConvictionsInLast7Years", func() {
		Describe("when at least one conviction occurred within the last 7 years", func() {
			BeforeEach(func() {
				conviction6 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1119999", DOB: birthDate, CodeSection: "187 PC", WasConvicted: true, CountOrder: "102001003300", DispositionDate: time.Date(2016, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction7 = data.DOJRow{SubjectID: "subj_id", Name: "SOUP,ZAK E", OFN: "1118888", DOB: birthDate, CodeSection: "191.5 PC", WasConvicted: true, CountOrder: "103001004300", DispositionDate: time.Date(2017, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}

				subject.PushRow(conviction6, sacramentoEligibilityFlow)
				subject.PushRow(conviction7, sacramentoEligibilityFlow)
			})

			It("returns the number of convictions that occurred in the last 7 years", func() {
				Expect(subject.NumberOfConvictionsInLast7Years()).To(Equal(2))
			})
		})

		It("returns 0 if no convictions occurred in the last 7 years", func() {
			Expect(subject.NumberOfConvictionsInLast7Years()).To(Equal(0))
		})
	})

	Describe("SuperstrikeCodeSections", func() {
		It("returns a list of superstrikes in the subject", func() {
			Expect(subject.SuperstrikeCodeSections()).To(ConsistOf("286(D)(1) PC", "187 PC"))
		})
	})

	Describe("PC290CodeSections", func() {
		It("returns a list of code sections that fall under PC290 for the subject", func() {
			Expect(subject.PC290CodeSections()).To(ConsistOf("286(D)(1) PC", "266J PC"))
		})
	})
})
