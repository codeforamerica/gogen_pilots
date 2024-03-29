package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gogen_pilots/matchers"
	"testing"
)

func TestMatchers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Matchers Suite")
}

func getMatchedCodeSection(codeSection string) string {
	_, section := matchers.ExtractProp64Section(codeSection)
	return section
}

func getMatchedRelatedChargeCodeSection(codeSection string) string {
	_, section := matchers.ExtractRelatedChargeSection(codeSection)
	return section
}

func getMatched11357SubSection(codeSection string) string {
	_, section := matchers.Extract11357SubSection(codeSection)
	return section
}

var _ = Describe("MatchedCodeSection", func() {
	It("returns the matched substring for a given Prop 64 code section", func() {
		Expect(getMatchedCodeSection("11358(c) HS")).To(Equal("11358"))

		Expect(getMatchedCodeSection("/11357 HS")).To(Equal("11357"))
	})

	It("returns empty string if there is no match", func() {
		Expect(getMatchedCodeSection("12345(c) HS")).To(Equal(""))
		Expect(getMatchedCodeSection("647(f) HS")).To(Equal(""))
		Expect(getMatchedCodeSection("4050.6 BP")).To(Equal(""))
		Expect(getMatchedCodeSection("14859 PC")).To(Equal(""))
	})

	It("returns empty string if the code section is for a related charge", func() {
		Expect(getMatchedCodeSection("647(f) PC")).To(Equal(""))
		Expect(getMatchedCodeSection("148.9 PC")).To(Equal(""))
		Expect(getMatchedCodeSection("4060    BP")).To(Equal(""))
		Expect(getMatchedCodeSection("--40508 VC--")).To(Equal(""))
		Expect(getMatchedCodeSection("1320(a) PC")).To(Equal(""))
	})

	It("recognizes attempted code sections for Prop 64", func() {
		Expect(getMatchedCodeSection("664.11357(c) HS")).To(Equal("11357"))
		Expect(getMatchedCodeSection("66411357(c) HS")).To(Equal("11357"))
		Expect(getMatchedCodeSection("664-11357(c) HS")).To(Equal("11357"))
		Expect(getMatchedCodeSection("664/11357(c) HS")).To(Equal("11357"))
	})
})

var _ = Describe("MatchedRelatedCodeSection", func() {
	It("returns the matched substring for a given related charge code section", func() {
		Expect(getMatchedRelatedChargeCodeSection("647(f) PC")).To(Equal("647(f) PC"))
		Expect(getMatchedRelatedChargeCodeSection("148.9 PC")).To(Equal("148.9 PC"))
		Expect(getMatchedRelatedChargeCodeSection("4060    BP")).To(Equal("4060    BP"))
		Expect(getMatchedRelatedChargeCodeSection("--40508 VC--")).To(Equal("40508 VC"))
		Expect(getMatchedRelatedChargeCodeSection("1320(a) PC")).To(Equal("1320(a) PC"))
	})

	It("returns empty string if there is no match", func() {
		Expect(getMatchedRelatedChargeCodeSection("12345(c) HS")).To(Equal(""))
		Expect(getMatchedRelatedChargeCodeSection("647(f) HS")).To(Equal(""))
		Expect(getMatchedRelatedChargeCodeSection("4050.6 BP")).To(Equal(""))
		Expect(getMatchedRelatedChargeCodeSection("14859 PC")).To(Equal(""))
	})

	It("returns empty string if the code section is for a Prop 64 charge", func() {
		Expect(getMatchedRelatedChargeCodeSection("11358(c) HS")).To(Equal(""))
		Expect(getMatchedRelatedChargeCodeSection("/11357 HS")).To(Equal(""))
	})
})

var _ = Describe("Matched11357SubSection", func() {
	It("returns the matched subsection for a given 11357 code section", func() {
		Expect(getMatched11357SubSection("11357(A)")).To(Equal("A"))
		Expect(getMatched11357SubSection("11357(C)")).To(Equal("C"))
		Expect(getMatched11357SubSection("Some Prefix 11357(C) Some Suffix")).To(Equal("C"))
	})

	It("returns empty string if there is no match", func() {
		Expect(getMatched11357SubSection("11357")).To(Equal(""))
		Expect(getMatched11357SubSection("647(f) HS")).To(Equal(""))
	})
})