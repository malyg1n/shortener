package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config for application
type Config struct {
	Addr            string
	BaseURL         string
	FileStoragePath string
	SecretKey       string
	DatabaseDSN     string
	EnableHTTPS     bool
	SSLCert         string
	SSLPrivateKey   string
}

const (
	envServerAddr      = "server_address"
	envBaseURL         = "base_url"
	envFileStoragePath = "file_storage_path"
	envSecretKey       = "app_key"
	envDatabaseDSN     = "database_dsn"
	envEnableHTTPS     = "enable_https"
	envSSLCert         = "ssl_cert"
	envSSLKey          = "ssl_key"

	defaultServerAddr      = ":8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = "links.json"
	defaultSecretKey       = "secret-key"
	defaultDatabaseDSN     = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
	defaultSSLCert         = "certs/cert.crt"
	defaultSSLKey          = "certs/key.key"
)

var instance *Config

func init() {
	viper.AutomaticEnv()
	pflag.StringP(envDatabaseDSN, "d", defaultDatabaseDSN, "database connection string")
	pflag.StringP(envServerAddr, "a", defaultServerAddr, "run address")
	pflag.StringP(envBaseURL, "b", defaultBaseURL, "base url")
	pflag.StringP(envFileStoragePath, "f", defaultFileStoragePath, "base url")
	pflag.StringP(envSecretKey, "k", defaultSecretKey, "secret key for app")
	pflag.BoolP(envEnableHTTPS, "s", false, "use https")
	pflag.StringP(envSSLCert, "c", defaultSSLCert, "ssl cert file")
	pflag.StringP(envSSLKey, "p", defaultSSLKey, "ssl private key file")
}

// GetConfig returns instance of Config
func GetConfig() *Config {
	if instance != nil {
		return instance
	}

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	instance = &Config{
		Addr:            viper.GetString(envServerAddr),
		BaseURL:         viper.GetString(envBaseURL),
		FileStoragePath: viper.GetString(envFileStoragePath),
		SecretKey:       viper.GetString(envSecretKey),
		DatabaseDSN:     viper.GetString(envDatabaseDSN),
		EnableHTTPS:     viper.GetBool(envEnableHTTPS),
		SSLCert:         viper.GetString(envSSLCert),
		SSLPrivateKey:   viper.GetString(envSSLKey),
	}

	return instance
}
