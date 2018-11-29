package data

import (
	"strings"
	"time"
)

type EligibilityInfo struct {
	QFinalSum float64
	Over1Lb   string
}

type CMSEntry struct {
	CourtNumber     string
	Level           string
	SSN             string
	CII             string
	Charge          string
	IncidentNumber  string
	Name            string
	CDL             string
	DateOfBirth     time.Time
	DispositionDate time.Time
	RawRow          []string
	EligibilityInfo EligibilityInfo
}

func NewCMSEntry(record []string) CMSEntry {
	const (
		COURTNO    int = 0
		INCIDENTNO int = 2
		NAME       int = 3
		CHARGE     int = 12
		LEVEL      int = 13
		SSN        int = 24
		CII        int = 22
		CDL        int = 25
		DOB        int = 20
		DISPODATE  int = 7
	)

	dob := parseDate(record[DOB])
	dispositionDate := parseDate(record[DISPODATE])

	return CMSEntry{
		CourtNumber:     record[COURTNO],
		Level:           record[LEVEL],
		SSN:             record[SSN],
		CII:             record[CII],
		Charge:          strings.TrimSpace(record[CHARGE]),
		IncidentNumber:  record[INCIDENTNO],
		Name:            strings.TrimSpace(record[NAME]),
		CDL:             strings.SplitN(record[CDL], " ", 2)[0],
		DateOfBirth:     dob,
		DispositionDate: dispositionDate,
		RawRow:          record,
	}
}

func (c CMSEntry) FormattedName() string {
	nameParts := strings.Split(c.Name, "/")

	if len(nameParts) > 1 {
		lastCommaFirst := strings.Join(nameParts[0:2], ",")
		return strings.Join(append([]string{lastCommaFirst}, nameParts[2:]...), " ")
	}

	return nameParts[0]
}

func parseDate(s string) time.Time {
	t, err := time.Parse("1/2/06", s)
	if err != nil {
		t = time.Time{}
	}
	if time.Now().Before(t) {
		t = t.AddDate(-100, 0, 0)
	}
	return t
}
