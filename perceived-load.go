package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strconv"
	"time"
)

func averages(ts *TimeSeries, days ...int) []float64 {
	ts.Resample(24 * time.Hour)
	ts.Interpolate()
	var today = time.Now().Truncate(24 * time.Hour)
	avgs := make([]float64, len(days))
	var since *TimeSeries
	for index, days := range days {
		d := time.Duration((days-1)*24) * time.Hour
		start := today.Add(-d)
		since = ts.Since(start)
		var sum float64
		for _, record := range since.records {
			sum += record.Datum
		}
		avgs[index] = sum / float64(len(since.records))
	}

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

	var day_list string
	for index, lookback := range days {
		if index > 0 {
			day_list += ", "
		}
		day_list += fmt.Sprint(lookback)
	}

	var avg_list string
	for index, avg := range avgs {
		if index > 0 {
			avg_list += ", "
		}
		avg_list += fmt.Sprintf("%.1f", avg)
	}

	fmt.Printf("Perceived task load average (%s days): %s\n",
		day_list, avg_list)
	fmt.Println("Optimum is 1.0; higher values mean delayed tasks")
}
