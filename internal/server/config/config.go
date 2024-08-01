package config

// Config holds the configuration settings for the server.
type Config struct {
	Address         string `env:"ADDRESS"`           // Server address
	StoreInterval   int    `env:"STORE_INTERVAL"`    // Interval for storing data to file
	FileStoragePath string `env:"FILE_STORAGE_PATH"` // Path to the file for storing data
	Restore         bool   `env:"RESTORE"`           // Flag to restore data from file on startup
	DataBaseDSN     string `env:"DATABASE_DSN"`      // Database connection string
	HashKey         string `env:"KEY"`               // Key for SHA256 signing
}
