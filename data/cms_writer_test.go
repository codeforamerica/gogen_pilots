package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gogen/data"

	"io/ioutil"
	path "path/filepath"
)

var _ = Describe("csvWriter", func() {
	var writer CMSWriter
	var entry CMSEntry
	var outputDir string
	var err error

	BeforeEach(func() {
		entry = CMSEntry{
			RawRow: []string{"A", "Slice", "Of", "Strings"},
		}

		outputDir, err = ioutil.TempDir("/tmp", "gogen")
		Expect(err).ToNot(HaveOccurred())
		writer = NewCMSWriter(path.Join(outputDir, "written_csv.csv"))
	})

	It("Writes CMS Entries to a csv file", func() {
		writer.WriteEntry(entry)
		writer.Flush()

		pathToExpectedResults, err := path.Abs(path.Join("..", "test_fixtures", "data", "expected_cms_writer_output.csv"))
		Expect(err).ToNot(HaveOccurred())
		expectedResultsBody, err := ioutil.ReadFile(pathToExpectedResults)
		Expect(err).ToNot(HaveOccurred())

		pathToOutput, err := path.Abs(path.Join(outputDir, "written_csv.csv"))
		Expect(err).ToNot(HaveOccurred())
		outputBody, err := ioutil.ReadFile(pathToOutput)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(outputBody)).To(Equal(string(expectedResultsBody)))
	})
})
