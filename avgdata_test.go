package main

import (
	"testing"
	"time"
)

func TestNewAvgData(t *testing.T) {
	avgdata := NewAvgData()
	if avgdata.records == nil {
		t.Errorf("incomplete initialization")
	}
}

func TestAvgDataRead(t *testing.T) {
	avgdata := NewAvgData()
	avgdata.Read("testdata/simple.csv")
	if len(avgdata.records) != 2 {
		t.Errorf("unexpected series length")
		return
	}

	testtime := time.Date(2018, 1, 1, 10, 0, 0, 0, time.UTC)
	record := avgdata.records[testtime]
	if record != 1.0 {
		t.Errorf("unexpected record value %r", record)
	}
}

func TestAvgDataResample(t *testing.T) {
	avgdata := NewAvgData()
	avgdata.Read("testdata/simple.csv")
	avgdata.Resample()
	testtime := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	record := avgdata.records[testtime]
	if record != 1.5 {
		t.Errorf("unexpected record value %r", record)
	}
}
