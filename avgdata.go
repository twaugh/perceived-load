package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type AvgData struct {
	records map[time.Time]float64
}

func NewAvgData() *AvgData {
	return &AvgData{
		records: make(map[time.Time]float64),
	}
}

func (a *AvgData) Read(DB string) error {
	csvFile, err := os.Open(DB)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		values, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		timestamp, err := time.Parse(time.RFC3339, values[0])
		if err != nil {
			return err
		}
		load, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return err
		}

		a.records[timestamp] = load
	}

	return nil
}

func (a *AvgData) Resample() {
	// Resample by date
	data := make(map[time.Time][]float64)

	for timestamp, datum := range a.records {
		timestamp = time.Date(
			timestamp.Year(),
			timestamp.Month(),
			timestamp.Day(),
			0, 0, 0, 0, time.UTC,
		)
		values, ok := data[timestamp]
		if !ok {
			data[timestamp] = make([]float64, 0)
		}
		data[timestamp] = append(values, datum)
	}

	records := make(map[time.Time]float64)
	for timestamp, values := range data {
		var total float64
		for _, value := range values {
			total += value
		}
		records[timestamp] = total / float64(len(values))
	}

	a.records = records
}
