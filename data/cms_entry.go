package data

import (
	"strings"
	"time"
)

type CMSEntry struct {
	CourtNumber     string
	Level           string
	SSN             string
	CII             string
	Charge          string
	IncidentNumber  string
	Name            string
	FormattedName   string
	CDL             string
	DateOfBirth     time.Time
	DispositionDate time.Time
	RawRow          []string
}

func NewCMSEntry(record []string) CMSEntry {
	const (
		COURTNO    int    = 0
		INCIDENTNO int    = 2
		NAME       int    = 3
		CHARGE     int    = 12
		LEVEL      int    = 13
		SSN        int    = 24
		CII        int    = 22
		CDL        int    = 25
		DOB        int    = 20
		DISPODATE  int    = 7
		DateFormat string = "1/2/06"
	)

	dob := parseDate(DateFormat, record[DOB])
	dispositionDate := parseDate(DateFormat, record[DISPODATE])
	formattedName := formatName(strings.TrimSpace(record[NAME]))
	cii := formatCII(record[CII])

	return CMSEntry{
		CourtNumber:     record[COURTNO],
		Level:           record[LEVEL],
		SSN:             record[SSN],
		CII:             formatCII(cii),
		Charge:          strings.TrimSpace(record[CHARGE]),
		IncidentNumber:  record[INCIDENTNO],
		Name:            strings.TrimSpace(record[NAME]),
		FormattedName:   formattedName,
		CDL:             strings.SplitN(record[CDL], " ", 2)[0],
		DateOfBirth:     dob,
		DispositionDate: dispositionDate,
		RawRow:          record,
	}
}

func formatName(name string) string {
	nameParts := strings.Split(name, "/")

	if len(nameParts) > 1 {
		lastCommaFirst := strings.Join(nameParts[0:2], ",")
		return strings.Join(append([]string{lastCommaFirst}, nameParts[2:]...), " ")
	}

	return nameParts[0]
}

func formatCII(cii string) string {
	if cii == "" {
	return cii
	}
	for len(cii) < 8 {
		cii = "0" + cii
	}
	return cii[len(cii)-8:]
}
