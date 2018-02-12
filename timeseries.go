package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

const DateFormat = "2006-01-02"

type Record struct {
	timestamp time.Time
	datum float64
}

type TimeSeries struct {
	granularity time.Duration

	// records must remain sorted by timestamp
	records []*Record
}

func (t *TimeSeries) sort() {
	sort.Slice(t.records, func (i, j int) bool {
		return t.records[i].timestamp.Before(t.records[j].timestamp)
	})
}

func NewTimeSeries() *TimeSeries {
	return &TimeSeries{
		granularity: 24 * time.Hour,
		records: make([]*Record, 0),
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

		t.records = append(t.records, &Record{
			timestamp: timestamp,
			datum: datum,
		})
	}

	t.sort()
	return nil
}

type InvalidTimestamp time.Time

func (t InvalidTimestamp) Error() string {
	return fmt.Sprintf("invalid timestamp %v", time.Time(t))
}

func (t *TimeSeries) LookUp(timestamp time.Time) (float64, error) {
	n := len(t.records)
	i := sort.Search(n, func(i int) bool {
		return !t.records[i].timestamp.Before(timestamp)
	})
	if i == n || t.records[i].timestamp != timestamp {
		return 0, InvalidTimestamp(timestamp)
	}
	return t.records[i].datum, nil
}

// New TimeSeries with data since timestamp
func (t *TimeSeries) Since(timestamp time.Time) *TimeSeries {
	i := sort.Search(len(t.records), func(i int) bool {
		return !t.records[i].timestamp.Before(timestamp)
	})
	return &TimeSeries{
		granularity: t.granularity,
		records: t.records[i:],
	}
}

// Resample by date
func (t *TimeSeries) Resample(d time.Duration) {

	// Make a list of values for each duration
	data := make(map[time.Time][]float64)
	for _, record := range t.records {
		timestamp := record.timestamp
		timestamp = timestamp.Truncate(d)
		values, ok := data[timestamp]
		if !ok {
			data[timestamp] = make([]float64, 0)
		}
		data[timestamp] = append(values, record.datum)
	}

	// Calculate mean for each duration
	records := make([]*Record, len(data))
	i := 0
	for timestamp, values := range data {
		var total float64
		for _, value := range values {
			total += value
		}
		records[i] = &Record{
			timestamp: timestamp,
			datum: total / float64(len(values)),
		}
		i++
	}

	// Use resampled values
	t.granularity = d
	t.records = records
	t.sort()
}

// Fill in missing days using linear interpolation
func (t *TimeSeries) Interpolate() {
	// Look for gaps (no records for a duration)
	duration := t.granularity
	var last *Record
	missing := make([]*Record, 0)
	for index, record := range t.records {
		at := record.timestamp.Round(duration)
		if index == 0 {
			last = record
			continue
		}

		interval := at.Sub(last.timestamp)
		periods := float64(interval / duration)
		if periods == 0 {
			last = record
			continue
		}

		start_value := last.datum
		end_value := record.datum
		step := (end_value - start_value) / float64(periods)
		for period := periods - 1; period > 0; period-- {
			at = at.Add(-duration)
			missing = append(missing, &Record{
				timestamp: at,
				datum: start_value + step*period,
			})
		}

		last = record
	}

	// Sort in the interpolated records
	t.records = append(t.records, missing...)
	t.sort()
}
