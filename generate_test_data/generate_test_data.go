package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	. "gogen/data"
	. "gogen/exporter"
	"math"
	"math/rand"
	"path/filepath"
	"sort"
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

	rand.Seed(time.Now().UnixNano())

	testWriter, _ := NewWriter(filepath.Join(opts.OutputFolder, "generated_test_data.csv"), DojFullHeaders)

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
	totalNumberOfCases := numberOfProp64Convictions + numberOfOtherConvictions + numberOfNonConvictions

	dates := generateDates(totalNumberOfCases)
	counter := 0

	caseInfos := make(map[int]CaseInfos)

	for i := 0; i < numberOfProp64Convictions; i++ {
		caseInfos[counter] = CaseInfos{isProp64: true, isConviction: true}
		counter++
	}

	for i := 0; i < numberOfOtherConvictions; i++ {
		caseInfos[counter] = CaseInfos{isProp64: false, isConviction: true}
		counter++
	}

	for i := 0; i < numberOfNonConvictions; i++ {
		caseInfos[counter] = CaseInfos{isProp64: false, isConviction: false}
		counter++
	}

	counter = 0
	for _, info := range caseInfos {
		info.cycleNumber = counter + 101
		info.date = dates[counter]
		rows = append(rows, generateCase(personalInfos, info)...)
		counter++
	}

	return rows
}

func generateDates(n int) []time.Time {
	randomNumbers := make([]float64, n)
	dates := make([]time.Time, n)

	for i := 0; i < n; i++ {
		randomNumbers[i] = rand.Float64()
	}

	sort.Float64s(randomNumbers)

	currentUnixTime := float64(time.Now().Unix())

	for i := 0; i < n; i++ {
		dates[i] = time.Unix(int64(randomNumbers[i]*currentUnixTime), 0)
	}

	return dates
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

func generateCase(personalInfos PersonalInfos, caseInfos CaseInfos) [][]string {
	var rows [][]string

	offenses := getOffenses(caseInfos.isProp64)
	county := getCounty(caseInfos.isProp64)

	// handle arrest lines
	for i, offense := range offenses {
		countOrder := getCountOrder(caseInfos.cycleNumber, 1, i+1)
		rows = append(rows, generateRow(countOrder, personalInfos, caseInfos.date, county, false, caseInfos.isConviction, offense))
	}

	// handle court lines
	for i, offense := range offenses {
		countOrder := getCountOrder(caseInfos.cycleNumber, 2, i+1)
		newRow := generateRow(countOrder, personalInfos, caseInfos.date, county, true, caseInfos.isConviction, offense)

		// only set OFN for first count in case
		if i == 0 {
			newRow[OFN] = randomNDigitNumber(8)
		}

		// only set sentence for last count in case
		if caseInfos.isConviction && i == len(offenses)-1 {
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

func getCountOrder(cycle int, step int, count int) string {
	countOrder := fmt.Sprintf("%03d%03d%03d000", cycle, step, count)
	return countOrder
}

func getCounty(prop64 bool) string {
	county := opts.County
	if !prop64 && rand.Float32() > 0.75 {
		county = "NARNIA"
	}
	return county
}

func getOffenses(prop64 bool) []string {
	numberOfCounts := 1 + rand.Intn(4)
	choices := otherOffenses
	if prop64 {
		choices = prop64Offenses
	}
	offenses := make([]string, numberOfCounts)
	for i := range offenses {
		offenses[i] = randomChoice(choices)
	}
	return offenses
}

func generateRow(countOrder string, infos PersonalInfos, cycleDate time.Time, county string, court bool, conviction bool, offense string) []string {
	row := make([]string, numberOfColumns)
	for i := range row {
		row[i] = "  -  "
	}

	disposition := ""
	eventType := "ARREST/DETAINED/CITED"
	eventDate := cycleDate
	if court {
		eventDate = cycleDate.AddDate(0, 1, 0)
		eventType = "COURT ACTION"
		if conviction {
			disposition = "CONVICTED"
		} else {
			disposition = "DISMISSED"
		}
	}

	row[SUBJECT_ID] = infos.SubjectId
	row[CII_NUMBER] = infos.CII
	row[PRI_NAME] = infos.Name
	row[PRI_DOB] = infos.DOB
	row[PRI_SSN] = infos.SSN
	row[PRI_CDL] = infos.CDL
	row[CYC_DATE] = cycleDate.Format(dateFormatString)
	row[STP_EVENT_DATE] = eventDate.Format(dateFormatString)
	row[STP_TYPE_DESCR] = eventType
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
	return randomDate.Format(dateFormatString)
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

type CaseInfos struct {
	cycleNumber  int
	date         time.Time
	isConviction bool
	isProp64     bool
}

const dateFormatString = "20060102"

var prop64Offenses = []string{
	"11357 HS-POSSESSION OF MARIJUANA",
	"11357(A) HS-POSSESS CONCENTRATED CANNABIS",
	"11357(C) HS-POSS MARIJUANA OVER 1 OZ/28.5 GRAM",
	"11357A HS-POSSESSION OF MARIJUANA",
	"11357B2 HS-POSSESSION OF MARIJUANA",
	"11358 HS-PLANT/CULTIVATE/ETC MARIJUANA/HASH",
	"11358(C) HS-CULTIVATE MARIJUANA 6+ PLANTS",
	"11359 HS-MARIJUANA POSSESS FOR SALE",
	"11359(A) HS-POSSESS MARIJUANA/HASH FOR SALE",
	"11359(B) HS-POSSESS MARIJUANA FOR SALE",
	"11359(C) HS-MARIJUANA POSSESS FOR SALE",
	"11359B HS-MARIJUANA POSSESS FOR SALE",
	"11360 HS-SELL/TRANSPORT/ETC MARIJUANA/HASH",
	"11360 HS-SELL/TRANSPORT/ETC MARIJUANA/HASH",
	"11360 HS-SELL/TRANSPORT/ETC MARIJUANA/HASH",
	"11360(A) HS-GIVE/ETC MARIJ OVER 1 OZ/28.5 GRM",
	"11360(B) SELL/FURNISH/ETC MARIJUANA/HASH",
	"11360A HS-SELL/FURNISH/ETC MARIJUANA/HASH",
	"11360A2 HS-SELL/FURNISH/ETC MARIJUANA/HASH",
}

var otherOffenses = []string{
	"11350(A) HS-POSSESS NARC CONTROL SUBSTANCE",
	"11351.5 HS-POSS/PURCHASE COCAINE BASE F/SALE",
	"11352(A) HS-TRANSPORT/SELL NARC/CNTL SUB",
	"3056 PC-VIOLATION OF PAROLE:FELONY",
	"11364 HS-POSSESS CONTROL SUBSTANCE PARAPHERNA",
	"459 PC-BURGLARY",
	"166(A)(4) PC-CONTEMPT:DISOBEY COURT ORDER/ETC",
	"496(A) PC-RECEIVE/ETC KNOWN STOLEN PROPERTY",
	"14601.1(A) VC-DRIVE WHILE LIC SUSPEND/ETC",
	"11377(A) HS-POSSESS CONTROLLED SUBSTANCE",
	"245(A)(1) PC-FORCE/ADW NOT FIREARM:GBI LIKELY",
	"182(A)(1) PC-CONSPIRACY:COMMIT CRIME",
	"11352 HS-TRANSPORT/SELL NARCOTIC/CNTL SUB",
	"1203.2(A) PC-PROBATION VIOL:REARREST/REVOKE",
	"148.9(A) PC-FALSE ID TO SPECIFIC PEACE OFICERS",
	"11360(A) HS-GIVE/ETC MARIJ OVER 1 OZ/28.5 GRM",
	"148 PC-OBSTRUCTS/RESISTS PUBLIC OFFICER",
	"182 PC-CRIMINAL CONSPIRACY",
	"148(A)(1) PC-OBSTRUCT/ETC PUBLIC OFFICER/ETC",
	"242 PC-BATTERY",
	"459 PC-BURGLARY:SECOND DEGREE",
	"11351 HS-POSS/PURCHASE FOR SALE NARC/CNTL SUB",
	"148(A) PC-OBSTRUCTS/RESISTS PUBLIC OFFICER/ETC",
	"853.7 PC-FAIL TO APPEAR AFTER WRITTEN PROMISE",
	"11378 HS-POSSESS CONTROL SUBSTANCE FOR SALE",
	"11350 HS-POSSESS NARCOTIC CONTROL SUBSTANCE",
	"466 PC-POSSESS/ETC BURGLARY TOOLS",
	"10851(A) VC-TAKE VEH W/O OWN CONSENT/VEH THEFT",
	"23152(A) VC-DUI ALCOHOL/DRUGS",
	"666 PC-PETTY THEFT W/PR JAIL:SPEC OFFENSES",
	"422 PC-THREATEN CRIME WITH INTENT TO TERRORIZE",
	"12500(A) VC-DRIVE W/O LICENSE",
	"40508(A) VC-FAIL TO APPEAR:WRITTEN PROMISE",
	"273.5 PC-INFLICT CORPORAL INJ ON SPOUSE/COHAB",
	"211 PC-ROBBERY",
	"647(F) PC-DISORDERLY CONDUCT:INTOX DRUG/ALCOH",
	"11550(A) HS-USE/UNDER INFL CONTRLD SUBSTANCE",
	"273.5(A) PC-INFLICT CORPORAL INJ SPOUSE/COHAB",
	"1203.2 PC-PROBATION VIOL:REARREST/REVOKE",
	"488 PC-PETTY THEFT",
	"11550 HS-USE/UNDER INFLUENCE CONTROL SUBST",
	"148.9 PC-FALSE IDENTIFICATION TO PEACE OFFICER",
	"23152(B) VC-DUI ALCOHOL/0.08 PERCENT",
	"243(E)(1) PC-BAT:SPOUSE/EX SP/DATE/ETC",
	"212.5(C) PC-ROBBERY:SECOND DEGREE",
	"32 PC-ACCESSORY",
	"10851 VC-TAKE VEH W/O OWN CONSENT/VEH THEFT",
	"11377(A) HS-POSSESS CNTL SUBSTANCE",
	"4140 BP-POSSESS HYPODERMIC NEEDLE/SYRINGE",
	"647(B) PC-DISORDERLY CONDUCT:PROSTITUTION",
	"243(B) PC-BATTERY PEACE OFCR/EMERG PERSNL/ETC",
	"484(A)/490.5 PC-THEFT/PETTY THEFT MERCHANDISE",
	"496.1 PC-RECEIVE/ETC KNOWN STOLEN PROPERTY",
	"25620 BP-POSS OPEN CONTAINER OF ALCOHOL:PUBLIC",
	"602(L) PC-TRESPASS:OCCUPY PROPERTY W/O CONSENT",
	"4149 BP-POSSESS HYPODERMIC NEEDLE/SYRINGE",
	"496 PC-RECEIVE/ETC KNOWN STOLEN PROPERTY",
	"288(A) PC-LEWD OR LASCIV ACTS", // 290 registerable offense
	"187 PC-MURDER:SECOND DEGREE",   // Superstrike
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
