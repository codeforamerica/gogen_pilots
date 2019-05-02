package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	. "gogen/data"
	. "gogen/processor"
	"math"
	"math/rand"
	"path/filepath"
	"strconv"
	"time"
)

var opts struct {
	OutputFolder string `long:"outputs" description:"The folder in which to place result files" required:"true"`
	TargetSize   int    `long:"target-size" description:"Desired number of lines in the output file" required:"true"`
	County       string `long:"county" description:"Desired county" required:"true"`
}

var counters struct {
	SubjectId int
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	testWriter := NewWriter(filepath.Join(opts.OutputFolder, "generated_test_data.csv"), DojFullHeaders)

	totalRows := 0

	counters.SubjectId = 10000000
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
	counters.SubjectId++

	return PersonalInfos{
		SubjectId: strconv.Itoa(counters.SubjectId),
		CII:       randomNDigitNumber(10),
		Name:      randomChoice(lastNames) + "," + randomChoice(firstNames),
		DOB:       randomDate(time.Date(1950, 1, 0, 0, 0, 0, 0, time.UTC), time.Date(2001, 1, 0, 0, 0, 0, 0, time.UTC)),
		SSN:       randomNDigitNumber(9),
		CDL:       randomLetter() + randomNDigitNumber(7),
	}
}

func randomDate(startDate time.Time, endDate time.Time) string {
	min := startDate.Unix()
	max := endDate.Unix()
	delta := max - min

	randomDateInSeconds := rand.Int63n(delta) + min
	randomDate := time.Unix(randomDateInSeconds, 0)
	return randomDate.Format("20060102")
}

func randomLetter() string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return string(alphabet[rand.Intn(len(alphabet))])
}

func randomNDigitNumber(n int) string {
	formatString := fmt.Sprintf("%%0%dd", n)
	return fmt.Sprintf(formatString, rand.Intn(int(math.Pow10(n))))
}

func randomChoice(values []string) string {
	return values[rand.Intn(len(values))]
}

type PersonalInfos struct {
	SubjectId string
	CII       string
	Name      string
	DOB       string
	SSN       string
	CDL       string
}

var firstNames = []string{
	"MICHAEL",
	"JENNIFER",
	"CHRISTOPHER",
	"AMY",
	"JASON",
	"MELISSA",
	"DAVID",
	"MICHELLE",
	"JAMES",
	"KIMBERLY",
	"JOHN",
	"LISA",
	"ROBERT",
	"ANGELA",
	"BRIAN",
	"HEATHER",
	"WILLIAM",
	"STEPHANIE",
	"MATTHEW",
	"NICOLE",
	"JOSEPH",
	"JESSICA",
	"DANIEL",
	"ELIZABETH",
	"KEVIN",
	"REBECCA",
	"ERIC",
	"KELLY",
	"JEFFREY",
	"MARY",
	"RICHARD",
	"CHRISTINA",
	"SCOTT",
	"AMANDA",
	"MARK",
	"JULIE",
	"STEVEN",
	"SARAH",
	"THOMAS",
	"LAURA",
	"TIMOTHY",
	"SHANNON",
	"ANTHONY",
	"CHRISTINE",
	"CHARLES",
	"TAMMY",
	"JOSHUA",
	"TRACY",
	"RYAN",
	"KAREN",
	"JEREMY",
	"DAWN",
	"PAUL",
	"SUSAN",
	"ANDREW",
	"ANDREA",
	"IAN",
	"MACKENZIE",
	"JASON",
	"KHLOE",
	"AYDEN",
	"SOPHIE",
	"ADAM",
	"KATHERINE",
	"PARKER",
	"PAISLEY",
	"COOPER",
	"MILA",
	"JUSTIN",
	"EVA",
	"XAVIER",
	"NAOMI",
	"NOLAN",
	"ELEANOR",
	"JACE",
	"GIANNA",
	"HUDSON",
	"MELANIE",
	"CARSON",
	"AUBREE",
	"BENTLEY",
	"FAITH",
	"LINCOLN",
	"KAYLA",
	"BLAKE",
	"PIPER",
	"EASTON",
	"MADELINE",
	"NATHANIEL",
	"LYDIA",
}

var lastNames = []string{
	"SMITH",
	"JOHNSON",
	"WILLIAMS",
	"JONES",
	"BROWN",
	"DAVIS",
	"MILLER",
	"WILSON",
	"MOORE",
	"TAYLOR",
	"ANDERSON",
	"THOMAS",
	"JACKSON",
	"WHITE",
	"HARRIS",
	"MARTIN",
	"THOMPSON",
	"GARCIA",
	"MARTINEZ",
	"ROBINSON",
	"CLARK",
	"RODRIGUEZ",
	"LEWIS",
	"LEE",
	"WALKER",
	"HALL",
	"ALLEN",
	"YOUNG",
	"HERNANDEZ",
	"KING",
	"WRIGHT",
	"LOPEZ",
	"HILL",
	"SCOTT",
	"GREEN",
	"ADAMS",
	"BAKER",
	"GONZALEZ",
	"NELSON",
	"CARTER",
	"MITCHELL",
	"PEREZ",
	"ROBERTS",
	"TURNER",
	"PHILLIPS",
	"CAMPBELL",
	"PARKER",
	"EVANS",
	"EDWARDS",
	"COLLINS",
	"STEWART",
	"SANCHEZ",
	"MORRIS",
	"ROGERS",
	"REED",
	"COOK",
	"MORGAN",
	"BELL",
	"MURPHY",
	"BAILEY",
	"RIVERA",
	"COOPER",
	"RICHARDSON",
	"COX",
	"HOWARD",
	"WARD",
	"TORRES",
	"PETERSON",
	"GRAY",
	"RAMIREZ",
	"JAMES",
	"WATSON",
	"BROOKS",
	"KELLY",
	"SANDERS",
	"PRICE",
	"BENNETT",
	"WOOD",
	"BARNES",
	"ROSS",
	"HENDERSON",
	"COLEMAN",
	"JENKINS",
	"PERRY",
	"POWELL",
	"LONG",
	"PATTERSON",
	"HUGHES",
	"FLORES",
	"WASHINGTON",
	"BUTLER",
	"SIMMONS",
	"FOSTER",
	"GONZALES",
	"BRYANT",
	"ALEXANDER",
	"RUSSELL",
	"GRIFFIN",
	"DIAZ",
	"HAYES ",
}
