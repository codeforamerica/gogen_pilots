package data_test

import (
	"encoding/csv"
	"os"
	path "path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gogen/data"
)

var _ = Describe("CMSInformation", func() {
	Describe("NewCMSInformation", func() {
		var pathToCMSFile string
		var expectedEntry CMSEntry
		var err error

		BeforeEach(func() {
			pathToCMSFile, err = path.Abs(path.Join("..", "test_fixtures", "felonies_sf.csv"))
			Expect(err).ToNot(HaveOccurred())

			expectedEntry = CMSEntry{}
		})

		It("returns a new weights information", func() {
			CMSCSV, err := os.Open(pathToCMSFile)
			Expect(err).ToNot(HaveOccurred())

			ci, err := NewCMSInformation(csv.NewReader(CMSCSV))
			Expect(err).ToNot(HaveOccurred())
			Expect(ci).ToNot(BeNil())
		})

		PIt("parses the csv and creates entries", func() {
			CMSCSV, err := os.Open(pathToCMSFile)
			Expect(err).ToNot(HaveOccurred())

			ci, err := NewCMSInformation(csv.NewReader(CMSCSV))
			Expect(err).ToNot(HaveOccurred())
			Expect(ci).ToNot(BeNil())
			Expect(ci.Entries).To(ContainElement(expectedEntry))
		})
	})
})
