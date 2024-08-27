package config

import (
	"encoding/json"
	goflag "flag"
	"os"

	"github.com/Alekseyt9/ypmetrics/internal/common/configh"
	flag "github.com/spf13/pflag"
)

var (
	addressDef         = "localhost:8080"
	fileStoragePathDef = "metrics-db.json"
	restoreDef         = true
	storeIntervalDef   = 300
)

// Config holds the configuration settings for the server.
type Config struct {
	Address         *string `json:"address"`      // Server address
	FileStoragePath *string `json:"store_file"`   // Path to the file for storing data
	DataBaseDSN     *string `json:"database_dsn"` // Database connection string
	HashKey         *string // Key for SHA256 signing
	Restore         *bool   `json:"restore"`        // Flag to restore data from file on startup
	StoreInterval   *int    `json:"store_interval"` // Interval for storing data to file
	CryptoKeyFile   *string `json:"crypto_key"`     // Key for RSA cypering
	TrustedSubnet   *string `json:"trusted_subnet"` // Trusted subnet
}

func Get() (*Config, error) {
	cfg := &Config{}
	fillConfig(cfg)
	return cfg, nil
}

func fillConfig(cfg *Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	addressCmd := flag.StringP("address", "a", "", "Address and port to run server")
	fileStoragePathCmd := flag.StringP("file-storage-path", "f", "",
		"full name of the file where the current values are saved")
	databaseCmd := flag.StringP("database", "d", "", "Database connection string")
	hkeyCmd := flag.StringP("key", "k", "", "key for SHA256 signing")
	restoreCmd := flag.BoolP("restore", "r", false,
		"load previously saved values from the file when the server starts")
	storeIntervalCmd := flag.IntP("store-interval", "i", 0,
		"time interval in seconds, based on which the current state of the server is displayed on disk")
	ckeyCmd := flag.StringP("crypto-key", "z", "", "key for RSA cypering")
	configCmd := flag.StringP("config", "c", "", "config path")
	subnetCmd := flag.StringP("trusted-subnet", "t", "", "trusted subnet")
	flag.Parse()

	if !flag.CommandLine.Changed("address") {
		addressCmd = nil
	}
	if !flag.CommandLine.Changed("file-storage-path") {
		fileStoragePathCmd = nil
	}
	if !flag.CommandLine.Changed("database") {
		databaseCmd = nil
	}
	if !flag.CommandLine.Changed("key") {
		hkeyCmd = nil
	}
	if !flag.CommandLine.Changed("restore") {
		restoreCmd = nil
	}
	if !flag.CommandLine.Changed("store-interval") {
		storeIntervalCmd = nil
	}
	if !flag.CommandLine.Changed("crypto-key") {
		ckeyCmd = nil
	}
	if !flag.CommandLine.Changed("config") {
		configCmd = nil
	}
	if !flag.CommandLine.Changed("trusted-subnet") {
		subnetCmd = nil
	}

	addressEnv := configh.GetEnvString("ADDRESS")
	fileStoragePathEnv := configh.GetEnvString("FILE_STORAGE_PATH")
	databaseEnv := configh.GetEnvString("DATABASE_DSN")
	hkeyEnv := configh.GetEnvString("KEY")
	restoreEnv := configh.GetEnvBool("RESTORE")
	storeIntervalEnv := configh.GetEnvInt("STORE_INTERVAL")
	ckeyEnv := configh.GetEnvString("CRYPTO_KEY")
	configEnv := configh.GetEnvString("CONFIG")
	subnetEnv := configh.GetEnvString("TRUSTED_SUBNET")

	configPar := configh.CombineParams(nil, configCmd, configEnv)
	if configPar != nil {
		setFromFile(cfg, *configPar)
	}

	cfg.Address = configh.CombineParams(&addressDef, cfg.Address, addressCmd, addressEnv)
	cfg.FileStoragePath = configh.CombineParams(&fileStoragePathDef, cfg.FileStoragePath, fileStoragePathCmd, fileStoragePathEnv)
	cfg.DataBaseDSN = configh.CombineParams(nil, cfg.DataBaseDSN, databaseCmd, databaseEnv)
	cfg.HashKey = configh.CombineParams(nil, cfg.HashKey, hkeyCmd, hkeyEnv)
	cfg.Restore = configh.CombineParams(&restoreDef, cfg.Restore, restoreCmd, restoreEnv)
	cfg.StoreInterval = configh.CombineParams(&storeIntervalDef, cfg.StoreInterval, storeIntervalCmd, storeIntervalEnv)
	cfg.CryptoKeyFile = configh.CombineParams(nil, cfg.CryptoKeyFile, ckeyCmd, ckeyEnv)
	cfg.TrustedSubnet = configh.CombineParams(nil, cfg.TrustedSubnet, subnetCmd, subnetEnv)
}

func setFromFile(cfg *Config, path string) error {
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}
