package main

import (
	"flag"
)

// GetStringFromCLI returns a string obtained from a flag.
// Example: ./goxkcd -s "hello world"
func GetStringFromCLI() string {
	var stringFromArgs string
	flag.StringVar(&stringFromArgs, "s", "", "A string flag")
	flag.Parse()
	return stringFromArgs
}
