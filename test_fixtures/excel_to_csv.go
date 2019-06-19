package test_fixtures

import (
	"encoding/csv"
	"fmt"
	"gogen/exporter"
	"io/ioutil"
	"os"
	path "path/filepath"
	"strconv"

	"github.com/tealeg/xlsx"
)

func ExportFullCSVFixtures(inputPathString string, outputPathString string) (string, string, error) {
	xlsxPath, err := path.Abs(path.Join(inputPathString))

	inputCSV := writeInputCSV(xlsxPath, outputPathString)
	expectedResultsCSV := writeFullResultsCSV(xlsxPath, outputPathString)

	return inputCSV, expectedResultsCSV, err
}

func ExtractFullCSVFixtures(inputPathString string) (string, string, error) {
	xlsxPath, err := path.Abs(path.Join(inputPathString))

	inputCSV := writeInputCSV(xlsxPath, "")
	expectedResultsCSV := writeFullResultsCSV(xlsxPath, "")

	return inputCSV, expectedResultsCSV, err
}

func ExtractCondensedCSVFixture(inputPathString string) (string, error) {
	xlsxPath, err := path.Abs(path.Join(inputPathString))

	expectedCondensedResultsCSV := writeCondensedColumnsCSV(xlsxPath)

	return expectedCondensedResultsCSV, err
}

func ExtractProp64ConvictionsCSVFixture(inputPathString string) (string, error) {
	xlsxPath, err := path.Abs(path.Join(inputPathString))

	expectedCondensedResultsCSV := writeProp64RowsCSV(xlsxPath)

	return expectedCondensedResultsCSV, err
}

func writeInputCSV(xlsxPath string, outputPath string) string {
	var (
		tmpCSVfile *os.File
		excelFile  *xlsx.File
	)

	if outputPath != "" {
		tmpCSVfile, excelFile = createExportFile(xlsxPath, path.Join(outputPath, "input.csv"))
	} else {
		tmpCSVfile, excelFile = createTempFile(xlsxPath)
	}

	inputCSV := csv.NewWriter(tmpCSVfile)
	for _, sheet := range excelFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			var rowSlice []string
			for cellIndex, cell := range row.Cells {
				lengthOfInputFile := len(exporter.DojFullHeaders) + 1
				if cellIndex >= lengthOfInputFile {
					break
				}
				firstCell := sheet.Rows[rowIndex].Cells[0]
				if cell == firstCell {
					continue
				} else {
					rowSlice = createRowSlice(cell, rowSlice)
				}
			}
			if rowIndex == 0 {
				continue
			} else {
				inputCSV.Write(rowSlice)
			}
		}
	}
	inputCSV.Flush()
	return tmpCSVfile.Name()
}

func writeFullResultsCSV(xlsxPath string, outputPath string) string {
	var (
		tmpCSVfile *os.File
		excelFile  *xlsx.File
	)

	if outputPath != "" {
		tmpCSVfile, excelFile = createExportFile(xlsxPath, path.Join(outputPath, "full_results.csv"))
	} else {
		tmpCSVfile, excelFile = createTempFile(xlsxPath)
	}

	fullResultsCSV := csv.NewWriter(tmpCSVfile)
	for _, sheet := range excelFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			var rowSlice []string
			for _, cell := range row.Cells {
				firstCell := sheet.Rows[rowIndex].Cells[0]
				if cell == firstCell {
					continue
				} else {
					rowSlice = createRowSlice(cell, rowSlice)
				}
			}
			if rowIndex == 0 {
				continue
			} else {
				fullResultsCSV.Write(rowSlice)
			}
		}
	}
	fullResultsCSV.Flush()
	return tmpCSVfile.Name()
}

func writeCondensedColumnsCSV(xlsxPath string) string {
	tmpCSVfile, excelFile := createTempFile(xlsxPath)
	condensedCSV := csv.NewWriter(tmpCSVfile)
	for _, sheet := range excelFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			var rowSlice []string
			for cellIndex, cell := range row.Cells {
				addToOutput, _ := strconv.ParseBool(sheet.Rows[0].Cells[cellIndex].String())
				firstCell := sheet.Rows[rowIndex].Cells[0]
				if addToOutput == true {
					if cell == firstCell {
						continue
					} else {
						rowSlice = createRowSlice(cell, rowSlice)
					}
				}
			}
			if rowIndex == 0 {
				continue
			} else {
				condensedCSV.Write(rowSlice)
			}
		}
	}
	condensedCSV.Flush()
	return tmpCSVfile.Name()
}

func writeProp64RowsCSV(xlsxPath string) string {
	tmpCSVfile, excelFile := createTempFile(xlsxPath)
	condensedCSV := csv.NewWriter(tmpCSVfile)
	for _, sheet := range excelFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			addToOutput := false
			var rowSlice []string
			for cellIndex, cell := range row.Cells {
				addToOutput, _ = strconv.ParseBool(sheet.Rows[rowIndex].Cells[0].String())
				if addToOutput == true {
					if cellIndex == 0 {
						continue
					} else {
						rowSlice = createRowSlice(cell, rowSlice)
					}
				}
			}
			if addToOutput == true {
				condensedCSV.Write(rowSlice)
			}
		}
	}
	condensedCSV.Flush()
	return tmpCSVfile.Name()
}

func createRowSlice(cell *xlsx.Cell, rowSlice []string) []string {
	text, err := cell.FormattedValue()
	if err != nil {
		rowSlice = append(rowSlice, err.Error())
	}
	rowSlice = append(rowSlice, text)
	return rowSlice
}

func createTempFile(xlsxPath string) (*os.File, *xlsx.File) {
	tmpCSVfile, err := ioutil.TempFile("", "temp_file.csv")
	excelFile, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		fmt.Println(err)
	}
	return tmpCSVfile, excelFile
}

func createExportFile(xlsxPath, outputPath string) (*os.File, *xlsx.File) {
	exportCSVFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err)
	}
	excelFile, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		fmt.Println(err)
	}
	return exportCSVFile, excelFile
}
