package data_test

import (
	"gogen/data"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DOJHistory", func() {


	Describe("match", func() {
		Context("An empty history", func() {
			var history data.DOJHistory

			BeforeEach(func() {
				rows := []data.DOJRow{
					{SubjectID: "subj_id", Name: "soup, zak e", CII: "12345678", SSN: "345678125", OFN: "1234", DOB: time.Time{}, CodeSection: "11357 HS", Convicted: true, DispositionDate: time.Date(2008,time.May,4, 0,0,0,0,nil)},
					{SubjectID: "subj_id", Name: "soup, zak e", CII: "12345678", SSN: "345678125", OFN: "1235", DOB: time.Time{}, CodeSection: "11357 HS", Convicted: false, DispositionDate: time.Date(2008,time.May,4, 0,0,0,0,nil) },
					{SubjectID: "subj_id", Name: "soup, zak e", CII: "12345678", SSN: "345678125", OFN: "1119999", DOB: time.Time{}, CodeSection: "11360 HS", Convicted: true, DispositionDate: time.Date(2008,time.May,4, 0,0,0,0,nil), County: "LOS ANGELES" },
					{SubjectID: "subj_id", Name: "soup, zak e", CII: "12345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: time.Time{}, CodeSection: "11360 HS", Convicted: true, DispositionDate: time.Date(2008,time.May,4, 0,0,0,0,nil), County: "SAN FRANCISCO" },
					{SubjectID: "subj_id", Name: "soup, zak e", CII: "12345678", SSN: "345678125", OFN: "1236 12345678-00", DOB: time.Time{}, CodeSection: "11360 HS", Convicted: true, DispositionDate: time.Date(2008,time.May,4, 0,0,0,0,nil) },
				}
				history = data.DOJHistory{}
				for _, row := range rows {
					history.PushRow(row)
				}
			})

			PIt("Sets the correct values on the history", func() {
				Expect(history.SubjectID).To(Equal("subj_id"))
				Expect(history.Name).To(Equal("soup, zak e"))
				Expect(history.CII).To(Equal("12345678"))
				Expect(history.SSN).To(Equal("345678125"))
				Expect(history.DOB).To(Equal(time.Time{}))
			})
		})
	})
})
