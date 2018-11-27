package data

import "strings"

type CMSEntry struct {
	CourtNumber    string
	Level          string
	SSN            string
	CII            string
	Charge         string
	IncidentNumber string
	Name           string
	CDL            string
	FormattedName  string
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
	)

	return CMSEntry{
		CourtNumber:    record[COURTNO],
		Level:          record[LEVEL],
		SSN:            record[SSN],
		CII:            record[CII],
		Charge:         strings.TrimSpace(record[CHARGE]),
		IncidentNumber: record[INCIDENTNO],
		Name:           strings.TrimSpace(record[NAME]),
		CDL:            strings.SplitN(record[CDL], " ", 2)[0],
	}
}
