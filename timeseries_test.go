package main

import (
	"io/ioutil"
	"os"
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
	record, err := ts.Lookup(testtime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record != 1.0 {
		t.Errorf("unexpected record value %#v", record)
	}

	testtime = time.Date(2018, 1, 1, 11, 0, 0, 0, time.UTC)
	record, err = ts.Lookup(testtime)
	err, ok := err.(InvalidTimestamp)
	if !ok {
		t.Fatalf("expected InvalidTimestamp")
	}
	err.Error()
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

func TestTimeSeriesRead_Error(t *testing.T) {
	ts := NewTimeSeries()
	err := ts.Read("testdata/csverror.csv")
	if err == nil {
		t.Fatal("expected error")
	}
	length := len(ts.records)
	if length != 1 {
		t.Errorf("unexpected series length %v", length)
	}
}

func TestTimeSeriesRead_DateError(t *testing.T) {
	ts := NewTimeSeries()
	err := ts.Read("testdata/dateerror.csv")
	if err == nil {
		t.Fatal("expected error")
	}
	length := len(ts.records)
	if length != 1 {
		t.Errorf("unexpected series length %v", length)
	}
}

func TestTimeSeriesRead_FloatError(t *testing.T) {
	ts := NewTimeSeries()
	err := ts.Read("testdata/floaterror.csv")
	if err == nil {
		t.Fatal("expected error")
	}
	length := len(ts.records)
	if length != 1 {
		t.Errorf("unexpected series length %v", length)
	}
}

func TestTimeSeriesWrite(t *testing.T) {
	ts := NewTimeSeries()
	if err := ts.Read("testdata/simple.csv"); err != nil {
		t.Fatal(err)
	}
	length := len(ts.records)
	tempfile, err := ioutil.TempDir("", "write")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = os.RemoveAll(tempfile); err != nil {
			t.Error(err)
		}
	}()
	tempfile += "/output.csv"
	err = ts.Write(tempfile)
	if err != nil {
		t.Fatal(err)
	}

	ts = NewTimeSeries()
	if err := ts.Read(tempfile); err != nil {
		t.Fatal(err)
	}
	got := len(ts.records)
	if got != length {
		t.Errorf("read %v records but expected %v", got, length)
	}
}

func TestTimeSeriesAdd(t *testing.T) {
	ts := NewTimeSeries()
	index := time.Date(2018, 5, 1, 0, 0, 0, 0, time.UTC)
	if _, err := ts.Lookup(index); err == nil {
		t.Fatal("new TimeSeries not empty")
	}
	expected := 5.1
	ts.Add(index, expected)
	value, err := ts.Lookup(index)
	if err != nil {
		t.Fatal(err)
	}
	if value != expected {
		t.Errorf("expected %v but got %v", expected, value)
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
	record, err := ts.Lookup(testtime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if record != 1.5 {
		t.Errorf("unexpected record value %#v", record)
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
	for mday, expectedValue := range expected {
		index := time.Date(2018, 2, mday, 0, 0, 0, 0, time.UTC)
		value, err := ts.Lookup(index)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		if value != expectedValue {
			t.Errorf("for mday %v expected %v but got %v",
				mday, expectedValue, value)
		}
	}
}
