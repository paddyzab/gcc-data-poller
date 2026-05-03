package config

import (
	"encoding/json"
	"os"
)

// Config holds all configuration for the service
type Config struct {
	ProjectID       string `json:"project_id"`
	GRPCPort        string `json:"grpc_port"`
	MetricsInterval string `json:"metrics_interval"` // e.g. "60s"
	CredentialsFile string `json:"credentials_file"` // Path to Google credentials JSON
}

// Load loads configuration from .config.json if it exists, otherwise from environment variables
func Load() (*Config, error) {
	// Default config from env vars
	cfg := &Config{
		ProjectID:       os.Getenv("GOOGLE_CLOUD_PROJECT"),
		GRPCPort:        os.Getenv("GRPC_PORT"),
		MetricsInterval: "60s",
		CredentialsFile: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}

	if cfg.GRPCPort == "" {
		cfg.GRPCPort = "50051"
	}

	// Try loading from .config.json
	fileBytes, err := os.ReadFile(".config.json")
	if err == nil {
		// File exists, decode JSON to overwrite default values
		if err := json.Unmarshal(fileBytes, cfg); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		// Error reading file
		return nil, err
	}

	return cfg, nil
}
