package data

import "encoding/csv"

type CMSInformation struct {
	Entries []CMSEntry
}

func NewCMSInformation(sourceCSV *csv.Reader) (*CMSInformation, error) {
	records, err := sourceCSV.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make([]CMSEntry, len(records))

	for i, record := range records {
		entries[i] = NewCMSEntry(record)
	}

	return &CMSInformation{entries}, nil
}
