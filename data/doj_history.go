package data

import "time"

type DOJHistory struct {
	Name              string
	WeakName          string
	SubjectID         string
	CII               string
	Convictions       [][]string
	Rows              []DOJRow
	DOB               time.Time
	SSN               string
	CDL               string
	PC290Registration bool
}

func (dh DOJHistory) PushRow(row []string) {}
