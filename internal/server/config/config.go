package config

import (
	"encoding/json"
	"flag"
	"os"
)

// Config holds the configuration settings for the server.
type Config struct {
	Address         *string `env:"ADDRESS" json:"address"`               // Server address
	FileStoragePath *string `env:"FILE_STORAGE_PATH" json:"store_file"`  // Path to the file for storing data
	DataBaseDSN     *string `env:"DATABASE_DSN" json:"database_dsn"`     // Database connection string
	HashKey         *string `env:"KEY"`                                  // Key for SHA256 signing
	Restore         *bool   `env:"RESTORE" json:"restore"`               // Flag to restore data from file on startup
	StoreInterval   *int    `env:"STORE_INTERVAL" json:"store_interval"` // Interval for storing data to file
	CryptoKeyFile   *string `env:"CRYPTO_KEY" json:"crypto_key"`         // Key for RSA cypering
}

func Get() (*Config, error) {
	cfg := &Config{}

	err := SetFromFile(cfg)
	if err != nil {
		return nil, err
	}
	SetFromFlags(cfg)
	SetFromEnv(cfg)

	return cfg, nil
}

func SetFromFile(cfg *Config) error {
	path := getConfigPath()
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

func getConfigPath() string {
	configFlag := flag.String("config", "", "Path to the config file")
	flag.StringVar(configFlag, "c", "", "Path to the config file (shorthand)")
	flag.Parse()

	configEnv := os.Getenv("CONFIG")
	if configEnv != "" {
		configFlag = &configEnv
	}
	return *configFlag
}
