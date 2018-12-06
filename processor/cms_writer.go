package processor

import (
	"encoding/csv"
	"gogen/data"
	"os"
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
	WriteEntry(data.CMSEntry, *data.DOJHistory, EligibilityInfo)
	Flush()
}

type csvWriter struct {
	outputFileWriter *csv.Writer
}

func NewCMSWriter(outputFilePath string) CMSWriter {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}

	w := csvWriter{
		outputFileWriter: csv.NewWriter(outputFile),
	}

	w.outputFileWriter.Write(headers)

	return w
}

func (cw csvWriter) WriteEntry(entry data.CMSEntry, history *data.DOJHistory, info EligibilityInfo) {
	var historyCols []string
	if history == nil {
		historyCols = []string{"no match", "no match", "no match", "no match", "no match"}
	} else {
		historyCols = []string{
			history.Name,
			history.DOB.Format("2006-01-02"),
			history.SubjectID,
			history.OriginalCII,
			history.SSN,
		}
	}

	eligibilityCols := []string{
		info.Superstrikes,
		info.SuperstrikeCodeSections,
		info.PC290Charges,
		info.PC290CodeSections,
		info.PC290Registration,
		info.TwoPriors,
		info.Over1Lb,
		info.QFinalSum,
		info.AgeAtConviction,
		info.YearsSinceEvent,
		info.YearsSinceMostRecentConviction,
		info.FinalRecommendation,
	}

	extraCols := append(historyCols, eligibilityCols...)

	_ = cw.outputFileWriter.Write(append(entry.RawRow, extraCols...))
}

func (cw csvWriter) Flush() {
	cw.outputFileWriter.Flush()
}
