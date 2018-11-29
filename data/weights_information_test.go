package data_test

import (
	"encoding/csv"
	"os"
	path "path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gogen/data"
)

var _ = Describe("WeightsInformation", func() {
	Describe("NewWeightsInformation", func() {
		var pathToWeights string
		var err error

		BeforeEach(func() {
			pathToWeights, err = path.Abs(path.Join("..", "test_fixtures", "conviction_weights.csv"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a new weights information", func() {
			weightsCSV, err := os.Open(pathToWeights)
			Expect(err).ToNot(HaveOccurred())

			wi, err := NewWeightsInformation(csv.NewReader(weightsCSV))
			Expect(err).ToNot(HaveOccurred())
			Expect(wi).ToNot(BeNil())
		})
	})

	Describe("#GetWeight", func() {
		var pathToWeights string
		var weightsInformation *WeightsInformation
		var err error

		BeforeEach(func() {
			pathToWeights, err = path.Abs(path.Join("..", "test_fixtures", "conviction_weights.csv"))
			Expect(err).ToNot(HaveOccurred())

			weightsCSV, err := os.Open(pathToWeights)
			Expect(err).ToNot(HaveOccurred())

			weightsInformation, err = NewWeightsInformation(csv.NewReader(weightsCSV))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Knows about convictions under 1 pound", func() {
			val := weightsInformation.GetWeight("305599")
			Expect(val.Weight).To(Equal(54.0))
			Expect(val.Found).To(Equal(true))
		})

		Context("When the key does not exist", func() {
			It("returns an error", func() {
				val := weightsInformation.GetWeight("thing")
				Expect(val.Weight).To(Equal(0.0))
				Expect(val.Found).To(Equal(false))
			})
		})
	})
})
