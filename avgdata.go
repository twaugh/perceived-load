package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type Record struct {
	timestamp time.Time
	load float64
}

type AvgData struct {
	records []Record
}

func NewAvgData(DB string) (*AvgData, error) {
	csvFile, err := os.Open(DB)
	if err != nil {
		return nil, err
	}

	var records []Record
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		values, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		timestamp, err := time.Parse(time.RFC3339, values[0])
		if err != nil {
			return nil, err
		}
		load, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return nil, err
		}

		records = append(records, Record{
			timestamp: timestamp,
			load: load,
		})
	}

	return &AvgData{
		records: records,
	}, nil
}
