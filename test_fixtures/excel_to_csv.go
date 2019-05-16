package test_fixtures

import (
	"encoding/csv"
	"fmt"
	"github.com/tealeg/xlsx"
	. "gogen/processor"
	"io/ioutil"
	path "path/filepath"
	"strconv"
)

func ExtractCSVFixtures(inputPathString string) (string, string, error) {
	xlsxPath, err := path.Abs(path.Join(inputPathString))

	inputCSV := extractInputCSV(xlsxPath)
	expectedResultsCSV := extractFullResultsCSV(xlsxPath)

	return inputCSV, expectedResultsCSV, err
}

func extractInputCSV(xlsxPath string) string {
	tmpCSVfile, err := ioutil.TempFile("", "cadoj_file.csv")
	xlFile, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		fmt.Println(err)
	}
	w := csv.NewWriter(tmpCSVfile)
	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			var rowSlice []string
			for j, cell := range row.Cells {
				if j >= len(DojFullHeaders) {
					break
				}
				text, err := cell.FormattedValue()
				if err != nil {
					rowSlice = append(rowSlice, err.Error())
				}
				rowSlice = append(rowSlice, text)
			}
			if rowIndex == 0 {
				continue
			} else {
				w.Write(rowSlice)
			}
		}
	}
	w.Flush()
	return tmpCSVfile.Name()
}

func extractFullResultsCSV(xlsxPath string) string {
	tmpCSVfile, err := ioutil.TempFile("", "cadoj_file.csv")
	xlFile, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		fmt.Println(err)
	}
	w := csv.NewWriter(tmpCSVfile)
	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			var rowSlice []string
			for _, cell := range row.Cells {
				text, err := cell.FormattedValue()
				if err != nil {
					rowSlice = append(rowSlice, err.Error())
				}
				rowSlice = append(rowSlice, text)
			}
			if rowIndex == 0 {
				continue
			} else {
				w.Write(rowSlice)
			}
		}
	}
	w.Flush()
	return tmpCSVfile.Name()
}



func ExtractCondensedCSVFixture(inputPathString string) (string, error) {
	xlsxPath, err := path.Abs(path.Join(inputPathString))

	expectedCondensedResultsCSV := extractCondensedResultsCSV(xlsxPath)

	return expectedCondensedResultsCSV, err
}



func extractCondensedResultsCSV(xlsxPath string) string {
	tmpCSVfile, err := ioutil.TempFile("", "cadoj_condensed_file.csv")
	xlFile, err := xlsx.OpenFile(xlsxPath)
	if err != nil {
		fmt.Println(err)
	}
	condensedCSV := csv.NewWriter(tmpCSVfile)
	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			var rowSlice []string
			for cellIndex, cell := range row.Cells {
				if cellIndex >= (len(DojFullHeaders) + len(EligiblityHeaders)) {
					break
				}
				inCondensedOutput, _ := strconv.ParseBool(sheet.Rows[0].Cells[cellIndex].String())
				if inCondensedOutput == true {
					text, err := cell.FormattedValue()
					if err != nil {
						rowSlice = append(rowSlice, err.Error())
					}
					rowSlice = append(rowSlice, text)
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
