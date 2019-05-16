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
			"236.1(B) PC",
			"236.1(B)(C) PC",
			"236.1(B)(1) PC",
			"236.1(C) PC",
			"236.1(C)(A) PC",
			"243.4(A) PC",
			"243.4(A)(C) PC",
			"261 PC",
			"261B PC",
			"262(A) PC",
			"262(A)(1) PC",
			"264.1 PC",
			"266 PC",
			"266C PC",
			"266H(B) PC",
			"267 PC",
			"269 PC",
			"269.8 PC",
			"272 PC",
			"285 PC",
			"286 PC",
			"287 PC",
			"288 PC",
			"288A PC",
			"288.2 PC",
			"288.3 PC",
			"288.4 PC",
			"288.5 PC",
			"288.7 PC",
			"289 PC",
			"311.1 PC",
			"311.2(B) PC",
			"311.2(C) PC",
			"311.2(D) PC",
			"311.3 PC",
			"311.4 PC",
			"311.10 PC",
			"311.11 PC",
			"314(1) PC",
			"314(2) PC",
			"451.5 PC",
			"647.6B PC",
			"647A PC",
			"653F(B) PC",
			"653F(C) PC",
		}

		for _, validPC290 := range validPC290s {
			Expect(data.IsPC290(validPC290)).To(BeTrue(), "Failed on example "+validPC290)
		}
	})

	It("returns false is the code section does not fall under PC 290", func() {
		nonPC290s := []string{
			"236.1(A) PC",
			"236.1B PC",
			"236.2(B) PC",
			"236.1(D)(B) PC",
			"243 PC",
			"262 PC",
			"264.11 PC",
			"266.5 PC",
			"266(C) PC",
			"266J(A) PC",
			"266J.11 PC",
			"266(J) PC",
			"267A PC",
			"2677 PC",
			"2699 PC",
			"288.6",
			"291 PC",
			"243 PC",
			"2434 PC",
			"24 PC",
			"647(A) PC",
			"647.1(A) PC",
			"6477A PC",
			"653F PC",
		}

		for _, nonPC290 := range nonPC290s {
			Expect(data.IsPC290(nonPC290)).To(BeFalse(), "Failed on example "+nonPC290)
		}
	})
})
