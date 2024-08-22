package config

import (
	"encoding/json"
	"os"
)

// Config holds the configuration settings for the client.
type Config struct {
	HashKey        *string `env:"KEY"`                                    // Key for hashing
	Address        *string `env:"ADDRESS" json:"address"`                 // Server address
	ReportInterval *int    `env:"REPORT_INTERVAL" json:"report_interval"` // Interval for reporting metrics
	PollInterval   *int    `env:"POLL_INTERVAL" json:"poll_interval"`     // Interval for polling metrics
	RateLimit      *int    `env:"RATE_LIMIT"`                             // Rate limit for sending metrics
	CryptoKeyFile  *string `env:"CRYPTO_KEY" json:"crypto_key"`           // Key for RSA cypering
	ConfigFile     *string `env:"CONFIG"`
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
		if cfg.ReportInterval == nil {
			cfg.ReportInterval = c.ReportInterval
		}
		if cfg.PollInterval == nil {
			cfg.PollInterval = c.PollInterval
		}
		if cfg.CryptoKeyFile == nil {
			cfg.CryptoKeyFile = c.CryptoKeyFile
		}
	}
	return nil
}
