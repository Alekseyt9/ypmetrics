package run

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
)

var FlagAddr *string = flag.StringP("address", "a", "localhost:8080", "address and port to run server")

func ParseFlags() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}
