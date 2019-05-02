package main

import (
	"github.com/jessevdk/go-flags"
	. "gogen/processor"
	. "gogen/data"
	"path/filepath"
)

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	TargetSize   int    `long:"target-size" description:"Desired number of lines in the output file" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	testWriter := NewWriter(filepath.Join(opts.OutputFolder, "generated_test_data.csv"), DojFullHeaders)

	totalRows := 0

	for totalRows < opts.TargetSize {
		rows := generateHistory()
		totalRows += len(rows)

		for _, row := range rows {
			testWriter.Write(row)
		}
	}

	testWriter.Flush()
}

func generateHistory() [][]string {
	numberOfColumns := len(DojFullHeaders)
	row := make([]string, numberOfColumns)

	row[SUBJECT_ID] = "10"
	row[CII_NUMBER] = "cii"
	row[PRI_NAME] = "SMITH,JOHN"
	row[PRI_DOB] = "19790620"
	row[PRI_SSN] = "123456789"
	row[PRI_CDL] = "B6320998"
	row[CYC_DATE] = "20040424"
	row[STP_EVENT_DATE] = "20050621"
	row[STP_TYPE_DESCR] = "COURT"
	row[STP_ORI_CNTY_NAME] = "ALAMEDA"
	row[CNT_ORDER] = "001002004000"
	row[OFN] = "1234K5n3"
	row[OFFENSE_DESCR] = "POSSESSION OF MARIJUANA"
	row[OFFENSE_TOC] = "F"
	row[FE_NUM_CRT_CASE] = "2342H8a8J"
	row[DISP_DESCR] = "CONVICTED"
	row[CONV_STAT_DESCR] = "M"
	row[SENT_LENGTH] = "45"
	row[SENT_TIME_CODE] = "D"
	row[COMMENT_TEXT] = "THIS IS A COMMENT"

	return [][]string{row}
}
