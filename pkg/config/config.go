package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultBaseUrl         = "http://localhost:8080"
	defaultFileStoragePath = "link.json"

	serverAddressVarName   = "served_address"
	baseUrlVarName         = "base_url"
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
		Addr:            viper.GetString("server_address"),
		BaseURL:         viper.GetString("base_url"),
		FileStoragePath: viper.GetString("file_storage_path"),
	}

	return instance
}

func processVars() {
	pflag.String("a", "localhost:8080", "Server address and port")
	pflag.String("b", "http://localhost:8080", "Server url")
	pflag.String("f", "links.json", "Path to file for links")
	pflag.Parse()

	_ = viper.BindPFlag(serverAddressVarName, pflag.CommandLine.Lookup("a"))
	_ = viper.BindPFlag(baseUrlVarName, pflag.CommandLine.Lookup("b"))
	_ = viper.BindPFlag(fileStoragePathVarName, pflag.CommandLine.Lookup("f"))

	_ = viper.BindEnv(serverAddressVarName)
	_ = viper.BindEnv(baseUrlVarName)
	_ = viper.BindEnv(fileStoragePathVarName)

	viper.SetDefault(serverAddressVarName, defaultServerAddress)
	viper.SetDefault(baseUrlVarName, defaultBaseUrl)
	viper.SetDefault(fileStoragePathVarName, defaultFileStoragePath)

	viper.AutomaticEnv()
}
