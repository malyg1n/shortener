package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = "link.json"

	serverAddressVarName   = "served_address"
	baseURLVarName         = "base_url"
	fileStoragePathVarName = "file_storage_path"
)

var (
	instance *Config = nil
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

	processVars()

	instance = &Config{
		Addr:            viper.GetString(serverAddressVarName),
		BaseURL:         viper.GetString(baseURLVarName),
		FileStoragePath: viper.GetString(fileStoragePathVarName),
	}

	return instance
}

func processVars() {
	pflag.String("a", "localhost:8080", "Server address and port")
	pflag.String("b", "http://localhost:8080", "Server URL")
	pflag.String("f", "links.json", "Path to file for links")
	pflag.Parse()

	_ = viper.BindPFlag(serverAddressVarName, pflag.CommandLine.Lookup("a"))
	_ = viper.BindPFlag(baseURLVarName, pflag.CommandLine.Lookup("b"))
	_ = viper.BindPFlag(fileStoragePathVarName, pflag.CommandLine.Lookup("f"))

	_ = viper.BindEnv(serverAddressVarName)
	_ = viper.BindEnv(baseURLVarName)
	_ = viper.BindEnv(fileStoragePathVarName)

	viper.SetDefault(serverAddressVarName, defaultServerAddress)
	viper.SetDefault(baseURLVarName, defaultBaseURL)
	viper.SetDefault(fileStoragePathVarName, defaultFileStoragePath)

	viper.AutomaticEnv()
}
