package processor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gogen/processor"

	"gogen/data"
)

var _ = Describe("EligibilityChecker", func() {
	var (
		entry      data.CMSEntry
		weightInfo data.WeightsEntry
	)

	BeforeEach(func() {
		weightInfo = data.WeightsEntry{
			Weight: 54.0,
			Found:  true,
		}
	})

	It("Checks for weight disqualifiers", func() {
		info := ComputeEligibility(entry, weightInfo)

		Expect(info.QFinalSum).To(Equal(54.0))
		Expect(info.Over1Lb).To(Equal("eligible"))
	})

	Context("a weight entry was not found", func() {
		BeforeEach(func() {
			weightInfo = data.WeightsEntry{
				Weight: 0,
				Found:  false,
			}
		})

		It("reports the not found weights entry", func() {
			info := ComputeEligibility(entry, weightInfo)

			Expect(info.QFinalSum).To(Equal(0.0))
			Expect(info.Over1Lb).To(Equal("not found"))
		})
	})

	Context("The CMSEntry is an 11357 charge", func() {
		BeforeEach(func() {
			weightInfo = data.WeightsEntry{
				Weight: 123.4,
				Found:  true,
			}
			entry = data.CMSEntry{
				Level: "M",
			}
		})

		It("reports the not found weights entry", func() {
			info := ComputeEligibility(entry, weightInfo)

			Expect(info.QFinalSum).To(Equal(123.4))
			Expect(info.Over1Lb).To(Equal("n/a"))
		})
	})
})
