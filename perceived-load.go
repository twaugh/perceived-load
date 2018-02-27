package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
)

func average(ts *TimeSeries, today *time.Time, lookback int) float64 {
	d := time.Duration((lookback-1)*24) * time.Hour
	start := today.Add(-d)
	since := ts.Since(start)
	var sum float64
	for _, record := range since.records {
		sum += record.Datum
	}
	return sum / float64(len(since.records))
}

func averages(ts *TimeSeries, days ...int) []float64 {
	ts.Resample(24 * time.Hour)
	ts.Interpolate()
	var today = time.Now().Truncate(24 * time.Hour)
	avgs := make([]float64, len(days))
	var wg sync.WaitGroup
	wg.Add(3)
	for index, lookback := range days {
		go func(index, lookback int) {
			avgs[index] = average(ts, &today, lookback)
			wg.Done()
		}(index, lookback)
	}
	wg.Wait()
	return avgs
}

func main() {
	var opts struct {
		DB string `long:"db" value-name:"FILE" description:"database file to use"`
	}

	args, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		log.Fatalf("Unexpected parameters: %v\n", args[1:])
	}

	if opts.DB == "" {
		opts.DB = os.ExpandEnv("${HOME}/.config/perceived-load.csv")
	}

	ts := NewTimeSeries()
	if err := ts.Read(opts.DB); err != nil {
		log.Fatal(err)
	}

	if len(args) > 0 {
		value, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		ts.Add(time.Now(), value)
		if err := ts.Write(opts.DB); err != nil {
			log.Fatal(err)
		}
	} else if len(ts.records) > 0 {
		ts.Add(time.Now(), ts.records[len(ts.records)-1].Datum)
	}

	days := []int{1, 5, 15}
	avgs := averages(ts, days...)

	const separator = ", "
	var day_list string
	for _, lookback := range days {
		day_list += fmt.Sprintf("%d%s", lookback, separator)
	}
	day_list = day_list[:len(day_list)-len(separator)]

	var avg_list string
	for _, avg := range avgs {
		avg_list += fmt.Sprintf("%.1f%s", avg, separator)
	}
	avg_list = avg_list[:len(avg_list)-len(separator)]

	fmt.Printf("Perceived task load average (%s days): %s\n",
		day_list, avg_list)
	fmt.Println("Optimum is 1.0; higher values mean delayed tasks")
}
