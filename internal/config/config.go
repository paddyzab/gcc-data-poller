package config

import "os"

// Config holds all configuration for the service
type Config struct {
	ProjectID   string
	GRPCPort    string
	MetricsInterval string // e.g. "60s"
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	return &Config{
		ProjectID:   os.Getenv("GOOGLE_CLOUD_PROJECT"),
		GRPCPort:    port,
		MetricsInterval: "60s",
	}, nil
}
