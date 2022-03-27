package config

import (
	"flag"
	"os"
)

const (
	envServerAddr      = "SERVER_ADDRESS"
	envBaseURL         = "BASE_URL"
	envFileStoragePath = "FILE_STORAGE_PATH"
	envSecretKey       = "APP_KEY"
	envDatabaseDSN     = "DATABASE_DSN"
	envEnableHTTPS     = "ENABLE_HTTPS"

	defaultServerAddr      = ":8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = "links.json"
	defaultSecretKey       = "secret-key"
	defaultDatabaseDSN     = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
)

var (
	addr            *string
	baseURL         *string
	fileStoragePath *string
	databaseDSN     *string
	instance        *Config = nil
	enableHTTPS     *bool
)

type Config struct {
	Addr            string
	BaseURL         string
	FileStoragePath string
	SecretKey       string
	DatabaseDSN     string
	EnableHTTPS     bool
}

func init() {
	addr = flag.String("a", "", "Server address and port (default localhost:8080)")
	baseURL = flag.String("b", "", "Server url (http://localhost:8080)")
	fileStoragePath = flag.String("f", "", "Path to file for links (default links.json)")
	databaseDSN = flag.String("d", "", "Database connection string")
	enableHTTPS = flag.Bool("s", false, "Enable https")
}

// GetConfig returns instance of Config
func GetConfig() *Config {
	if instance != nil {
		return instance
	}

	flag.Parse()
	instance = &Config{
		Addr:            getAddr(),
		BaseURL:         getBaseURL(),
		FileStoragePath: getFileStoragePath(),
		SecretKey:       getEnv(envSecretKey, defaultSecretKey),
		DatabaseDSN:     getDatabaseDSN(),
		EnableHTTPS:     getEnableHTTPS(),
	}

	return instance
}

func getAddr() string {
	if *addr != "" {
		return *addr
	}
	if os.Getenv(envServerAddr) != "" {
		return os.Getenv(envServerAddr)
	}

	return defaultServerAddr
}

func getBaseURL() string {
	if *baseURL != "" {
		return *baseURL
	}
	if os.Getenv(envBaseURL) != "" {
		return os.Getenv(envBaseURL)
	}

	return defaultBaseURL
}

func getFileStoragePath() string {
	if *fileStoragePath != "" {
		return *fileStoragePath
	}
	if os.Getenv(envFileStoragePath) != "" {
		return os.Getenv(envFileStoragePath)
	}

	return defaultFileStoragePath
}

func getDatabaseDSN() string {
	if *databaseDSN != "" {
		return *databaseDSN
	}
	if os.Getenv(envDatabaseDSN) != "" {
		return os.Getenv(envDatabaseDSN)
	}

	return defaultDatabaseDSN
}

func getEnableHTTPS() bool {
	if *enableHTTPS == true {
		return *enableHTTPS
	}

	if os.Getenv(envEnableHTTPS) == "true" || os.Getenv(envEnableHTTPS) == "1" {
		return true
	}

	return false
}

func getEnv(envName, defaultValue string) string {
	if os.Getenv(envName) != "" {
		return os.Getenv(envName)
	}

	return defaultValue
}
