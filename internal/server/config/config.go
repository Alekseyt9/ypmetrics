package config

// Config holds the configuration settings for the server.
type Config struct {
	Address         string `env:"ADDRESS"`           // Server address
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // Path to the file for storing data
	DataBaseDSN     string `env:"DATABASE_DSN"`      // Database connection string
	HashKey         string `env:"KEY"`               // Key for SHA256 signing
	Restore         bool   `env:"RESTORE"`           // Flag to restore data from file on startup
	StoreInterval   int    `env:"STORE_INTERVAL"`    // Interval for storing data to file
	CryptoKeyFile   string `env:"CRYPTO_KEY"`        // Key for RSA cypering
}
