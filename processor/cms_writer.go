package processor

import (
	"encoding/csv"
	"gogen/data"
	"os"
)

var cmsHeaders = []string{
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
	"Booked charge",
	"Booked charge level",
	"Booked charge date",
	"Race",
	"Sex",
	"DOB",
	"SFNO",
	"CII",
	"FBI",
	"SSN",
	"DL Number",
	"EOR",
}

var dojHistoryHeaders = []string{
	"PRI_NAME",
	"PRI_DOB",
	"SUBJECT_ID",
	"CII_NUMBER",
	"PRI_SSN",
}

var eligiblityHeaders = []string{
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

var dojFullHeaders = []string{
	"RECORD_ID",
	"SUBJECT_STATUS",
	"SUBJECT_ID",
	"REQ_SEG_SEP",
	"REQ_CII_NUMBER",
	"REQ_NAME",
	"REQ_GENDER",
	"REQ_DOB",
	"REQ_CDL",
	"REQ_SSN",
	"PII_SEG_SEP",
	"CII_NUMBER",
	"PRI_NAME",
	"GENDER",
	"PRI_DOB",
	"PRI_SSN",
	"PRI_CDL",
	"PRI_IDN",
	"PRI_INN",
	"FBI_NUMBER",
	"PDR_SEG_SEP",
	"RACE_CODE",
	"RACE_DESCR",
	"EYE_COLOR_CODE",
	"EYE_COLOR_DESCR",
	"HAIR_COLOR_CODE",
	"HAIR_COLOR_DESCR",
	"HEIGHT",
	"WEIGHT",
	"SINGLE_SOURCE",
	"MULTI_SOURCE",
	"POB_CODE",
	"POB_NAME",
	"POB_TYPE",
	"CITIZENSHIP_LIST",
	"CYC_SEG_SEP",
	"CYC_ORDER",
	"CYC_DATE",
	"STP_SEG_SEP",
	"STP_ORDER",
	"STP_EVENT_DATE",
	"STP_TYPE_CODE",
	"STP_TYPE_DESCR",
	"STP_ORI_TYPE",
	"STP_ORI_TYPE_DESCR",
	"STP_ORI_CODE",
	"STP_ORI_DESCR",
	"STP_ORI_CNTY_CODE",
	"STP_ORI_CNTY_NAME",
	"CNT_SEG_SEP",
	"CNT_ORDER",
	"DISP_DATE",
	"OFN",
	"OFFENSE_CODE",
	"OFFENSE_DESCR",
	"OFFENSE_TOC",
	"OFFENSE_QUAL_LST",
	"DISP_OFFENSE_CODE",
	"DISP_OFFENSE_DESCR",
	"DISP_OFFENSE_TOC",
	"DISP_OFFENSE_QUAL_LST",
	"CONV_OFFENSE_ORDER",
	"CONV_OFFENSE_CODE",
	"CONV_OFFENSE_DESCR",
	"CONV_OFFENSE_TOC",
	"CONV_OFFENSE_QUAL_LST",
	"FE_NUM_ORDER",
	"FE_NUM_ARR_AGY",
	"FE_NUM_BNCH_WARR",
	"FE_NUM_CITE",
	"FE_NUM_DOCKET",
	"FE_NUM_INCIDENT",
	"FE_NUM_BOOKING",
	"FE_NUM_NUMBER",
	"FE_NUM_REMAND",
	"FE_NUM_OOS_INN",
	"FE_NUM_CRT_CASE",
	"FE_NUM_WARRANT",
	"DISP_ORDER",
	"DISP_CODE",
	"DISP_DESCR",
	"CONV_STAT_CODE",
	"CONV_STAT_DESCR",
	"SENT_SEG_SEP",
	"SENT_ORDER",
	"SENT_LOC_CODE",
	"SENT_LOC_DESCR",
	"SENT_LENGTH",
	"SENT_TIME_CODE",
	"SENT_TIME_DESCR",
	"CYC_AGE",
	"CII_TYPE",
	"CII_TYPE_ALPHA",
	"COMMENT_TEXT",
	"END_OF_REC",
	"",
}

type CMSWriter interface {
	WriteEntry(data.CMSEntry, *data.DOJHistory, EligibilityInfo)
	WriteDOJEntry([]string, EligibilityInfo)
	Flush()
}

type csvWriter struct {
	outputFileWriter *csv.Writer
	filename string
}

func NewCMSWriter(outputFilePath string) CMSWriter {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}

	w := new(csvWriter)
	w.outputFileWriter = csv.NewWriter(outputFile)
	w.filename = outputFilePath

	headers := append(cmsHeaders, dojHistoryHeaders...)
	headers = append(headers, eligiblityHeaders...)

	w.outputFileWriter.Write(headers)

	return w
}

func NewDOJWriter(outputFilePath string) CMSWriter {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}

	w := new(csvWriter)
	w.outputFileWriter = csv.NewWriter(outputFile)
	w.filename = outputFilePath

	headers := append(dojFullHeaders, eligiblityHeaders...)

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

func (cw csvWriter) WriteDOJEntry(entry []string, info EligibilityInfo) {
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

	_ = cw.outputFileWriter.Write(append(entry, eligibilityCols...))
}

func (cw csvWriter) Flush() {
	cw.outputFileWriter.Flush()
}
