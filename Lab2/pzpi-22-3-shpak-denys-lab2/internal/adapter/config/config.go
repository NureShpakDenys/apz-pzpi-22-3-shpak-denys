// Package config provides utilities for loading application configuration
// using YAML files and environment variables.
package config // import "wayra/internal/adapter/config"

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config defines the main configuration structure.
// This structure includes storage paths, HTTP server configuration,
// authentication settings, and database credentials.
type Config struct {
	StoragePath   string     `yaml:"storage_path" env-required:"true"`   // Path to the storage directory.
	Http          HttpConfig `yaml:"http"`                               // HTTP server configuration.
	AuthConfig    AuthConfig `yaml:"auth"`                               // Authentication configuration.
	DBPassword    string     `yaml:"db_password" env-required:"true"`    // Database password.
	EncryptionKey string     `yaml:"encryption_key" env-required:"true"` // Encryption key.
}

// HttpConfig defines the HTTP server configuration.
// This structure includes the port and timeout settings for the HTTP server.
type HttpConfig struct {
	Port    int           `yaml:"port"`    // Port number for the HTTP server.
	Timeout time.Duration `yaml:"timeout"` // Timeout duration for the HTTP server.
}

// AuthConfig defines the authentication configuration.
// This structure includes the secret key and token expiration time.
type AuthConfig struct {
	SecretKey   string        `yaml:"secret"`    // Secret key for authentication.
	TokenExpiry time.Duration `yaml:"token_ttl"` // Token expiration duration.
}

// MustLoad loads the configuration file specified by the CONFIG_PATH
// environment variable or the --config flag and panics if any error occurs.
// This function ensures the configuration is properly loaded or terminates the application.
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

// MustLoadPath loads the configuration from the given path and panics if any error occurs.
// This function reads the YAML configuration file and parses it into a Config structure.
func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches the path to the configuration file from the --config flag
// or the CONFIG_PATH environment variable.
// This function prioritizes the command-line flag over the environment variable.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
