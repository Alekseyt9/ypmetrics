package config

import (
	goflag "flag"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

// ParseFlags parses command-line flags and sets the corresponding fields in the given Config.
// Parameters:
//   - cfg: the configuration structure to populate with parsed flag values
func SetFromFlags(cfg *Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	address := flag.StringP("address", "a", "localhost:8080", "Address and port to run server")
	storeInterval := flag.IntP("store-interval", "i", 300,
		"time interval in seconds, based on which the current state of the server is displayed on disk")
	fileStoragePath := flag.StringP("file-storage-path", "f", "metrics-db.json",
		"full name of the file where the current values are saved")
	restore := flag.BoolP("restore", "r", true,
		"load previously saved values from the file when the server starts")
	database := flag.StringP("database", "d", "", "Database connection string")
	key := flag.StringP("key", "k", "", "key for SHA256 signing")
	ckey := flag.StringP("-crypto-key", "z", "", "key for RSA cypering")
	flag.Parse()

	if flag.CommandLine.Changed("address") || cfg.Address == nil {
		cfg.Address = address
	}
	if flag.CommandLine.Changed("store-interval") || cfg.StoreInterval == nil {
		cfg.StoreInterval = storeInterval
	}
	if flag.CommandLine.Changed("file-storage-path") || cfg.FileStoragePath == nil {
		cfg.FileStoragePath = fileStoragePath
	}
	if flag.CommandLine.Changed("restore") || cfg.Restore == nil {
		cfg.Restore = restore
	}
	if flag.CommandLine.Changed("database") {
		cfg.DataBaseDSN = database
	}
	if flag.CommandLine.Changed("key") {
		cfg.HashKey = key
	}
	if flag.CommandLine.Changed("crypto-key") {
		cfg.CryptoKeyFile = ckey
	}
}

// SetEnv parses environment variables and sets the corresponding fields in the given Config.
// Parameters:
//   - cfg: the configuration structure to populate with parsed environment variable values
func SetFromEnv(cfg *Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
