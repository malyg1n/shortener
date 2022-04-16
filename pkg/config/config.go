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
	TrustedSubnet   string
	APIType         string
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
	envConfigFilePath  = "config"
	envTrustedSubnet   = "trusted_subnet"
	envAPIType         = "api_type"

	defaultServerAddr      = ":8080"
	defaultBaseURL         = "http://localhost:8080"
	defaultFileStoragePath = "links.json"
	defaultSecretKey       = "secret-key"
	defaultDatabaseDSN     = "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"
	defaultSSLCert         = "certs/cert.crt"
	defaultSSLKey          = "certs/key.key"
	defaultConfigFilePath  = ""
	defaultTrustedSubnet   = ""
	defaultAPIType         = "rest"
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
	pflag.StringP(envSSLCert, "r", defaultSSLCert, "ssl cert file")
	pflag.StringP(envSSLKey, "p", defaultSSLKey, "ssl private key file")
	pflag.StringP(envConfigFilePath, "c", defaultConfigFilePath, "path to config")
	pflag.StringP(envTrustedSubnet, "t", defaultTrustedSubnet, "trusted subnet")
	pflag.StringP(envAPIType, "g", defaultAPIType, "type of api")
}

// GetConfig returns instance of Config
func GetConfig() *Config {
	if instance != nil {
		return instance
	}

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	if s := viper.GetString(envConfigFilePath); s != "" {
		viper.SetConfigFile(s)
		viper.SetConfigType("json")
		viper.ReadInConfig()
	}

	instance = &Config{
		Addr:            viper.GetString(envServerAddr),
		BaseURL:         viper.GetString(envBaseURL),
		FileStoragePath: viper.GetString(envFileStoragePath),
		SecretKey:       viper.GetString(envSecretKey),
		DatabaseDSN:     viper.GetString(envDatabaseDSN),
		EnableHTTPS:     viper.GetBool(envEnableHTTPS),
		SSLCert:         viper.GetString(envSSLCert),
		SSLPrivateKey:   viper.GetString(envSSLKey),
		TrustedSubnet:   viper.GetString(envTrustedSubnet),
		APIType:         viper.GetString(envAPIType),
	}

	return instance
}
