package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strconv"
	"time"
)

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
	}

	ts.Resample(24 * time.Hour)
	ts.Interpolate()
	var today = time.Now().Truncate(24 * time.Hour)
	days := []int{1, 5, 15}
	avgs := make([]float64, len(days))
	var since *TimeSeries
	for index, days := range days {
		d := time.Duration((days - 1) * 24) * time.Hour
		start := today.Add(-d)
		since = ts.Since(start)
		for _, record := range since.records {
			log.Printf("%v: %v\n", record.timestamp.String(), record.datum)
		}
		log.Printf("--")
		var sum float64
		for _, record := range since.records {
			sum += record.datum
		}
		avgs[index] = sum / float64(len(since.records))
	}

	fmt.Printf("Perceived task load average (%v days): %v\n", days, avgs)
}
