package config

import (
	"flag"
	"os"
)

const (
	envServerAddr      = "SERVER_ADDRESS"
	envBaseURL         = "BASE_URL"
	envFileStoragePath = "FILE_STORAGE_PATH"

	defaultServerAddr      = "localhost:8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = "links.json"
)

var (
	addr            *string
	baseURL         *string
	fileStoragePath *string
	instance        *Config = nil
)

type Config struct {
	Addr            string
	BaseURL         string
	FileStoragePath string
}

// GetConfig returns instance of Config
func GetConfig() *Config {
	if instance != nil {
		return instance
	}

	addr = flag.String("a", "", "Server address and port (default localhost:8080)")
	baseURL = flag.String("b", "", "Server url (http://localhost:8080)")
	fileStoragePath = flag.String("f", "", "Path to file for links (default links.json)")
	flag.Parse()

	instance = &Config{
		Addr:            getAddr(),
		BaseURL:         getBaseURL(),
		FileStoragePath: getFileStoragePath(),
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
