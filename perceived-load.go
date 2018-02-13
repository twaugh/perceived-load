package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strconv"
	"time"
)

func list(format string, items []float64) string {
	var result string
	for index, item := range items {
		if index != 0 {
			result += ", "
		}
		result += fmt.Sprintf(format, item)
	}
	return result
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
		ts.Add(time.Now(), ts.records[len(ts.records) - 1].Datum)
	}

	ts.Resample(24 * time.Hour)
	ts.Interpolate()
	var today = time.Now().Truncate(24 * time.Hour)
	days := []float64{1, 5, 15}
	avgs := make([]float64, len(days))
	var since *TimeSeries
	for index, days := range days {
		d := time.Duration((days - 1) * 24) * time.Hour
		start := today.Add(-d)
		since = ts.Since(start)
		var sum float64
		for _, record := range since.records {
			sum += record.Datum
		}
		avgs[index] = sum / float64(len(since.records))
	}

	fmt.Printf("Perceived task load average (%v days): %v\n",
		list("%.0f", days), list("%.1f", avgs))
	fmt.Println("Optimum is 1.0; higher values mean delayed tasks")
}
