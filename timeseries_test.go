package main

import (
	"testing"
	"time"
)

func TestNewTimeSeries(t *testing.T) {
	ts := NewTimeSeries()
	if ts.records == nil {
		t.Errorf("incomplete initialization")
	}
}

func TestTimeSeriesRead(t *testing.T) {
	ts := NewTimeSeries()
	if err := ts.Read("testdata/simple.csv"); err != nil {
		t.Fatal(err)
	}
	length := len(ts.records)
	if length != 4 {
		t.Errorf("unexpected series length %v", length)
		return
	}

	testtime := time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC)
	record, err := ts.LookUp(testtime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record != 1.0 {
		t.Errorf("unexpected record value %r", record)
	}
}

func TestTimeSeriesRead_DateOnly(t *testing.T) {
	ts := NewTimeSeries()
	if err := ts.Read("testdata/dateonly.csv"); err != nil {
		t.Fatal(err)
	}
	length := len(ts.records)
	if length != 1 {
		t.Errorf("unexpected series length %v", length)
	}
}

func TestTimeSeriesSince(t *testing.T) {
	ts := NewTimeSeries()
	if err := ts.Read("testdata/simple.csv"); err != nil {
		t.Fatal(err)
	}
	ts = ts.Since(time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC))
	length := len(ts.records)
	if length != 2 {
		t.Errorf("unexpected series length %v", length)
	}
}

func TestTimeSeriesResample(t *testing.T) {
	ts := NewTimeSeries()
	if err := ts.Read("testdata/simple.csv"); err != nil {
		t.Fatal(err)
	}
	ts.Resample(24 * time.Hour)
	testtime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	record, err := ts.LookUp(testtime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record != 1.5 {
		t.Errorf("unexpected record value %r", record)
	}
}

func TestTimeSeriesInterpolate(t *testing.T) {
	ts := NewTimeSeries()
	if err := ts.Read("testdata/simple.csv"); err != nil {
		t.Fatal(err)
	}
	ts.Interpolate()
	expected := map[int]float64{
		1: 11,
		2: 12,
		3: 13,
		4: 14,
		5: 15,
	}
	for mday, expected_value := range expected {
		index := time.Date(2018, 2, mday, 0, 0, 0, 0, time.UTC)
		value, err := ts.LookUp(index)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		if value != expected_value {
			t.Errorf("for mday %v expected %v but got %v",
				mday, expected_value, value)
		}
	}
}
