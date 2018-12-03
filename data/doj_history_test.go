package data_test

import (
	"gogen/data"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DOJHistory", func() {

	Describe("PushRow", func() {
		Context("An empty history", func() {
			var history data.DOJHistory
			var conviction1 data.DOJRow
			var conviction2 data.DOJRow
			var conviction3 data.DOJRow
			var conviction4 data.DOJRow
			var nonConviction data.DOJRow

			BeforeEach(func() {
				conviction1 = data.DOJRow{SubjectID: "subj_id", Name: "soup, zak e", CDL: "testcdl", CII: "12345678", SSN: "345678125", OFN: "1234", DOB: time.Time{}, CodeSection: "11357 HS", Convicted: true, DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC)}
				nonConviction = data.DOJRow{SubjectID: "subj_id", Name: "soup, zak e", CDL: "testcdl", CII: "12345678", SSN: "345678125", OFN: "1235", DOB: time.Time{}, CodeSection: "11357 HS", Convicted: false, DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC)}
				conviction2 = data.DOJRow{SubjectID: "subj_id", Name: "soup, zak e", CDL: "testcdl", CII: "12345678", SSN: "345678125", OFN: "1119999", DOB: time.Time{}, CodeSection: "11360 HS", Convicted: true, DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC), County: "LOS ANGELES"}
				conviction3 = data.DOJRow{SubjectID: "subj_id", Name: "soup, zak e", CDL: "testcdl", CII: "12345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: time.Time{}, CodeSection: "11360 HS", Convicted: true, DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC), County: "SAN FRANCISCO"}
				conviction4 = data.DOJRow{SubjectID: "subj_id", Name: "soup, zak e", CDL: "testcdl", CII: "12345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: time.Time{}, CodeSection: "11360 HS", Convicted: true, DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC)}
				registration := data.DOJRow{SubjectID: "subj_id", Name: "soup, zak e", CDL: "testcdl", CII: "12345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: time.Time{}, CodeSection: "290 PC", Convicted: false, DispositionDate: time.Date(2008, time.May, 4, 0, 0, 0, 0, time.UTC), PC290Registration: true}

				rows := []data.DOJRow{conviction1, nonConviction, conviction2, registration, conviction3, conviction4}
				history = data.DOJHistory{}
				for _, row := range rows {
					history.PushRow(row)
				}
			})

			It("Sets the correct values on the history", func() {
				Expect(history.SubjectID).To(Equal("subj_id"))
				Expect(history.Name).To(Equal("soup, zak e"))
				Expect(history.CII).To(Equal("12345678"))
				Expect(history.SSN).To(Equal("345678125"))
				Expect(history.DOB).To(Equal(time.Time{}))
				Expect(history.CDL).To(Equal("testcdl"))
				Expect(history.PC290Registration).To(BeTrue())
				Expect(history.Convictions).To(ConsistOf(&conviction1, &conviction2, &conviction3, &conviction4))
			})
		})
	})
})
