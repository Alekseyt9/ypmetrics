package config

import (
	"encoding/json"
	goflag "flag"
	"os"

	"github.com/Alekseyt9/ypmetrics/internal/common/configh"
	flag "github.com/spf13/pflag"
)

var (
	addressDef        = "localhost:8080"
	reportIntervalDef = 10
	pollIntervalDef   = 2
	rateLimitDef      = 5
)

// Config holds the configuration settings for the client.
type Config struct {
	HashKey        *string // Key for hashing
	Address        *string `json:"address"`         // Server address
	ReportInterval *int    `json:"report_interval"` // Interval for reporting metrics
	PollInterval   *int    `json:"poll_interval"`   // Interval for polling metrics
	RateLimit      *int    // Rate limit for sending metrics
	CryptoKeyFile  *string `json:"crypto_key"` // Key for RSA cypering
	GRPCAddress    *string // GRPC server address
}

func Get() (*Config, error) {
	cfg := &Config{}
	fillConfig(cfg)
	return cfg, nil
}

func fillConfig(cfg *Config) {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	hkeyCmd := flag.StringP("key", "k", "", "key for SHA256 signing")
	addressCmd := flag.StringP("address", "a", "", "Address and port to run server")
	reportIntervalCmd := flag.IntP("reportInterval", "r", 10, "frequency of sending metrics to the server")
	pollIntervalCmd := flag.IntP("pollInterval", "p", 2, "frequency of polling metrics from the runtime package")
	rateLimitCmd := flag.IntP("rateLimit", "l", 5, "upper limit on the number of outgoing requests")
	ckeyCmd := flag.StringP("crypto-key", "z", "", "key for RSA cypering")
	configCmd := flag.StringP("config", "c", "", "config path")
	grpcCmd := flag.StringP("grpc", "g", "", "grpc address")
	flag.Parse()

	if !flag.CommandLine.Changed("key") {
		hkeyCmd = nil
	}
	if !flag.CommandLine.Changed("address") {
		addressCmd = nil
	}
	if !flag.CommandLine.Changed("reportInterval") {
		reportIntervalCmd = nil
	}
	if !flag.CommandLine.Changed("pollInterval") {
		pollIntervalCmd = nil
	}
	if !flag.CommandLine.Changed("rateLimit") {
		rateLimitCmd = nil
	}
	if !flag.CommandLine.Changed("crypto-key") {
		ckeyCmd = nil
	}
	if !flag.CommandLine.Changed("config") {
		configCmd = nil
	}
	if !flag.CommandLine.Changed("grpc") {
		grpcCmd = nil
	}

	hkeyEnv := configh.GetEnvString("KEY")
	addressEnv := configh.GetEnvString("ADDRESS")
	reportIntervalEnv := configh.GetEnvInt("REPORT_INTERVAL")
	pollIntervalEnv := configh.GetEnvInt("POLL_INTERVAL")
	rateLimitEnv := configh.GetEnvInt("RATE_LIMIT")
	ckeyEnv := configh.GetEnvString("CRYPTO_KEY")
	configEnv := configh.GetEnvString("CONFIG")
	grpcEnv := configh.GetEnvString("GRPC_ADDRESS")

	configPar := configh.CombineParams(nil, configCmd, configEnv)
	if configPar != nil {
		setFromFile(cfg, *configPar)
	}

	cfg.HashKey = configh.CombineParams(nil, cfg.HashKey, hkeyCmd, hkeyEnv)
	cfg.Address = configh.CombineParams(&addressDef, cfg.Address, addressCmd, addressEnv)
	cfg.ReportInterval = configh.CombineParams(&reportIntervalDef, cfg.ReportInterval, reportIntervalCmd, reportIntervalEnv)
	cfg.PollInterval = configh.CombineParams(&pollIntervalDef, cfg.PollInterval, pollIntervalCmd, pollIntervalEnv)
	cfg.RateLimit = configh.CombineParams(&rateLimitDef, cfg.RateLimit, rateLimitCmd, rateLimitEnv)
	cfg.CryptoKeyFile = configh.CombineParams(nil, cfg.CryptoKeyFile, ckeyCmd, ckeyEnv)
	cfg.GRPCAddress = configh.CombineParams(nil, cfg.GRPCAddress, grpcCmd, grpcEnv)
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
