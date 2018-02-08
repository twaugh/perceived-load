package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Load float64
}

type AvgData struct {
	Records []Record
}

func NewAvgData(DB string) (*AvgData, error) {
	csvFile, err := os.Open(DB)
	if err != nil {
		return nil, err
	}

	var records []Record
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		load, err := strconv.ParseFloat(line[0], 64)
		if err != nil {
			return nil, err
		}
		records = append(records, Record{
			Load: load,
		})
	}

	return &AvgData{
		Records: records,
	}, nil
}
