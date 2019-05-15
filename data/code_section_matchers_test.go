package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gogen/data"
)

var _ = Describe("IsSuperstrike", func() {
	It("returns true if the code section is a superstrike", func() {
		validSuperstrikes := []string{
			"187 PC",
			"189 PC",
			"189.1 PC",
			"189.5 PC",
			"190 PC",
			"190.03 PC",
			"190.2(A)(4) PC",
			"190.25(B) PC",
			"191 PC",
			"191.5 PC",
			"205 PC",
			"207 PC",
			"209 PC",
			"209.5 PC",
			"217.1 PC",
			"218 PC",
			"219 PC",
			"220 PC",
			"245(D)(3) PC",
			"261(A)(2) PC",
			"262(A)(4) PC",
			"262 PC",
			"264.1 PC",
			"269 PC",
			"273AB PC",
			"273AB(A) PC",
			"286(C)(2)(A) PC",
			"286(C)(1) PC",
			"286(C)(2)(B) PC",
			"286(C)(2)(C) PC",
			"286(C)(3) PC",
			"286(D)(1) PC",
			"286(D)(2) PC",
			"286(D)(3) PC",
			"286 PC",
			"287 PC",
			"288 PC",
			"288(A) PC",
			"288(B)(1) PC",
			"288(B)(2) PC",
			"288A PC",
			"288A(C)(1) PC",
			"288A(C)(2)(A) PC",
			"288A(C)(2)(B) PC",
			"288A(C)(2)(C) PC",
			"288A(D) PC",
			"288.5 PC",
			"288.5(A) PC",
			"289(A)(1)(A) PC",
			"289(A)(1)(B) PC",
			"289(A)(1)(C) PC",
			"289(A)(2)(C) PC",
			"289(J) PC",
			"289 PC",
			"451.5 PC",
			"653F PC",
			"667.61 PC",
			"667.7 PC",
			"667.71 PC",
			"4500 PC",
			"11418(A)(1) PC",
			"11418(B)(1) PC",
			"11418(B)(2) PC",
			"12308 PC",
			"12310 PC",
			"18745 PC",
			"18755 PC",
			"1672(A) MV",
			//"187-664 PC",
			//"191.5-664 PC",
		}

		for _, validSuperstrike := range validSuperstrikes {
			Expect(data.IsSuperstrike(validSuperstrike)).To(BeTrue(), "Failed on example "+validSuperstrike)
		}
	})

	It("returns true if the code section is a superstrike with gang enhancement", func() {
		validSuperstrikes := []string{
			"246/186.22(B)(4) PC",
			"519-186.22(B)(4) PC",
			"136.1/186.22(B)(4) PC",
			"136.1-186.22(B)(4) PC",
			"136.1+186.22(B)(4) PC",
			"215 + 186.22(B)(4) PC",
			"12022.55 + 186.22(B)(4) PC",
			"213(A)(1)(A) + 186.22(B)(4) PC",
			"213(A)(1)(A)/186.22(B)(4) PC",
			"186.22(B)(4)/136.1 PC",
			"186.22(B)(4)-136.1 PC",
			"186.22(B)(4)+136.1 PC",
			"186.22(B)(4)-213(A)(1)(A) PC",
			"186.22(B)(4)/12022.55 PC",
			"186.22(B)(4)-246 PC",
		}

		for _, validSuperstrike := range validSuperstrikes {
			Expect(data.IsSuperstrike(validSuperstrike)).To(BeTrue(), "Failed on example "+validSuperstrike)
		}
	})

	It("returns true if the code section is a superstrike with a -664 attempted", func() {
		attemptedSuperstrikeCodeSection := "187-664 PC"

		Expect(data.IsSuperstrike(attemptedSuperstrikeCodeSection)).To(BeTrue(), "Failed on example ")
	})

	It("returns true if the code section is a superstrike with a -182 conspiracy", func() {
		attemptedSuperstrikeCodeSection := "182 + 190.2(A)(4) PC"

		Expect(data.IsSuperstrike(attemptedSuperstrikeCodeSection)).To(BeTrue(), "Failed on example ")
	})


	It("returns false is the code section is not a superstrike", func() {
		nonSuperstrikes := []string{
			"186.22(B)(3)+136.1 PC",
			"186.22(A)+136.1 PC",
			"186.22(B)(4) PC",
			"136.1 PC",
			"189.2 PC",
			"1900 PC",
			"190.123 PC",
			"191(A) PC",
			"191.55 PC",
			"205.1 PC",
			"209.6 PC",
			"246 PC",
			"12022.5 PC",
			"217.1(A) PC",
			"245 PC",
			"261.5 PC",
			"264.11 PC",
			"273 PC",
			"273(A) PC",
			"273A(B)",
			"286(C)(4)(A) PC",
			"286(E)(2)(B) PC",
			"286(C) PC",
			"288(B) PC",
			"288(A)(C)(1) PC",
			"288B PC",
			"289(A)(2)(A) PC",
			"451.1 PC",
			"653(F) PC",
			"667.6 PC",
			"667.72 PC",
			"667 PC",
			"667/664 PC",
			"4500(A) PC",
			"11418(A)(2) PC",
			"18745(A) PC",
			"1672(B) MV",
		}

		for _, nonSuperstrike := range nonSuperstrikes {
			Expect(data.IsSuperstrike(nonSuperstrike)).To(BeFalse(), "Failed on example "+nonSuperstrike)
		}
	})
})

var _ = Describe("IsPC290", func() {
	It("returns true if the code section falls under PC 290", func() {
		validPC290s := []string{
			"266 PC",
			"266C PC",
			"267 PC",
			"285 PC",
			"288 PC",
			"290 PC",
			"261 PC",
			"269 PC",
			"314 PC",
			"290(A) PC",
			"290.1 PC",
			"261B PC",
			"269 PC",
			"269.8 PC",
			"314(A)(2)(C)(1) PC",
			"236.1(B) PC",
			"236.1(C) PC",
			"236.1(C)(A) PC",
			"236.1(B)(C) PC",
			"236.1(B)(1) PC",
			"243.4(A) PC",
			"264.11 PC",
			"311.1(2) PC",
			"647.6B PC",
			"243.4(A)(C) PC",
			"266J(A) PC",
			"266J.11 PC",
			"647A(2) PC",
			"647AB PC",
			"647A.2 PC",
		}

		for _, validPC290 := range validPC290s {
			Expect(data.IsPC290(validPC290)).To(BeTrue(), "Failed on example "+validPC290)
		}
	})

	It("returns false is the code section does not fall under PC 290", func() {
		nonPC290s := []string{
			"266.5 PC",
			"266(C) PC",
			"267A PC",
			"2677 PC",
			"2699 PC",
			"262 PC",
			"291 PC",
			"236.1(A) PC",
			"236.1B PC",
			"236.2(B) PC",
			"236.1(D)(B) PC",
			"243 PC",
			"2434 PC",
			"24 PC",
			"647(A) PC",
			"266(J) PC",
			"647.1(A) PC",
			"6477A PC",
		}

		for _, nonPC290 := range nonPC290s {
			Expect(data.IsPC290(nonPC290)).To(BeFalse(), "Failed on example "+nonPC290)
		}
	})
})


var _ = Describe("StripFlags", func() {
	It("strips the trailing 182 and punctuation and replaces it with a space", func() {
		validConspiredSuperstrikes := []string{
			"187-182 PC",
			"187 182 PC",
			"189/182 PC",
			"189.1+182 PC",
			"190.25(B) + 182 PC",
			"191 / 182 PC",
			"191.5 - 182 PC",
			"37/182 PC",
			"37 182 PC",
		}
		StrippedSuperstrikes := []string{
			"187 PC",
			"187 PC",
			"189 PC",
			"189.1 PC",
			"190.25(B) PC",
			"191 PC",
			"191.5 PC",
			"37 PC",
			"37 PC",
		}

		for i, validConspiredSuperstrike := range validConspiredSuperstrikes {
			Expect(data.StripFlags(validConspiredSuperstrike, `182`)).To(Equal(StrippedSuperstrikes[i]), "Failed on example "+validConspiredSuperstrike)
		}
	})

	It("strips the prepended 182 and punctuation and replaces it with a space", func() {
		validConspiredSuperstrikes := []string{
			"182-189.5 PC",
			"182/190 PC",
			"182+190.03 PC",
			"182 + 190.2(A)(4) PC",
		}
		StrippedSuperstrikes := []string{
			" 189.5 PC",
			" 190 PC",
			" 190.03 PC",
			" 190.2(A)(4) PC",
		}

		for i, validConspiredSuperstrike := range validConspiredSuperstrikes {
			Expect(data.StripFlags(validConspiredSuperstrike, `182`)).To(Equal(StrippedSuperstrikes[i]), "Failed on example "+validConspiredSuperstrike)
		}
	})

	It("strips the trailing 664 and punctuation and replaces it with a space", func() {
		validAttemptedSuperstrikes := []string{
			"187-664 PC",
			"189/664 PC",
			"189.1+664 PC",
			"189.1 664 PC",
			"190.25(B) + 664 PC",
			"191 / 664 PC",
			"191.5 - 664 PC",
			"37/664 PC",
			"37 664 PC",
		}
		StrippedSuperstrikes := []string{
			"187 PC",
			"189 PC",
			"189.1 PC",
			"189.1 PC",
			"190.25(B) PC",
			"191 PC",
			"191.5 PC",
			"37 PC",
			"37 PC",
		}

		for i, validAttemptedSuperstrike := range validAttemptedSuperstrikes {
			Expect(data.StripFlags(validAttemptedSuperstrike, `664`)).To(Equal(StrippedSuperstrikes[i]), "Failed on example "+validAttemptedSuperstrike)
		}
	})

	It("strips the prepended 664 and punctuation and replaces it with a space", func() {
		validAttemptedSuperstrikes := []string{
			"664-189.5 PC",
			"664/190 PC",
			"664+190.03 PC",
			"664 + 190.2(A)(4) PC",
		}
		StrippedSuperstrikes := []string{
			" 189.5 PC",
			" 190 PC",
			" 190.03 PC",
			" 190.2(A)(4) PC",
		}

		for i, validAttemptedSuperstrike := range validAttemptedSuperstrikes {
			Expect(data.StripFlags(validAttemptedSuperstrike, `664`)).To(Equal(StrippedSuperstrikes[i]), "Failed on example "+validAttemptedSuperstrike)
		}
	})
})
