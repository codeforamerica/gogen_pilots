package data_test

import (
	"gogen/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("contraCostaEligibilityFlow", func() {

	var flow data.EligibilityFlow

	BeforeEach(func() {
		flow = data.EligibilityFlows["CONTRA COSTA"]
	})

	Describe("MatchedCodeSection", func() {
		It("returns the matched substring for a given Prop 64 code section", func() {
			Expect(flow.MatchedCodeSection("11358(c) HS")).To(Equal("11358"))
			Expect(flow.MatchedCodeSection("/11357 HS")).To(Equal("11357"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedCodeSection("12345(c) HS")).To(Equal(""))
			Expect(flow.MatchedCodeSection("647(f) HS")).To(Equal(""))
			Expect(flow.MatchedCodeSection("4050.6 BP")).To(Equal(""))
			Expect(flow.MatchedCodeSection("14859 PC")).To(Equal(""))
		})

		It("returns empty string if the code section is for a related charge", func() {
			Expect(flow.MatchedCodeSection("647(f) PC")).To(Equal(""))
			Expect(flow.MatchedCodeSection("148.9 PC")).To(Equal(""))
			Expect(flow.MatchedCodeSection("4060    BP")).To(Equal(""))
			Expect(flow.MatchedCodeSection("--40508 VC--")).To(Equal(""))
			Expect(flow.MatchedCodeSection("1320(a) PC")).To(Equal(""))
		})
	})

	Describe("MatchedRelatedCodeSection", func() {
		It("returns the matched substring for a given related charge code section", func() {
			Expect(flow.MatchedRelatedCodeSection("647(f) PC")).To(Equal("647(f) PC"))
			Expect(flow.MatchedRelatedCodeSection("148.9 PC")).To(Equal("148.9 PC"))
			Expect(flow.MatchedRelatedCodeSection("4060    BP")).To(Equal("4060    BP"))
			Expect(flow.MatchedRelatedCodeSection("--40508 VC--")).To(Equal("40508 VC"))
			Expect(flow.MatchedRelatedCodeSection("1320(a) PC")).To(Equal("1320(a) PC"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedRelatedCodeSection("12345(c) HS")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("647(f) HS")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("4050.6 BP")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("14859 PC")).To(Equal(""))
		})

		It("returns empty string if the code section is for a Prop 64 charge", func() {
			Expect(flow.MatchedRelatedCodeSection("11358(c) HS")).To(Equal(""))
			Expect(flow.MatchedRelatedCodeSection("/11357 HS")).To(Equal(""))
		})
	})
})
