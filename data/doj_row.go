package data

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DOJRow struct {
	Name                   string
	WeakName               string
	SubjectID              string
	CII                    string
	SSN                    string
	DOB                    time.Time
	CDL                    string
	Convicted              bool
	CodeSection            string
	DispositionDate        time.Time
	OFN                    string
	Type                   string
	PC290Registration      bool
	County                 string
	Felony                 bool
	NumCrtCase             string
	CycleDate              time.Time
	RawRow                 []string
	CourtNoParts           []string
	CountOrder             string
	Index                  int
	SentenceEndDate        time.Time
	SentencePartDuration   time.Duration
	HasProp64ChargeInCycle bool
}

const dateFormat = "20060102"

func NewDOJRow(rawRow []string, index int) DOJRow {

	var isFelony bool

	if rawRow[CONV_STAT_DESCR] == "" {
		isFelony = rawRow[OFFENSE_TOC] == "F"
	} else {
		isFelony = rawRow[CONV_STAT_DESCR] == "FELONY"
	}

	return DOJRow{
		Name:                 rawRow[PRI_NAME],
		WeakName:             strings.Split(rawRow[PRI_NAME], " ")[0],
		SubjectID:            rawRow[SUBJECT_ID],
		CII:                  rawRow[CII_NUMBER],
		SSN:                  rawRow[PRI_SSN],
		DOB:                  parseDate(dateFormat, rawRow[PRI_DOB]),
		CDL:                  rawRow[PRI_CDL],
		Convicted:            strings.HasPrefix(rawRow[DISP_DESCR], "CONVICTED"),
		CodeSection:          findCodeSection(rawRow),
		DispositionDate:      parseDate(dateFormat, rawRow[STP_EVENT_DATE]),
		OFN:                  rawRow[OFN],
		NumCrtCase:           rawRow[FE_NUM_CRT_CASE],
		CycleDate:            parseDate(dateFormat, rawRow[CYC_DATE]),
		Type:                 rawRow[STP_TYPE_DESCR],
		PC290Registration:    rawRow[STP_TYPE_DESCR] == "REGISTRATION" && strings.HasPrefix(rawRow[OFFENSE_DESCR], "290"),
		County:               rawRow[STP_ORI_CNTY_NAME],
		Felony:               isFelony,
		CountOrder:           rawRow[CNT_ORDER],
		Index:                index,
		SentenceEndDate:      getSentenceEndDate(rawRow),
		SentencePartDuration: getSentencePartDuration(rawRow),
	}
}

func getSentenceEndDate(rawRow []string) time.Time {
	dispDate := parseDate(dateFormat, rawRow[STP_EVENT_DATE])
	return dispDate.Add(getSentencePartDuration(rawRow))
}

func getSentencePartDuration(rawRow []string) time.Duration {
	sentenceLength, _ := strconv.Atoi(rawRow[SENT_LENGTH])

	days := time.Duration(24) * (time.Hour)
	years := time.Date(2012, 03, 04, 0, 0, 0, 0, time.UTC).Sub(time.Date(2011, 03, 04, 0, 0, 0, 0, time.UTC))
	months := years / 12

	switch rawRow[SENT_TIME_CODE] {
	case "D":
		return time.Duration(sentenceLength) * days
	case "M":
		return time.Duration(sentenceLength) * months
	case "Y":
		return time.Duration(sentenceLength) * years
	}
	return time.Duration(0)
}

func findCodeSection(rawRow []string) string {
	switch offenseDescription := rawRow[OFFENSE_DESCR]; offenseDescription {
	case "SEE COMMENT FOR CHARGE":
		return strings.Split(rawRow[COMMENT_TEXT], "-")[0]
	default:
		return strings.Split(rawRow[OFFENSE_DESCR], "-")[0]
	}
}

func (row *DOJRow) MatchingCourtNumber(courtNumber string) bool {
	if courtNumber == row.OFN || courtNumber == row.NumCrtCase {
		return true
	}

	if row.CourtNoParts == nil {
		row.CourtNoParts = regexp.MustCompile("[ ,-]").Split(row.OFN, -1)
	}

	for _, part := range row.CourtNoParts {
		if part == courtNumber {
			return true
		}
	}
	return false
}

func (row *DOJRow) OccurredInLast7Years() bool {
	sevenYearsAgo := time.Now().AddDate(-7,0,0)

	if row.DispositionDate.After(sevenYearsAgo) {
		return true
	} else {
		return false
	}
}

const (
	RECORD_ID = iota
	SUBJECT_STATUS
	SUBJECT_ID
	REQ_SEG_SEP
	REQ_CII_NUMBER
	REQ_NAME
	REQ_GENDER
	REQ_DOB
	REQ_CDL
	REQ_SSN
	PII_SEG_SEP
	CII_NUMBER
	PRI_NAME
	GENDER
	PRI_DOB
	PRI_SSN
	PRI_CDL
	PRI_IDN
	PRI_INN
	FBI_NUMBER
	PDR_SEG_SEP
	RACE_CODE
	RACE_DESCR
	EYE_COLOR_CODE
	EYE_COLOR_DESCR
	HAIR_COLOR_CODE
	HAIR_COLOR_DESCR
	HEIGHT
	WEIGHT
	SINGLE_SOURCE
	MULTI_SOURCE
	POB_CODE
	POB_NAME
	POB_TYPE
	CITIZENSHIP_LIST
	CYC_SEG_SEP
	CYC_ORDER
	CYC_DATE
	STP_SEG_SEP
	STP_ORDER
	STP_EVENT_DATE
	STP_TYPE_CODE
	STP_TYPE_DESCR
	STP_ORI_TYPE
	STP_ORI_TYPE_DESCR
	STP_ORI_CODE
	STP_ORI_DESCR
	STP_ORI_CNTY_CODE
	STP_ORI_CNTY_NAME
	CNT_SEG_SEP
	CNT_ORDER
	DISP_DATE
	OFN
	OFFENSE_CODE
	OFFENSE_DESCR
	OFFENSE_TOC
	OFFENSE_QUAL_LST
	DISP_OFFENSE_CODE
	DISP_OFFENSE_DESCR
	DISP_OFFENSE_TOC
	DISP_OFFENSE_QUAL_LST
	CONV_OFFENSE_ORDER
	CONV_OFFENSE_CODE
	CONV_OFFENSE_DESCR
	CONV_OFFENSE_TOC
	CONV_OFFENSE_QUAL_LST
	FE_NUM_ORDER
	FE_NUM_ARR_AGY
	FE_NUM_BNCH_WARR
	FE_NUM_CITE
	FE_NUM_DOCKET
	FE_NUM_INCIDENT
	FE_NUM_BOOKING
	FE_NUM_NUMBER
	FE_NUM_REMAND
	FE_NUM_OOS_INN
	FE_NUM_CRT_CASE
	FE_NUM_WARRANT
	DISP_ORDER
	DISP_CODE
	DISP_DESCR
	CONV_STAT_CODE
	CONV_STAT_DESCR
	SENT_SEG_SEP
	SENT_ORDER
	SENT_LOC_CODE
	SENT_LOC_DESCR
	SENT_LENGTH
	SENT_TIME_CODE
	SENT_TIME_DESCR
	CYC_AGE
	CII_TYPE
	CII_TYPE_ALPHA
	COMMENT_TEXT
	END_OF_REC
)
