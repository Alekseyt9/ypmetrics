package config

import (
	"encoding/json"
	"flag"
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
