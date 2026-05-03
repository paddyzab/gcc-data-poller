package models

import "time"

// RawMetric represents a single metric polled from GCP
type RawMetric struct {
	API       string
	Timestamp time.Time
	LatencyMs float64
	ErrorRate float64
}

// ActionSignal represents a generated signal sent to downstream services
type ActionSignal struct {
	API       string
	Timestamp time.Time
	Level     SignalLevel
	Message   string
}

type SignalLevel string

const (
	LevelInfo     SignalLevel = "INFO"
	LevelWarning  SignalLevel = "WARNING"
	LevelCritical SignalLevel = "CRITICAL"
)
