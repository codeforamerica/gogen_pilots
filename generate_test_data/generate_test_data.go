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

var numberOfColumns = len(DojFullHeaders)

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

	cycleCounter := 101

	for i := 0; i < numberOfProp64Convictions; i++ {
		rows = append(rows, generateCase(cycleCounter, personalInfos, true, true)...)
		cycleCounter++
	}

	for i := 0; i < numberOfOtherConvictions; i++ {
		rows = append(rows, generateCase(cycleCounter, personalInfos, false, true)...)
		cycleCounter++
	}

	for i := 0; i < numberOfNonConvictions; i++ {
		rows = append(rows, generateCase(cycleCounter, personalInfos, false, false)...)
		cycleCounter++
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

func generateCase(cycle int, infos PersonalInfos, prop64 bool, conviction bool) [][]string {
	var rows [][]string
	numberOfCounts := 1 + rand.Intn(4)
	county := opts.County
	if !prop64 && rand.Float32() > 0.75 {
		county = "NARNIA"
	}

	step := 1

	for i := 0; i < numberOfCounts; i++ {
		countOrder := fmt.Sprintf("%03d%03d%03d000", cycle, step, i+1)
		newRow := generateRow(countOrder, infos, county, conviction, prop64)

		// only set OFN for first count in case
		if i == 0 {
			newRow[OFN] = randomNDigitNumber(8)
		}

		// only set sentence for last count in case
		if conviction && i == numberOfCounts-1 {
			numSentenceParts := rand.Intn(3) + 1
			for j := 0; j < numSentenceParts; j++ {
				sentenceRow := make([]string, numberOfColumns)
				copy(sentenceRow, newRow)

				choice := rand.Intn(3)
				if choice == 0 {
					sentenceRow[SENT_LENGTH] = strconv.Itoa(rand.Intn(100) + 1)
					sentenceRow[SENT_TIME_CODE] = "D"
				} else if choice == 1 {
					sentenceRow[SENT_LENGTH] = strconv.Itoa(rand.Intn(24) + 1)
					sentenceRow[SENT_TIME_CODE] = "M"
				} else {
					sentenceRow[SENT_LENGTH] = strconv.Itoa(rand.Intn(10) + 1)
					sentenceRow[SENT_TIME_CODE] = "Y"
				}

				rows = append(rows, sentenceRow)
			}
		} else {
			rows = append(rows, newRow)
		}
	}

	return rows
}

func generateRow(countOrder string, infos PersonalInfos, county string, conviction bool, prop64 bool) []string {
	row := make([]string, numberOfColumns)
	for i := range row {
		row[i] = "  -  "
	}

	disposition := "DISMISSED"
	if conviction {
		disposition = "CONVICTED"
	}
	offense := randomChoice(otherOffenses)
	if prop64 {
		offense = randomChoice(prop64Offenses)
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
	row[STP_ORI_CNTY_NAME] = county
	row[CNT_ORDER] = countOrder
	row[OFN] = ""
	row[OFFENSE_DESCR] = offense
	row[OFFENSE_TOC] = "F"
	row[FE_NUM_CRT_CASE] = ""
	row[DISP_DESCR] = disposition
	row[CONV_STAT_DESCR] = randomSeverity()
	row[SENT_LENGTH] = ""
	row[SENT_TIME_CODE] = ""
	row[COMMENT_TEXT] = "THIS IS A COMMENT"
	return row
}

func generatePersonalInfos() PersonalInfos {
	counters.SubjectId++

	return PersonalInfos{
		SubjectId: strconv.Itoa(counters.SubjectId),
		CII:       randomNDigitNumber(10),
		Name:      randomChoice(lastNames) + "," + randomChoice(firstNames) + " " + randomLetter(),
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

var prop64Offenses = []string{
	"11357 HS-POSSESSION OF MARIJUANA",
	"11358 HS-CULTIVATION OF MARIJUANA",
	"11359 HS-MARIJUANA POSSESS FOR SALE",
	"11360 HS-SELL/TRANSPORT/ETC MARIJUANA/HASH",
}

var otherOffenses = []string{
	"11350(a) HS-POSSESS NARC CONTROL SUBSTANCE",
	"11351.5 HS-POSS/PURCHASE COCAINE BASE F/SALE",
	"11352(a) HS-TRANSPORT/SELL NARC/CNTL SUB",
	"3056 PC-VIOLATION OF PAROLE:FELONY",
	"11364 HS-POSSESS CONTROL SUBSTANCE PARAPHERNA",
	"459 PC-BURGLARY",
	"166(a)(4) PC-CONTEMPT:DISOBEY COURT ORDER/ETC",
	"496(a) PC-RECEIVE/ETC KNOWN STOLEN PROPERTY",
	"14601.1(a) VC-DRIVE WHILE LIC SUSPEND/ETC",
	"11377(a) HS-POSSESS CONTROLLED SUBSTANCE",
	"245(a)(1) PC-FORCE/ADW NOT FIREARM:GBI LIKELY",
	"182(a)(1) PC-CONSPIRACY:COMMIT CRIME",
	"11352 HS-TRANSPORT/SELL NARCOTIC/CNTL SUB",
	"1203.2(a) PC-PROBATION VIOL:REARREST/REVOKE",
	"148.9(a) PC-FALSE ID TO SPECIFIC PEACE OFICERS",
	"11360(a) HS-GIVE/ETC MARIJ OVER 1 OZ/28.5 GRM",
	"148 PC-OBSTRUCTS/RESISTS PUBLIC OFFICER",
	"182 PC-CRIMINAL CONSPIRACY",
	"148(a)(1) PC-OBSTRUCT/ETC PUBLIC OFFICER/ETC",
	"242 PC-BATTERY",
	"459 PC-BURGLARY:SECOND DEGREE",
	"11351 HS-POSS/PURCHASE FOR SALE NARC/CNTL SUB",
	"148(a) PC-OBSTRUCTS/RESISTS PUBLIC OFFICER/ETC",
	"853.7 PC-FAIL TO APPEAR AFTER WRITTEN PROMISE",
	"11378 HS-POSSESS CONTROL SUBSTANCE FOR SALE",
	"11350 HS-POSSESS NARCOTIC CONTROL SUBSTANCE",
	"466 PC-POSSESS/ETC BURGLARY TOOLS",
	"10851(a) VC-TAKE VEH W/O OWN CONSENT/VEH THEFT",
	"23152(a) VC-DUI ALCOHOL/DRUGS",
	"666 PC-PETTY THEFT W/PR JAIL:SPEC OFFENSES",
	"422 PC-THREATEN CRIME WITH INTENT TO TERRORIZE",
	"12500(a) VC-DRIVE W/O LICENSE",
	"40508(a) VC-FAIL TO APPEAR:WRITTEN PROMISE",
	"273.5 PC-INFLICT CORPORAL INJ ON SPOUSE/COHAB",
	"211 PC-ROBBERY",
	"647(f) PC-DISORDERLY CONDUCT:INTOX DRUG/ALCOH",
	"11550(a) HS-USE/UNDER INFL CONTRLD SUBSTANCE",
	"273.5(a) PC-INFLICT CORPORAL INJ SPOUSE/COHAB",
	"1203.2 PC-PROBATION VIOL:REARREST/REVOKE",
	"488 PC-PETTY THEFT",
	"11550 HS-USE/UNDER INFLUENCE CONTROL SUBST",
	"148.9 PC-FALSE IDENTIFICATION TO PEACE OFFICER",
	"23152(b) VC-DUI ALCOHOL/0.08 PERCENT",
	"243(e)(1) PC-BAT:SPOUSE/EX SP/DATE/ETC",
	"212.5(c) PC-ROBBERY:SECOND DEGREE",
	"32 PC-ACCESSORY",
	"10851 VC-TAKE VEH W/O OWN CONSENT/VEH THEFT",
	"11377(a) HS-POSSESS CNTL SUBSTANCE",
	"4140 BP-POSSESS HYPODERMIC NEEDLE/SYRINGE",
	"647(b) PC-DISORDERLY CONDUCT:PROSTITUTION",
	"243(b) PC-BATTERY PEACE OFCR/EMERG PERSNL/ETC",
	"484(a)/490.5 PC-THEFT/PETTY THEFT MERCHANDISE",
	"496.1 PC-RECEIVE/ETC KNOWN STOLEN PROPERTY",
	"25620 BP-POSS OPEN CONTAINER OF ALCOHOL:PUBLIC",
	"602(l) PC-TRESPASS:OCCUPY PROPERTY W/O CONSENT",
	"4149 BP-POSSESS HYPODERMIC NEEDLE/SYRINGE",
	"496 PC-RECEIVE/ETC KNOWN STOLEN PROPERTY",
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
