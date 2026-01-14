package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Nostr NostrConfig `yaml:"nostr"`
	Rakuten RakutenConfig `yaml:"rakuten"`
	Database DatabaseConfig `yaml:"database"`
	Bot BotConfig `yaml:"bot"`
}

// NostrConfig holds Nostr client settings
type NostrConfig struct {
	SecretKey string   `yaml:"secret_key"`
	Relays    []string `yaml:"relays"`
}

// RakutenConfig holds Rakuten API settings
type RakutenConfig struct {
	ApplicationID string `yaml:"application_id"`
}

// DatabaseConfig holds database settings
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// BotConfig holds bot settings
type BotConfig struct {
	CheckInterval         string `yaml:"check_interval"`
	AnnounceNewReleases   bool   `yaml:"announce_new_releases"`
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.Database.Path == "" {
		cfg.Database.Path = "data/komikan.db"
	}

	return &cfg, nil
}

// LoadFromEnv loads config values from environment variables
// These override values from the config file
func (c *Config) LoadFromEnv() {
	if key := os.Getenv("NOSTR_SECRET_KEY"); key != "" {
		c.Nostr.SecretKey = key
	}
	if appID := os.Getenv("RAKUTEN_APP_ID"); appID != "" {
		c.Rakuten.ApplicationID = appID
	}
}
