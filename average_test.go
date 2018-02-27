package main

import (
	"testing"
	"time"
)

func newAverageTestData(t *testing.T) (ts *TimeSeries, today time.Time) {
	ts = NewTimeSeries()
	if err := ts.Read("testdata/average.csv"); err != nil {
		t.Fatal(err)
	}
	today = time.Date(2018, 1, 15, 10, 0, 0, 0, time.UTC)
	return
}

func TestAverage(t *testing.T) {
	ts, today := newAverageTestData(t)
	for lookback, expect := range map[int]float64{
		1: 15,
		2: 14.5,
		3: 14,
		4: 13.5,
		5: 13,
	} {
		avg := average(ts, &today, lookback)
		if avg != expect {
			t.Errorf("lookback %d: incorrect avg %f (expected %f)",
				lookback, avg, expect)
		}
	}
}

func TestAverages(t *testing.T) {
	ts, today := newAverageTestData(t)
	for days, expect := range map[[3]int]([]float64){
		[3]int{1, 3, 5}: []float64{15, 14, 13},
	} {
		avgs := averages(ts, &today, days[:3]...)
		if len(avgs) != len(expect) {
			t.Errorf("incorrect length: %d (expected %d)",
				len(avgs), len(expect))
			continue
		}
		for index, value := range avgs {
			if value != expect[index] {
				t.Errorf("days %v: incorrect avgs %v (expected %v)",
					days, avgs, expect)
			}
		}
	}
}
