package main

import (
	"github.com/jessevdk/go-flags"
	"fmt"
	"os"
)

func main() {
	var opts struct {
		DB string `long:"db" value-name:"FILE" description:"database file to use"`
	}

	args, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Printf("Unexpected parameters: %v\n", args[1:])
		os.Exit(1)
	}

	ts := NewTimeSeries()
	if err := ts.Read(opts.DB); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ts, args, opts.DB)
}
