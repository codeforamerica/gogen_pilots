package data

import (
	"encoding/csv"
	"errors"
	"strconv"
)

const COURTNO int = 2
const WEIGHT int = 3

type WeightsInformation struct {
	cases map[string]float64
}

func NewWeightsInformation(sourceCSV *csv.Reader) (*WeightsInformation, error) {
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

func (w WeightsInformation) Under1LB(courtNumber string) (bool, error) {
	weight, ok := w.cases[courtNumber]
	if !ok {
		return false, errors.New("court number did not exist")
	}
	return weight <= 453.592, nil
}

func (w WeightsInformation) GetWeight(courtNumber string) (float64, error) {
	weight, ok := w.cases[courtNumber]
	if !ok {
		return 0, errors.New("court number did not exist")
	}

	return weight, nil
}
