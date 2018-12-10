package data

import (
	"strings"
	"time"
)

type CMSEntry struct {
	CourtNumber     string
	FormattedCourtNumber string
	Level           string
	SSN             string
	CII             string
	Charge          string
	IncidentNumber  string
	Name            string
	FormattedName   string
	WeakName        string
	CDL             string
	DateOfBirth     time.Time
	DispositionDate time.Time
	BookingDate     time.Time
	RawRow          []string
}

func NewCMSEntry(record []string) CMSEntry {
	const DateFormat string = "010206"
	const (
		COURTNO = iota
		IND
		INCIDENTNO
		NAME
		CASE_CLAS
		DISPO
		DISPO_DESCRIPTION
		DISPO_DATE
		ACTION_NUMER
		FILED_CHARGE
		FILED_CHARGE_CLASS
		FILED_CHARGE_DATE
		CURRENT_CHARGE
		CURRENT_CHARGE_CLASS
		CURRENT_CHARGE_DESCRIPTION
		CHARGE_DISPO
		CHARGE_DISPO_DESCRIPTION
		CHARGE_DISPO_DATE
		BOOKED_CHARGE
		BOOKED_CHARGE_LEVEL
		BOOKED_CHARGE_DATE
		RACE
		SEX
		DOB
		SFNO
		CII
		FBI
		SSN
		CDL
	)

	dob := parseDate(DateFormat, record[DOB])
	dispositionDate := parseDate(DateFormat, record[DISPO_DATE])
	bookingDate := parseDate(DateFormat, record[BOOKED_CHARGE_DATE])
	formattedName := formatName(strings.TrimSpace(record[NAME]))
	firstLast := strings.Split(formattedName, " ")[0]
	cii := formatCII(record[CII])

	return CMSEntry{
		CourtNumber:     record[COURTNO],
		FormattedCourtNumber: formatCourtNumber(record[COURTNO]),
		Level:           record[CURRENT_CHARGE_CLASS],
		SSN:             record[SSN],
		CII:             formatCII(cii),
		Charge:          strings.TrimSpace(record[CURRENT_CHARGE]),
		IncidentNumber:  record[INCIDENTNO],
		Name:            strings.TrimSpace(record[NAME]),
		FormattedName:   formattedName,
		WeakName:        firstLast,
		CDL:             strings.SplitN(record[CDL], " ", 2)[0],
		DateOfBirth:     dob,
		DispositionDate: dispositionDate,
		BookingDate:     bookingDate,
		RawRow:          record,
	}
}

func (entry CMSEntry) MJCharge() bool {
	return strings.HasPrefix(entry.Charge, "11357") ||
		strings.HasPrefix(entry.Charge, "11358") ||
		strings.HasPrefix(entry.Charge, "11359") ||
		strings.HasPrefix(entry.Charge, "11360")
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

func formatCourtNumber(number string) string {
	if number == "" {
		return number
	}

	for len(number) < 8 {
		number = "0" + number
	}
	return number[len(number)-8:]
}
