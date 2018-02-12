package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

const DateFormat = "2006-01-02"

type TimeSeries struct {
	granularity time.Duration
	records map[time.Time]float64
}

func NewTimeSeries() *TimeSeries {
	return &TimeSeries{
		granularity: 24 * time.Hour,
		records: make(map[time.Time]float64),
	}
}

func (t *TimeSeries) Read(db string) error {
	csvFile, err := os.Open(db)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.FieldsPerRecord = 2
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
			timestamp, err = time.Parse(DateFormat, values[0])
			if err != nil {
				return err
			}
		}
		datum, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return err
		}

		t.records[timestamp] = datum
	}

	return nil
}

// New TimeSeries with data since timestamp
func (t *TimeSeries) Since(ts time.Time) *TimeSeries {
	records := make(map[time.Time]float64)
	for timestamp, datum := range t.records {
		if !timestamp.Before(ts) {
			records[timestamp] = datum
		}
	}

	return &TimeSeries{
		granularity: t.granularity,
		records: records,
	}
}

// Resample by date
func (t *TimeSeries) Resample(d time.Duration) {

	// Make a list of values for each duration
	data := make(map[time.Time][]float64)
	for timestamp, datum := range t.records {
		timestamp = timestamp.Truncate(d)
		values, ok := data[timestamp]
		if !ok {
			data[timestamp] = make([]float64, 0)
		}
		data[timestamp] = append(values, datum)
	}

	// Calculate mean for each duration
	records := make(map[time.Time]float64)
	for timestamp, values := range data {
		var total float64
		for _, value := range values {
			total += value
		}
		records[timestamp] = total / float64(len(values))
	}

	// Use resampled values
	t.granularity = d
	t.records = records
}

// Fill in missing days using linear interpolation
func (t *TimeSeries) Interpolate() {
	// Put the timestamps in chronological order
	timestamps := make([]time.Time, len(t.records))
	i := 0
	for key := range t.records {
		timestamps[i] = key
		i++
	}

	sort.Slice(timestamps, func(i, j int) bool {
		return timestamps[i].Before(timestamps[j])
	})

	// Look for gaps (no records for a duration)
	duration := t.granularity
	var last time.Time
	for index, timestamp := range timestamps {
		timestamp = timestamp.Round(duration)
		if index == 0 {
			last = timestamp
			continue
		}

		interval := timestamp.Sub(last)
		periods := float64(interval / duration)
		if periods == 0 {
			last = timestamp
			continue
		}

		start_value := t.records[last]
		end_value := t.records[timestamp]
		step := (end_value - start_value) / float64(periods)
		for at, period := last, 0.0; period < periods; period++ {
			if period > 0 {
				t.records[at] = start_value + step*period
			}

			at = at.Add(duration)
		}

		last = timestamp
	}
}
