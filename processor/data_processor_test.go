package processor_test

import (
	"encoding/csv"
	//"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"gogen/data"
	. "gogen/processor"
	"io/ioutil"
	"os"
	path "path/filepath"
	"time"
)

var _ = Describe("DataProcessor", func() {
	var (
		outputDir     string
		dataProcessor DataProcessor
		err           error
	)

	Describe("Sacramento", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "sacramento", "cadoj_sacramento.csv"))
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, "SACRAMENTO")

			dojWriter := NewDOJWriter(path.Join(outputDir, "doj_sacramento_results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "doj_sacramento_results_condensed.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("SACRAMENTO")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "doj_sacramento_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			pathToExpectedDOJResults, err := path.Abs(path.Join("..", "test_fixtures", "sacramento", "doj_sacramento_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			for i, row := range outputDOJCSV {
				//fmt.Printf("output file %#v", outputDOJCSV)
				//for j, item := range row {
				//	Expect(item).To(Equal(expectedDOJResultsCSV[i][j]))
				//}
				Expect(row).To(Equal(expectedDOJResultsCSV[i]))
			}

			//Expect(outputDOJCSV).To(Equal(expectedDOJResultsCSV))
		})
	})

	Describe("San Joaquin", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "san_joaquin", "cadoj_san_joaquin.csv"))
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, "SAN JOAQUIN")

			dojWriter := NewDOJWriter(path.Join(outputDir, "doj_san_joaquin_results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "doj_san_joaquin_results_condensed.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("SAN JOAQUIN")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "doj_san_joaquin_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			pathToExpectedDOJResults, err := path.Abs(path.Join("..", "test_fixtures", "san_joaquin", "doj_san_joaquin_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			for i, row := range outputDOJCSV {
				//fmt.Printf("output file %#v", outputDOJCSV)
				//for j, item := range row {
				//	Expect(item).To(Equal(expectedDOJResultsCSV[i][j]))
				//}
				Expect(row).To(Equal(expectedDOJResultsCSV[i]))
			}

			//Expect(outputDOJCSV).To(Equal(expectedDOJResultsCSV))
		})
	})

	Describe("Contra Costa", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "contra_costa", "cadoj_contra_costa.csv"))
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA")

			dojResultsPath := path.Join(outputDir, "doj_contra_costa_results.csv")
			dojCondensedResultsPath := path.Join(outputDir, "doj_contra_costa_results_condensed.csv")

			dojWriter := NewDOJWriter(dojResultsPath)
			dojCondensedWriter := NewCondensedDOJWriter(dojCondensedResultsPath)

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter)
		})

		It("runs and has full output", func() {
			dataProcessor.Process("CONTRA COSTA")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "doj_contra_costa_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())

			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			pathToExpectedDOJResults, err := path.Abs(path.Join("..", "test_fixtures", "contra_costa", "doj_contra_costa_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			for i, row := range outputDOJCSV {
				//fmt.Printf("output file %#v", outputDOJCSV)
				//for j, item := range row {
				//	Expect(item).To(Equal(expectedDOJResultsCSV[i][j]))
				//}
				Expect(row).To(Equal(expectedDOJResultsCSV[i]))
			}

			//Expect(outputDOJCSV).To(Equal(expectedDOJResultsCSV))
		})
	})

	Describe("Los Angeles", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "los_angeles", "cadoj_los_angeles.csv"))
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, "LOS ANGELES")

			dojWriter := NewDOJWriter(path.Join(outputDir, "doj_los_angeles_results.csv"))
			dojCondensedWriter := NewDOJWriter(path.Join(outputDir, "doj_los_angeles_results_condensed.csv"))

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter)
		})

		It("runs and has output", func() {
			dataProcessor.Process("LOS ANGELES")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "doj_los_angeles_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			pathToExpectedDOJResults, err := path.Abs(path.Join("..", "test_fixtures", "los_angeles", "doj_los_angeles_results.csv"))
			Expect(err).ToNot(HaveOccurred())
			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			for i, row := range outputDOJCSV {
				//fmt.Printf("output file %#v", outputDOJCSV)
				//for j, item := range row {
				//	Expect(item).To(Equal(expectedDOJResultsCSV[i][j]))
				//}
				Expect(row).To(Equal(expectedDOJResultsCSV[i]))
			}

			//Expect(outputDOJCSV).To(Equal(expectedDOJResultsCSV))
		})
	})

	Describe("Condensed output file", func() {
		BeforeEach(func() {
			outputDir, err = ioutil.TempDir("/tmp", "gogen")
			Expect(err).ToNot(HaveOccurred())

			pathToDOJ, err := path.Abs(path.Join("..", "test_fixtures", "contra_costa", "cadoj_contra_costa_condensed_input.csv"))
			Expect(err).ToNot(HaveOccurred())

			comparisonTime := time.Date(2019, time.November, 11, 0, 0, 0, 0, time.UTC)

			dojInformation, _ := data.NewDOJInformation(pathToDOJ, comparisonTime, "CONTRA COSTA")

			dojResultsPath := path.Join(outputDir, "doj_contra_costa_results.csv")
			dojCondensedResultsPath := path.Join(outputDir, "doj_contra_costa_results_condensed.csv")

			dojWriter := NewDOJWriter(dojResultsPath)
			dojCondensedWriter := NewCondensedDOJWriter(dojCondensedResultsPath)

			dataProcessor = NewDataProcessor(dojInformation, dojWriter, dojCondensedWriter)
		})

		It("runs and has condensed output", func() {
			dataProcessor.Process("CONTRA COSTA")
			format.TruncatedDiff = false

			pathToDOJOutput, err := path.Abs(path.Join(outputDir, "doj_contra_costa_results_condensed.csv"))
			Expect(err).ToNot(HaveOccurred())
			OutputDOJFile, err := os.Open(pathToDOJOutput)
			Expect(err).ToNot(HaveOccurred())
			outputDOJCSV, err := csv.NewReader(OutputDOJFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			pathToExpectedDOJResults, err := path.Abs(path.Join("..", "test_fixtures", "contra_costa", "doj_contra_costa_results_condensed.csv"))
			Expect(err).ToNot(HaveOccurred())
			ExpectedDOJResultsFile, err := os.Open(pathToExpectedDOJResults)
			Expect(err).ToNot(HaveOccurred())
			expectedDOJResultsCSV, err := csv.NewReader(ExpectedDOJResultsFile).ReadAll()
			Expect(err).ToNot(HaveOccurred())

			for i, row := range outputDOJCSV {
				Expect(row).To(Equal(expectedDOJResultsCSV[i]))
			}
		})
	})

})
