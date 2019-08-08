package exporter

import (
	"encoding/csv"
	"fmt"
	"gogen/data"
	"os"
	"time"
)

var EligiblityHeaders = []string{
	"Case Number",
	"# of convictions on record",
	"Superstrike Code Section(s)",
	"PC290 Code Section(s)",
	"PC290 Registration",
	"Date of Conviction",
	"Years Since This Conviction",
	"Years Since Any Conviction",
	"# of Prop 64 convictions",
	"# of HS 11357 convictions",
	"# of HS 11358 convictions",
	"# of HS 11359 convictions",
	"# of HS 11360 convictions",
	"Deceased",
	"Eligibility Determination",
	"Eligibility Reason",
}

var DojFullHeaders = []string{
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
}
var DojCondensedHeaders = []string{
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
	WriteEntryWithEligibilityInfo([]string, *data.EligibilityInfo)
	WriteCondensedEntryWithEligibilityInfo([]string, *data.EligibilityInfo)
	Write([]string)
	Flush()
}

type csvWriter struct {
	outputFileWriter *csv.Writer
	filename         string
}

func NewWriter(outputFilePath string, headers []string) (DOJWriter, error) {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return nil, err
	}

	w := new(csvWriter)
	w.outputFileWriter = csv.NewWriter(outputFile)
	w.filename = outputFilePath

	err = w.outputFileWriter.Write(headers)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func NewDOJWriter(outputFilePath string) (DOJWriter, error) {
	headers := append(DojFullHeaders, EligiblityHeaders...)
	return NewWriter(outputFilePath, headers)
}

func NewCondensedDOJWriter(outputFilePath string) (DOJWriter, error) {
	headers := append(DojCondensedHeaders, EligiblityHeaders...)
	return NewWriter(outputFilePath, headers)
}

func (cw csvWriter) WriteEntryWithEligibilityInfo(entry []string, info *data.EligibilityInfo) {
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
			writeInt(info.NumberOf11357Convictions),
			writeInt(info.NumberOf11358Convictions),
			writeInt(info.NumberOf11359Convictions),
			writeInt(info.NumberOf11360Convictions),
			info.Deceased,
			info.EligibilityDetermination,
			info.EligibilityReason,
		}
	} else {
		eligibilityCols = make([]string, len(EligiblityHeaders))
	}

	cw.Write(append(entry, eligibilityCols...))
}

func (cw csvWriter) WriteCondensedEntryWithEligibilityInfo(entry []string, info *data.EligibilityInfo) {
	var condensedRow []string

	includedColumns := []int{
		data.CII_NUMBER,
		data.PRI_NAME,
		data.GENDER,
		data.PRI_DOB,
		data.RACE_DESCR,
		data.CYC_DATE,
		data.STP_EVENT_DATE,
		data.STP_ORI_DESCR,
		data.STP_ORI_CNTY_NAME,
		data.DISP_DATE,
		data.OFN,
		data.OFFENSE_DESCR,
		data.DISP_DESCR,
		data.CONV_STAT_DESCR,
		data.SENT_LOC_DESCR,
		data.SENT_LENGTH,
		data.SENT_TIME_CODE,
		data.CYC_AGE,
		data.COMMENT_TEXT,
		data.END_OF_REC,
	}

	for _, col := range includedColumns {
		condensedRow = append(condensedRow, entry[col])
	}

	cw.WriteEntryWithEligibilityInfo(condensedRow, info)
}

func writeDate(val time.Time) string {
	return val.Format("01/02/2006")
}

func (cw csvWriter) Flush() {
	cw.outputFileWriter.Flush()
}

func (cw csvWriter) Write(line []string) {
	_ = cw.outputFileWriter.Write(line)
}

func writeFloat(val float64) string {
	return fmt.Sprintf("%.1f", val)
}

func writeInt(val int) string {
	return fmt.Sprintf("%d", val)
}
