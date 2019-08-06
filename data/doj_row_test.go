package data_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gogen/data"
)

var _ = Describe("DojRow", func() {
	var rawRow []string

	BeforeEach(func() {
		rawRow = []string{
			"x", "x", "18675309", "#", "1008675309", "x", "x", "x", "x", "x", "#", "1008675309", "SKYWALKER,LUKE S", "x", "19600314", "123456789", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "19790525", "x", "ARREST/DETAINED/CITED", "x", "x", "x", "CAPDSAN FRANCISCO", "x", "SAN FRANCISCO", "x", "x", "19790525", "12 140189-B", "x", "503 VC-TAKE CAR W/OUT OWNERS CONSENT", "F", "              ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "REL/TOT OTHER JURIS/AUTH", "", "FELONY", "#", "", "", "", "", "", "                  ", "23", "", "", "", "#", "",
		}
	})

	It("Sets values on initialization", func() {
		expectedDob := time.Date(1960, time.March, 14, 0, 0, 0, 0, time.UTC)

		row := NewDOJRow(rawRow, 1)
		Expect(row.Name).To(Equal("SKYWALKER,LUKE S"))
		Expect(row.SubjectID).To(Equal("18675309"))
		Expect(row.DOB).To(Equal(expectedDob))
		Expect(row.IsPC290Registration).To(BeFalse())
		Expect(row.IsFelony).To(BeTrue())
	})

	Context("The row is a registration event", func() {
		BeforeEach(func() {
			rawRow = []string{
				"x", "x", "17954908", "#", "8690594867", "x", "x", "x", "x", "x", "#", "8690594867", "BIRD,BIG", "        x", "19850822", "987654321", "F1234567", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "19790601", "x", "REGISTRATION", "         x", "x", "x", "CASCSAN FRANCISCO", "x", "SAN FRANCISCO", "x", "x", "19790601", "678544", "x", "290 PC-REGISTRATION OF SEX OFFENDER", "", "                ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "                        ", "", "#", "", "", "", "", "", "																		", "  ", "", "", "", "#", "",
			}
		})

		It("recognizes the registration", func() {
			row := NewDOJRow(rawRow, 1)

			Expect(row.IsPC290Registration).To(BeTrue())
		})
	})

	Describe("OccurredInLast7Years", func() {
		Context("when the disposition date occurred in the last 7 years", func() {
			BeforeEach(func() {
				rawRow = []string{
					"x", "x", "18675309", "#", "1008675309", "x", "x", "x", "x", "x", "#", "1008675309", "SKYWALKER,LUKE S", "x", "19600314", "123456789", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "20160525", "x", "ARREST/DETAINED/CITED", "x", "x", "x", "CAPDSAN FRANCISCO", "x", "SAN FRANCISCO", "x", "x", "20160525", "12 140189-B", "x", "503 VC-TAKE CAR W/OUT OWNERS CONSENT", "F", "              ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "REL/TOT OTHER JURIS/AUTH", "", "FELONY", "#", "", "", "", "", "", "                  ", "23", "", "", "", "#", "",
				}
			})
			It("returns true", func() {
				row := NewDOJRow(rawRow, 1)

				Expect(row.DispositionDate).To(Equal(time.Date(2016, time.May, 25, 0, 0, 0, 0, time.UTC)))
				Expect(row.OccurredInLast7Years()).To(BeTrue())
			})
		})

		Context("when the disposition date occurred in the last 7 years", func() {
			BeforeEach(func() {
				rawRow = []string{
					"x", "x", "18675309", "#", "1008675309", "x", "x", "x", "x", "x", "#", "1008675309", "SKYWALKER,LUKE S", "x", "19600314", "123456789", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "x", "19790525", "x", "ARREST/DETAINED/CITED", "x", "x", "x", "CAPDSAN FRANCISCO", "x", "SAN FRANCISCO", "x", "x", "19790525", "12 140189-B", "x", "503 VC-TAKE CAR W/OUT OWNERS CONSENT", "F", "              ", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "REL/TOT OTHER JURIS/AUTH", "", "FELONY", "#", "", "", "", "", "", "                  ", "23", "", "", "", "#", "",
				}
			})
			It("returns true", func() {
				row := NewDOJRow(rawRow, 1)

				Expect(row.DispositionDate).To(Equal(time.Date(1979, time.May, 25, 0, 0, 0, 0, time.UTC)))
				Expect(row.OccurredInLast7Years()).To(BeFalse())
			})
		})
	})

	Describe("Determines the code section", func() {

		It("detects the code section when it is explicitly specified in OFFENSE_DESCR", func() {
			rawRow[OFFENSE_DESCR] = "503 VC-TAKE CAR W/OUT OWNERS CONSENT"
			row := NewDOJRow(rawRow, 1)
			Expect(row.CodeSection).To(Equal("503 VC"))
		})

		It("detects the code section from COMMENT_TEXT when OFFENSE_DESCR reads 'SEE COMMENT FOR CHANGE'", func() {
			rawRow[OFFENSE_DESCR] = "SEE COMMENT FOR CHARGE"
			rawRow[COMMENT_TEXT] = "503 VC-TAKE CAR W/OUT OWNERS CONSENT"
			row := NewDOJRow(rawRow, 1)
			Expect(row.CodeSection).To(Equal("503 VC"))
		})

		It("detects the code section from COMMENT_TEXT when OFFENSE_DESCR is empty", func() {
			rawRow[OFFENSE_DESCR] = ""
			rawRow[COMMENT_TEXT] = "503 VC-TAKE CAR W/OUT OWNERS CONSENT"
			row := NewDOJRow(rawRow, 1)
			Expect(row.CodeSection).To(Equal("503 VC"))
		})

		It("detects the code section from COMMENT_TEXT when OFFENSE_DESCR is all blank", func() {
			rawRow[OFFENSE_DESCR] = "   "
			rawRow[COMMENT_TEXT] = "503 VC-TAKE CAR W/OUT OWNERS CONSENT"
			row := NewDOJRow(rawRow, 1)
			Expect(row.CodeSection).To(Equal("503 VC"))
		})

		It("doesn't detect code section when COMMENT_TEXT and OFFENSE_DESCR are empty", func() {
			rawRow[OFFENSE_DESCR] = ""
			rawRow[COMMENT_TEXT] = ""
			row := NewDOJRow(rawRow, 1)
			Expect(row.CodeSection).To(Equal(""))
		})
	})
})
