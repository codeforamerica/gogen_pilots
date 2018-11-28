package data_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gogen/data"
)

var _ = Describe("CMSEntry", func() {
	var (
		record []string
		ce     CMSEntry
	)

	BeforeEach(func() {
		record = []string{"305563", "", "A1567564", "BIRD/BIG", "MISD", "190", "COUNTY JAIL W/ PROBATION CONDITION  ", "4/20/99", "M66654", "          ", " ", "      ", "11357(C)HS", "M", "", "190", "COUNTY JAIL W/ PROBATION CONDITION   ", "      ", "", "", "9/14/65", "S554423", "A123456780", "", "123456789", "F1234567 CA", "EOR"}
		ce = NewCMSEntry(record)
	})

	Describe("NewCMSEntry", func() {
		It("Returns a CMSEntry", func() {
			Expect(ce).ToNot(BeNil())
			Expect(ce.CourtNumber).To(Equal("305563"))
			Expect(ce.Level).To(Equal("M"))
			Expect(ce.SSN).To(Equal("123456789"))
			Expect(ce.CII).To(Equal("A123456780"))
			Expect(ce.Charge).To(Equal("11357(C)HS"))
			Expect(ce.IncidentNumber).To(Equal("A1567564"))
			Expect(ce.Name).To(Equal("BIRD/BIG"))
			Expect(ce.DateOfBirth).To(Equal(time.Date(1965, time.September, 14, 0, 0, 0, 0, time.UTC)))
		})

		It("Parses states out of DL numbers", func() {
			Expect(ce.CDL).To(Equal("F1234567"))
		})

		Describe("#FormattedName", func() {
			It("Formats the name", func() {
				Expect(ce.FormattedName()).To(Equal("BIRD,BIG"))
			})

			Context("There is a middle name", func() {
				BeforeEach(func() {
					record = []string{"305563", "", "A1567564", "BIRD/BIG/FLAPPY/YELLOW", "MISD", "190", "COUNTY JAIL W/ PROBATION CONDITION  ", "4/20/99", "M66654", "          ", " ", "      ", "11357(C)HS", "M", "", "190", "COUNTY JAIL W/ PROBATION CONDITION   ", "      ", "", "", "9/14/65", "S554423", "A123456780", "", "123456789", "F1234567 CA", "EOR"}
					ce = NewCMSEntry(record)
				})

				It("Formats the name", func() {
					Expect(ce.FormattedName()).To(Equal("BIRD,BIG FLAPPY YELLOW"))
				})
			})
		})

		Context("There is whitespace in the columns", func() {
			BeforeEach(func() {
				record = []string{"305563", "", "A1567564", "BIRD/BIG        ", "MISD", "190", "COUNTY JAIL W/ PROBATION CONDITION  ", "4/20/99", "M66654", "          ", " ", "      ", "	11357(C)HS     ", "M", "", "190", "COUNTY JAIL W/ PROBATION CONDITION   ", "      ", "", "", "9/14/65", "S554423", "A123456780", "", "123456789", "F1234567 CA", "EOR"}
				ce = NewCMSEntry(record)
			})

			It("Trims the whitepsace", func() {
				Expect(ce).ToNot(BeNil())
				Expect(ce.Charge).To(Equal("11357(C)HS"))
				Expect(ce.Name).To(Equal("BIRD/BIG"))
			})
		})
	})
})
