package data

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

const COURTNO int = 2
const WEIGHT int = 3

type WeightsInformation struct {
	cases map[string]float64
}

func NewWeightsInformation(sourceCSV *csv.Reader) (*WeightsInformation, error) {
	wi := WeightsInformation{}
	wi.cases = make(map[string]float64)

	//ignore headers
	_, err := sourceCSV.Read()
	if err != nil {
		return nil, err
	}

	for {
		record, err := sourceCSV.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		weight, err := strconv.ParseFloat(record[WEIGHT], 64)
		if err != nil {
			return nil, err
		}

		wi.cases[record[COURTNO]] = weight
	}

	return &wi, nil
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
