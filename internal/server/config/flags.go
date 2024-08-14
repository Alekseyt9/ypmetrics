package config

import (
	goflag "flag"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

// ParseFlags parses command-line flags and sets the corresponding fields in the given Config.
// Parameters:
//   - cfg: the configuration structure to populate with parsed flag values
func ParseFlags(cfg *Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	address := flag.StringP("address", "a", "localhost:8080", "Address and port to run server")
	storeInterval := flag.IntP("store-interval", "i", 300,
		"time interval in seconds, based on which the current state of the server is displayed on disk")
	fileStoragePath := flag.StringP("file-storage-path", "f", "/tmp/metrics-db.json",
		"full name of the file where the current values are saved")
	restore := flag.BoolP("restore", "r", true,
		"load previously saved values from the file when the server starts")
	database := flag.StringP("database", "d", "", "Database connection string")
	key := flag.StringP("key", "k", "", "key for SHA256 signing")
	ckey := flag.StringP("-crypto-key", "z", "", "key for RSA cypering")
	cfgFile := flag.StringP("-config", "c", "", "config file")
	flag.Parse()

	cfg.Address = *address
	cfg.StoreInterval = *storeInterval
	cfg.FileStoragePath = *fileStoragePath
	cfg.Restore = *restore
	cfg.DataBaseDSN = *database
	cfg.HashKey = *key
	cfg.CryptoKeyFile = *ckey
	cfg.ConfigFile = *cfgFile
}

// SetEnv parses environment variables and sets the corresponding fields in the given Config.
// Parameters:
//   - cfg: the configuration structure to populate with parsed environment variable values
func SetEnv(cfg *Config) {
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}
