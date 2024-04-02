package main

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
)

var flagAddr *string = flag.StringP("address", "a", "localhost:8080", "address and port to run server")

func parseFlags() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
