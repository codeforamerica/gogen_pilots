package data_test

import (
	"gogen/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sanJoaquinEligibilityFlow", func() {

	var flow data.EligibilityFlow

	BeforeEach(func() {
		flow = data.EligibilityFlows["SAN JOAQUIN"]
	})

	Describe("MatchedCodeSection", func() {
		It("returns the matched substring for a given code section", func() {
			Expect(flow.MatchedCodeSection("11358(c) HS")).To(Equal("11358"))
			Expect(flow.MatchedCodeSection("/11357 HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("647(f) PC")).To(Equal("647(f) PC"))
			Expect(flow.MatchedCodeSection("148.9 PC")).To(Equal("148.9 PC"))
			Expect(flow.MatchedCodeSection("4060    BP")).To(Equal("4060    BP"))
			Expect(flow.MatchedCodeSection("--40508 VC--")).To(Equal("40508 VC"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedCodeSection("12345(c) HS")).To(Equal(""))
			Expect(flow.MatchedCodeSection("647(f) HS")).To(Equal(""))
			Expect(flow.MatchedCodeSection("4050.6 BP")).To(Equal(""))
			Expect(flow.MatchedCodeSection("14859 PC")).To(Equal(""))

		})
	})
})
