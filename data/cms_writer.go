package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

var headers = []string{
	"Court Number",
	"Ind",
	"Incident Number",
	"Truename",
	"Case Level",
	"Case Dispo",
	"Case Disposition Description ",
	"Dispo Date",
	"Action Number",
	"1st Filed",
	"Charge Level",
	"Charge Date",
	"Current Charge",
	"Current Level",
	"Current Charge Description",
	"Chg Disp",
	"Charge Disposition Description",
	"Ch Dispo Date",
	"Race",
	"Sex",
	"DOB",
	"SFNO",
	"CII",
	"FBI",
	"SSN",
	"DL Number",
	"EOR",
	"PRI_NAME",
	"PRI_DOB",
	"SUBJECT_ID",
	"CII_NUMBER",
	"PRI_SSN",
	"Superstrikes",
	"Superstrike Code Section(s)",
	"PC290 Charges",
	"PC290 Code Section(s)",
	"PC290 Registration",
	"Two Priors",
	"Over 1lb",
	"Q_final_sum",
	"Age at Conviction",
	"Years Since Event",
	"Years Since Most Recent Conviction",
	"Final Recommendation",
}

type CMSWriter interface {
	WriteEntry(CMSEntry, EligibilityInfo)
	Flush()
}

type csvWriter struct {
	outputFileWriter *csv.Writer
}

func NewCMSWriter(outputFilePath string) CMSWriter {
	outputFile, err := os.Create(filepath.Join(outputFilePath))
	if err != nil {
		panic(err)
	}

	w := csvWriter{
		outputFileWriter: csv.NewWriter(outputFile),
	}

	w.outputFileWriter.Write(headers)

	return w
}

func (cw csvWriter) WriteEntry(entry CMSEntry, info EligibilityInfo) {
	cw.outputFileWriter.Write(append(entry.RawRow, info.Over1Lb, fmt.Sprintf("%.1f", info.QFinalSum)))
}

func (cw csvWriter) Flush() {
	cw.outputFileWriter.Flush()
}
