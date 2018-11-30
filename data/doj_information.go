package data

import "encoding/csv"

type DOJInformation struct {
	Rows         [][]string
	DOJHistories map[string]DOJHistory
}

type MatchData struct {
	history           DOJHistory
	cii               bool
	ssn               bool
	courtno           bool
	name_and_dob      bool
	weak_name_and_dob bool
	cdl               bool
	match_strength    int
}

func (information DOJInformation) findDOJHistory(entry CMSEntry) DOJHistory {
	return DOJHistory{}
}

func NewDOJInformation(sourceCSV *csv.Reader) (*DOJInformation, error) {
	const SubjectIDIndex int = 0

	rows, err := sourceCSV.ReadAll()
	if err != nil {
		panic(err)
	}

	info := DOJInformation{
		Rows:         rows,
		DOJHistories: make(map[string]DOJHistory),
	}

	for _, row := range rows {
		info.DOJHistories[row[SubjectIDIndex]].PushRow(NewDOJRow(row))
	}

	return &info, nil
}
