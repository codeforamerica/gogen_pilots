package data

import (
	"strings"
	"time"
)

type CMSEntry struct {
	CourtNumber    string
	Level          string
	SSN            string
	CII            string
	Charge         string
	IncidentNumber string
	Name           string
	CDL            string
	DateOfBirth    time.Time
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
	)

	dob, err := time.Parse("1/2/06", record[DOB])
	if err != nil {
		panic(err)
	}
	if time.Now().Before(dob) {
		dob = dob.AddDate(-100, 0, 0)
	}

	return CMSEntry{
		CourtNumber:    record[COURTNO],
		Level:          record[LEVEL],
		SSN:            record[SSN],
		CII:            record[CII],
		Charge:         strings.TrimSpace(record[CHARGE]),
		IncidentNumber: record[INCIDENTNO],
		Name:           strings.TrimSpace(record[NAME]),
		CDL:            strings.SplitN(record[CDL], " ", 2)[0],
		DateOfBirth:    dob,
	}
}

func (c CMSEntry) FormattedName() string {
	nameParts := strings.Split(c.Name, "/")
	lastCommaFirst := strings.Join(nameParts[0:2], ",")
	return strings.Join(append([]string{lastCommaFirst}, nameParts[2:]...), " ")
}
