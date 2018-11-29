package data

import (
	"encoding/csv"
	"strconv"
)

type WeightsInformation struct {
	cases map[string]float64
}

type WeightsEntry struct {
	Weight float64
	Found  bool
}

func NewWeightsInformation(sourceCSV *csv.Reader) (*WeightsInformation, error) {
	const COURTNO int = 2
	const WEIGHT int = 3

	cases := make(map[string]float64)

	records, err := sourceCSV.ReadAll()
	if err != nil {
		return nil, err
	}

	// Take slice from 1: because csv headers
	for _, record := range records[1:] {
		weight, err := strconv.ParseFloat(record[WEIGHT], 64)
		if err != nil {
			return nil, err
		}

		cases[record[COURTNO]] = weight
	}

	return &WeightsInformation{cases}, nil
}

func (w WeightsInformation) GetWeight(courtNumber string) WeightsEntry {
	weight, ok := w.cases[courtNumber]

	return WeightsEntry{weight, ok}
}
