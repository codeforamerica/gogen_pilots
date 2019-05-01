package processor

import (
	"encoding/csv"
	"fmt"
	"gogen/data"
	"os"
	"time"
)

var eligiblityHeaders = []string{
	"Case Number",
	"# of convictions on record",
	"Superstrike Code Section(s)",
	"PC290 Code Section(s)",
	"PC290 Registration",
	"Date of Conviction",
	"Years Since This Conviction",
	"Years Since Any Conviction",
	"# of Prop 64 convictions",
	"Deceased",
	"Eligibility Determination",
	"Eligibility Reason",
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
	"CII_NUMBER",//11
	"PRI_NAME",//12
	"GENDER",//13
	"PRI_DOB",//14
	"PRI_SSN",
	"PRI_CDL",
	"PRI_IDN",
	"PRI_INN",
	"FBI_NUMBER",
	"PDR_SEG_SEP",
	"RACE_CODE",
	"RACE_DESCR",//22
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
	"CYC_DATE",//37
	"STP_SEG_SEP",
	"STP_ORDER",
	"STP_EVENT_DATE",//40
	"STP_TYPE_CODE",
	"STP_TYPE_DESCR",
	"STP_ORI_TYPE",
	"STP_ORI_TYPE_DESCR",
	"STP_ORI_CODE",
	"STP_ORI_DESCR",//46
	"STP_ORI_CNTY_CODE",
	"STP_ORI_CNTY_NAME",//48
	"CNT_SEG_SEP",
	"CNT_ORDER",
	"DISP_DATE",//51
	"OFN",//52
	"OFFENSE_CODE",
	"OFFENSE_DESCR",//54
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
	"DISP_DESCR",//80
	"CONV_STAT_CODE",
	"CONV_STAT_DESCR",//82
	"SENT_SEG_SEP",
	"SENT_ORDER",
	"SENT_LOC_CODE",
	"SENT_LOC_DESCR",//86
	"SENT_LENGTH",//87
	"SENT_TIME_CODE",//88
	"SENT_TIME_DESCR",
	"CYC_AGE",//90
	"CII_TYPE",
	"CII_TYPE_ALPHA",
	"COMMENT_TEXT",//93
	"END_OF_REC",//94
}
	var dojCondensedHeaders = []string{
	"CII_NUMBER",
	"PRI_NAME",
	"GENDER",
	"PRI_DOB",
	"RACE_DESCR",
	"CYC_DATE",
	"STP_EVENT_DATE",
	"STP_ORI_DESCR",
	"STP_ORI_CNTY_NAME",
	"DISP_DATE",
	"OFN",
	"OFFENSE_DESCR",
	"DISP_DESCR",
	"CONV_STAT_DESCR",
	"SENT_LOC_DESCR",
	"SENT_LENGTH",
	"SENT_TIME_CODE",
	"CYC_AGE",
	"COMMENT_TEXT",
	"END_OF_REC",
}

type DOJWriter interface {
	WriteDOJEntry([]string, *data.EligibilityInfo)
	Flush()
}

type csvWriter struct {
	outputFileWriter *csv.Writer
	filename         string
}

func NewWriter(outputFilePath string, headers []string) DOJWriter {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}

	w := new(csvWriter)
	w.outputFileWriter = csv.NewWriter(outputFile)
	w.filename = outputFilePath

	err = w.outputFileWriter.Write(headers)
	if err != nil {
		panic(err)
	}

	return w
}

func NewDOJWriter(outputFilePath string) DOJWriter {
	headers := append(dojFullHeaders, eligiblityHeaders...)
	return NewWriter(outputFilePath, headers)
}

func NewCondensedDOJWriter(outputFilePath string) DOJWriter {
	headers := append(dojCondensedHeaders, eligiblityHeaders...)
	return NewWriter(outputFilePath, headers)
}

func (cw csvWriter) WriteDOJEntry(entry []string, info *data.EligibilityInfo) {
	var eligibilityCols []string

	if info != nil {
		eligibilityCols = []string{
			info.CaseNumber,
			writeInt(info.NumberOfConvictionsOnRecord),
			info.Superstrikes,
			info.PC290CodeSections,
			info.PC290Registration,
			writeDate(info.DateOfConviction),
			writeFloat(info.YearsSinceThisConviction),
			writeFloat(info.YearsSinceMostRecentConviction),
			writeInt(info.NumberOfProp64Convictions),
			info.Deceased,
			info.EligibilityDetermination,
			info.EligibilityReason,
		}
	} else {
		eligibilityCols = make([]string, len(eligiblityHeaders))
	}

	_ = cw.outputFileWriter.Write(append(entry, eligibilityCols...))
}

func writeDate(val time.Time) string {
	return val.Format("01/02/2006")
}

func (cw csvWriter) Flush() {
	cw.outputFileWriter.Flush()
}

func writeFloat(val float64) string {
	return fmt.Sprintf("%.1f", val)
}

func writeInt(val int) string {
	return fmt.Sprintf("%d", val)
}
