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
	numberOfProp64Convictions := getNumberOfProp64Convictions()
	numberOfOtherConvictions := getNumberOfOtherConvictions()
	numberOfNonConvictions := getNumberOfNonConvictions()

	for i := 0; i < numberOfProp64Convictions; i++ {
		rows = append(rows, generateCase(personalInfos, true, true)...)
	}

	for i := 0; i < numberOfOtherConvictions; i++ {
		rows = append(rows, generateCase(personalInfos, false, true)...)
	}

	for i := 0; i < numberOfNonConvictions; i++ {
		rows = append(rows, generateCase(personalInfos, false, false)...)
	}

	return rows
}

func getNumberOfProp64Convictions() int {
	randomNumber := rand.Float64()
	result := 1
	if randomNumber > 0.95 {
		result = 3
	} else if randomNumber > 0.8 {
		result = 2
	}
	return result
}

func getNumberOfOtherConvictions() int {
	return int(math.Round(6 * rand.ExpFloat64()))
}

func getNumberOfNonConvictions() int {
	return int(math.Round(10 * rand.ExpFloat64()))
}

func generateCase(infos PersonalInfos, prop64 bool, conviction bool) [][]string {
	numberOfColumns := len(DojFullHeaders)
	row := make([]string, numberOfColumns)

	disposition := "DISMISSED"
	if conviction {
		disposition = "CONVICTED"
	}

	offense := "123 PC-OTHER OFFENSE"
	if prop64 {
		offense = "11357 HS-POSSESSION OF MARIJUANA"
	}

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
	row[OFN] = randomNDigitNumber(8)
	row[OFFENSE_DESCR] = offense
	row[OFFENSE_TOC] = "F"
	row[FE_NUM_CRT_CASE] = randomNDigitNumber(10)
	row[DISP_DESCR] = disposition
	row[CONV_STAT_DESCR] = randomSeverity()
	row[SENT_LENGTH] = "45"
	row[SENT_TIME_CODE] = "D"
	row[COMMENT_TEXT] = "THIS IS A COMMENT"
	return [][]string{row}
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

func randomSeverity() string {
	if rand.Float32() > 0.5 {
		return "FELONY"
	} else {
		return "MISDEMEANOR"
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
