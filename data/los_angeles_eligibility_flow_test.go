package data_test

import (
	"gogen/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("losAngelesEligibilityFlow", func() {

	var flow data.EligibilityFlow

	BeforeEach(func() {
		flow = data.EligibilityFlows["LOS ANGELES"]
	})

	Describe("MatchedCodeSection", func() {
		It("returns the matched substring for a given code section", func() {
			Expect(flow.MatchedCodeSection("11358(c) HS")).To(Equal("11358"))
		})

		It("returns empty string if there is no match", func() {
			Expect(flow.MatchedCodeSection("12345(c) HS")).To(Equal(""))
		})

		It("recognizes attempted code sections for Prop 64", func() {
			Expect(flow.MatchedCodeSection("664.11357(c) HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("66411357(c) HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("664-11357(c) HS")).To(Equal("11357"))
			Expect(flow.MatchedCodeSection("664/11357(c) HS")).To(Equal("11357"))
		})
	})
})