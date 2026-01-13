package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

// Config represents the complete application configuration
type Config struct {
	Server    ServerConfig    `toml:"server"`
	Endpoints EndpointsConfig `toml:"endpoints"`
	CORS      CORSConfig      `toml:"cors"`
	Frontend  FrontendConfig  `toml:"frontend"`

	// DatabaseURL is loaded from environment variable
	DatabaseURL string `toml:"-"`
}

// ServerConfig holds server-related settings
type ServerConfig struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

// EndpointsConfig holds API endpoint URLs
type EndpointsConfig struct {
	APIBase string `toml:"api_base"`
	WSBase  string `toml:"ws_base"`
}

// CORSConfig holds CORS settings
type CORSConfig struct {
	AllowedOrigins []string `toml:"allowed_origins"`
}

// FrontendConfig holds frontend-related settings
type FrontendConfig struct {
	Port int `toml:"port"`
}

// AppConfig is the global configuration instance
var AppConfig *Config

// Load loads configuration from TOML file and environment variables
// env should be "dev", "staging", or "prod"
func Load(env string) (*Config, error) {
	// Determine config file path
	configPath := getConfigPath(env)

	// Load TOML config
	var cfg Config
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	// Load environment variables from .env file
	envPath := getEnvPath(env)
	if err := godotenv.Load(envPath); err != nil {
		// Not fatal - env vars might be set directly in Docker
		log.Printf("Note: Could not load %s (this is OK in Docker): %v", envPath, err)
	}

	// Load database URL from environment
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	AppConfig = &cfg
	return &cfg, nil
}

// GetDatabaseURL returns the PostgreSQL connection string from environment
func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}

// GetServerAddr returns the server address in host:port format
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

// getConfigPath returns the path to the config file for the given environment
func getConfigPath(env string) string {
	// First check if CONFIG_PATH is set (for Docker)
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	// Try to find config relative to working directory or project root
	filename := fmt.Sprintf("config.%s.toml", env)

	// Check current directory
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	// Check parent directory (for when running from backend/)
	parentPath := filepath.Join("..", filename)
	if _, err := os.Stat(parentPath); err == nil {
		return parentPath
	}

	// Check two levels up (for when running from backend/counter/)
	grandparentPath := filepath.Join("..", "..", filename)
	if _, err := os.Stat(grandparentPath); err == nil {
		return grandparentPath
	}

	// Default to current directory
	return filename
}

// getEnvPath returns the path to the .env file for the given environment
func getEnvPath(env string) string {
	// First check if ENV_PATH is set (for Docker)
	if path := os.Getenv("ENV_PATH"); path != "" {
		return path
	}

	filename := fmt.Sprintf(".env.%s", env)

	// Check current directory
	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	// Check parent directory
	parentPath := filepath.Join("..", filename)
	if _, err := os.Stat(parentPath); err == nil {
		return parentPath
	}

	// Check two levels up
	grandparentPath := filepath.Join("..", "..", filename)
	if _, err := os.Stat(grandparentPath); err == nil {
		return grandparentPath
	}

	return filename
}

// GetEnv returns the current environment from CONFIG_ENV or defaults to "dev"
func GetEnv() string {
	env := os.Getenv("CONFIG_ENV")
	if env == "" {
		return "dev"
	}
	return env
}
