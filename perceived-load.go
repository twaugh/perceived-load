package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Duration int

func (d Duration) String() string {
	return fmt.Sprintf("%v", int(d))
}

type LoadAvg float64

func (l LoadAvg) String() string{
	return fmt.Sprintf("%.1f", float64(l))
}

func list(items []fmt.Stringer) string {
	item_strings := make([]string, len(items))
	for index, item := range items {
		item_strings[index] = item.String()
	}
	return strings.Join(item_strings, ", ")
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

	ts.Resample(24 * time.Hour)
	ts.Interpolate()
	var today = time.Now().Truncate(24 * time.Hour)
	lookbacks := []int{1, 5, 15}
	days := make([]fmt.Stringer, len(lookbacks))
	avgs := make([]fmt.Stringer, len(lookbacks))
	var since *TimeSeries
	for index, lookback := range lookbacks {
		d := time.Duration((lookback-1)*24) * time.Hour
		start := today.Add(-d)
		since = ts.Since(start)
		var sum float64
		for _, record := range since.records {
			sum += record.Datum
		}
		days[index] = Duration(lookback)
		avgs[index] = LoadAvg(sum / float64(len(since.records)))
	}

	fmt.Printf("Perceived task load average (%v days): %v\n",
		list(days), list(avgs))
	fmt.Println("Optimum is 1.0; higher values mean delayed tasks")
}
