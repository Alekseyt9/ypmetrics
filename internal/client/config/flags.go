package config

import (
	goflag "flag"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

// ParseFlags parses command-line flags and sets the corresponding fields in the given Config.
// Parameters:
//   - cfg: the configuration structure to populate with parsed flag values
func ParseFlags(cfg *Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	address := flag.StringP("address", "a", "localhost:8080", "address and port to run server")
	reportInterval := flag.IntP("reportInterval", "r", 10, "frequency of sending metrics to the server")
	pollInterval := flag.IntP("pollInterval", "p", 2, "frequency of polling metrics from the runtime package")
	rateLimit := flag.IntP("rateLimit", "l", 5, "upper limit on the number of outgoing requests")
	key := flag.StringP("key", "k", "", "key for SHA256 signing")
	ckey := flag.StringP("-crypto-key", "z", "", "key for RSA cypering")

	flag.Parse()

	if flag.CommandLine.Changed("address") || cfg.Address == nil {
		cfg.Address = address
	}
	if flag.CommandLine.Changed("reportInterval") || cfg.ReportInterval == nil {
		cfg.ReportInterval = reportInterval
	}
	if flag.CommandLine.Changed("pollInterval") || cfg.PollInterval == nil {
		cfg.PollInterval = pollInterval
	}
	if flag.CommandLine.Changed("rateLimit") || cfg.RateLimit == nil {
		cfg.RateLimit = rateLimit
	}
	cfg.HashKey = key
	cfg.CryptoKeyFile = ckey
}

// SetEnv parses environment variables and sets the corresponding fields in the given Config.
// Parameters:
//   - cfg: the configuration structure to populate with parsed environment variable values
func SetEnv(cfg *Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
