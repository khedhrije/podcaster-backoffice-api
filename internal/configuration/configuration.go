// Package configuration handles the loading and storing of configuration settings for the application.
// It utilizes the viper library to load settings from environment variables.
package configuration

import (
	"github.com/spf13/viper"
	"time"
)

// Config is a global variable that holds the application's configuration settings.
var Config *AppConfig

// init initializes the configuration by loading settings from the environment.
func init() {
	Config = loadFromEnv()
}

// AppConfig defines the structure of the application's configuration settings.
// It includes server details, database, token, cache, and email configurations.
type AppConfig struct {
	Name           string         // Name of the application
	Env            string         // Environment (e.g., development, production)
	HostAddress    string         // Server host address
	HostPort       int            // Server port number
	DocsAddress    string         // Address for API documentation
	DatabaseConfig DatabaseConfig // Configuration settings for the database
	CacheConfig    CacheConfig    // Configuration settings for caching
}

// DatabaseConfig defines the configuration settings for the database connection.
type DatabaseConfig struct {
	Name                   string        // Database name
	Driver                 string        // Database driver type
	DSN                    string        // Data source name for the database
	MaxConnections         int           // Maximum number of open connections to the database
	MaxIdleConnections     int           // Maximum number of idle connections to the database
	MaxLifetimeConnections time.Duration // Maximum amount of time a connection may be reused
}

// CacheConfig defines the configuration for cache connections.
type CacheConfig struct {
	DSN string // Data source name for the cache
}

// loadFromEnv loads configuration settings from environment variables and returns an AppConfig instance.
// It uses viper to handle the environment variables and sets defaults if specific configurations are not provided.
func loadFromEnv() *AppConfig {
	viper.AutomaticEnv()
	return &AppConfig{
		Name:        viper.GetString("APP_NAME"),
		Env:         viper.GetString("APP_ENV"),
		HostAddress: viper.GetString("APP_HOST_ADDRESS"),
		HostPort:    viper.GetInt("APP_HOST_PORT"),
		DocsAddress: viper.GetString("APP_DOCS_HOST_ADDRESS"),
		DatabaseConfig: DatabaseConfig{
			Driver:                 viper.GetString("MYSQL_DRIVER"),
			Name:                   viper.GetString("MYSQL_NAME"),
			DSN:                    viper.GetString("MYSQL_DSN"),
			MaxConnections:         viper.GetInt("MYSQL_MAX_CONNECTIONS"),
			MaxIdleConnections:     viper.GetInt("MYSQL_MAX_IDLE_CONNECTIONS"),
			MaxLifetimeConnections: viper.GetDuration("MYSQL_MAX_CONNECTION_LIFETIME"),
		},
		CacheConfig: CacheConfig{
			DSN: viper.GetString("REDIS_DSN"),
		},
	}
}
