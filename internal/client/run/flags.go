package run

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
)

var FlagAddr *string = flag.StringP("address", "a", "localhost:8080", "address and port to connect to server")
var FlagReportInterval *int = flag.IntP("reportInterval", "r", 10, "frequency of sending metrics to the server")
var FlagPollInterval *int = flag.IntP("pollInterval", "p", 2, "frequency of polling metrics from the runtime package")

func ParseFlags() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
