package main

import (
	"flag"
	"strings"
)

// GetStringFromCLI returns a string obtained from command line arguments or a flag.
// Example: ./goxkcd -s "hello world"
// Example: ./goxkcd hello world
func GetStringFromCLI() string {
	var stringFromArgs string
	flag.StringVar(&stringFromArgs, "s", "", "A string flag")
	flag.Parse()
	if stringFromArgs != "" {
		return stringFromArgs
	}
	args := flag.Args()
	return strings.Join(args, " ")
}
