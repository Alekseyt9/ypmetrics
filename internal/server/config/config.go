package config

import (
	"encoding/json"
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
	ConfigFile      *string `env:"CONFIG"`
}

func Get() (*Config, error) {
	cfg := &Config{}
	ParseFlags(cfg)
	SetEnv(cfg)
	err := MergeConfigFromFile(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func FromFile(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func MergeConfigFromFile(cfg *Config) error {
	if cfg.ConfigFile != nil {
		c, err := FromFile(*cfg.ConfigFile)
		if err != nil {
			return err
		}
		if cfg.Address == nil {
			cfg.Address = c.Address
		}
		if cfg.Restore == nil {
			cfg.Restore = c.Restore
		}
		if cfg.StoreInterval == nil {
			cfg.StoreInterval = c.StoreInterval
		}
		if cfg.FileStoragePath == nil {
			cfg.FileStoragePath = c.FileStoragePath
		}
		if cfg.DataBaseDSN == nil {
			cfg.DataBaseDSN = c.DataBaseDSN
		}
		if cfg.CryptoKeyFile == nil {
			cfg.CryptoKeyFile = c.CryptoKeyFile
		}
	}
	return nil
}
