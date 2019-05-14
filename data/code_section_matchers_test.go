package data_test

import (
	"gogen/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsSuperstrike", func() {
	It("returns true if the code section is a superstrike", func() {
		validSuperstrikes := []string{
			"187 PC",
			"191.5 PC",
			"187-664 PC",
			"191.5-664 PC",
			"209 PC",
			"220 PC",
			"245(D)(3) PC",
			"261(A)(2) PC",
			"264.1 PC",
			"269 PC",
			"286(C)(2)(A) PC",
			"286(C)(1) PC",
			"286(C)(2)(B) PC",
			"286(C)(2)(C) PC",
			"286(C)(3) PC",
			"286(D)(1) PC",
			"286(D)(2) PC",
			"286(D)(3) PC",
			"288(A) PC",
			"288(B)(1) PC",
			"288(B)(2) PC",
			"288A(C)(1) PC",
			"288A(C)(2)(A) PC",
			"288A(C)(2)(B) PC",
			"288A(C)(2)(C) PC",
			"288A(D) PC",
			"288.5(A) PC",
			"289(A)(1)(A) PC",
			"289(A)(1)(B) PC",
			"289(A)(1)(C) PC",
			"289(A)(2)(C) PC",
			"289(J) PC",
			"653F PC",
			"11418(A)(1) PC",
		}

		for _, validSuperstrike := range validSuperstrikes {
			Expect(data.IsSuperstrike(validSuperstrike)).To(BeTrue(), "Failed on example "+validSuperstrike)
		}
	})

	It("returns false is the code section is not a superstrike", func() {
		nonSuperstrikes := []string{
			"189 PC",
			"187A PC",
			"191.55 PC",
			"219 PC",
			"245 PC",
			"261 PC",
			"264.11 PC",
			"555 PC",
			"55 PC",
			"269.1 PC",
			"286(C)(2) PC",
			"286(D)(1)(2) PC",
			"288 PC",
			"288(B) PC",
			"288(A)(C)(1) PC",
			"288B PC",
			"653(F) PC",
			"289 PC",
			"289(A)(2)(A) PC",
		}

		for _, nonSuperstrike := range nonSuperstrikes {
			Expect(data.IsSuperstrike(nonSuperstrike)).To(BeFalse(), "Failed on example "+nonSuperstrike)
		}
	})
})
