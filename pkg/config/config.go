package config

import (
	"flag"
	"os"
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
	if os.Getenv("SERVER_ADDRESS") != "" {
		return os.Getenv("SERVER_ADDRESS")
	}
	return "localhost:8080"
}

func getBaseURL() string {
	if *baseURL != "" {
		return *baseURL
	}
	if os.Getenv("BASE_URL") != "" {
		return os.Getenv("BASE_URL")
	}
	return "http://localhost:8080"
}

func getFileStoragePath() string {
	if *fileStoragePath != "" {
		return *fileStoragePath
	}
	if os.Getenv("FILE_STORAGE_PATH") != "" {
		return os.Getenv("FILE_STORAGE_PATH")
	}
	return "links.json"
}
