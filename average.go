package main

import (
	"sync"
	"time"
)

func average(ts *TimeSeries, today *time.Time, lookback int) float64 {
	d := time.Duration(lookback*24) * time.Hour
	start := today.Add(-d)
	since := ts.Since(start)
	var sum float64
	for _, record := range since.records {
		sum += record.Datum
	}
	return sum / float64(len(since.records))
}

func averages(ts *TimeSeries, today *time.Time, days ...int) []float64 {
	ts.Resample(24 * time.Hour)
	ts.Interpolate()
	today.Truncate(24 * time.Hour)
	avgs := make([]float64, len(days))
	var wg sync.WaitGroup
	wg.Add(3)
	for index, lookback := range days {
		go func(index, lookback int) {
			avgs[index] = average(ts, today, lookback)
			wg.Done()
		}(index, lookback)
	}
	wg.Wait()
	return avgs
}
