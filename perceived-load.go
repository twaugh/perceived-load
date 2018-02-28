package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
)

func getTimeSeries() *TimeSeries {
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

	return ts
}

func main() {
	ts := getTimeSeries()
	days := []int{1, 5, 15}
	today := time.Now()
	avgs := averages(ts, &today, days...)

	const separator = ", "
	var dayList string
	for _, lookback := range days {
		dayList += fmt.Sprintf("%d%s", lookback, separator)
	}
	dayList = dayList[:len(dayList)-len(separator)]

	var avgList string
	for _, avg := range avgs {
		avgList += fmt.Sprintf("%.1f%s", avg, separator)
	}
	avgList = avgList[:len(avgList)-len(separator)]

	fmt.Printf("Perceived task load average (%s days): %s\n",
		dayList, avgList)
	fmt.Println("Optimum is 1.0; higher values mean delayed tasks")
}
