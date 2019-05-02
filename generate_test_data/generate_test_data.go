package main

import (
	"github.com/jessevdk/go-flags"
	. "gogen/data"
	. "gogen/processor"
	"math/rand"
	"path/filepath"
)

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	TargetSize   int    `long:"target-size" description:"Desired number of lines in the output file" required:"true"`
	County		 string `long:"county" description:"Desired county" required:"true"`
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
	var rows [][]string

	personalInfos := generatePersonalInfos()

	randomNumber := rand.Float64()
	numberOfProp64Convictions := 1
	if randomNumber > 0.95 {
		numberOfProp64Convictions = 3
	} else if randomNumber > 0.8 {
		numberOfProp64Convictions = 2
	}

	for i := 0; i < numberOfProp64Convictions; i++ {
		rows = append(rows, generateProp64Row(personalInfos))
	}

	return rows
}

func generateProp64Row(infos PersonalInfos) []string {
	numberOfColumns := len(DojFullHeaders)
	row := make([]string, numberOfColumns)

	row[SUBJECT_ID] = infos.SubjectId
	row[CII_NUMBER] = infos.CII
	row[PRI_NAME] = infos.Name
	row[PRI_DOB] = infos.DOB
	row[PRI_SSN] = infos.SSN
	row[PRI_CDL] = infos.CDL
	row[CYC_DATE] = "20040424"
	row[STP_EVENT_DATE] = "20050621"
	row[STP_TYPE_DESCR] = "COURT"
	row[STP_ORI_CNTY_NAME] = opts.County
	row[CNT_ORDER] = "001002004000"
	row[OFN] = "1234K5n3"
	row[OFFENSE_DESCR] = "11357 HS-POSSESSION OF MARIJUANA"
	row[OFFENSE_TOC] = "F"
	row[FE_NUM_CRT_CASE] = "2342H8a8J"
	row[DISP_DESCR] = "CONVICTED"
	row[CONV_STAT_DESCR] = "M"
	row[SENT_LENGTH] = "45"
	row[SENT_TIME_CODE] = "D"
	row[COMMENT_TEXT] = "THIS IS A COMMENT"
	return row
}

func generatePersonalInfos() PersonalInfos {
	return PersonalInfos{
		SubjectId: "10",
		CII: "cii",
		Name: "SMITH,JOHN",
		DOB: "19790620",
		SSN: "123456789",
		CDL: "B6320998",
	}
}

type PersonalInfos struct {
	SubjectId string
	CII       string
	Name      string
	DOB       string
	SSN       string
	CDL       string
}
