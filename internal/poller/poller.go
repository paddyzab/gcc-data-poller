package poller

import (
	"context"
	"log"
	"time"

	"github.com/paddyzab/gcc-data-poller/internal/config"
	"github.com/paddyzab/gcc-data-poller/internal/models"
	"github.com/paddyzab/gcc-data-poller/internal/processor"
)

// DataPoller is responsible for fetching metrics from Google Cloud APIs.
type DataPoller struct {
	cfg       *config.Config
	processor *processor.SignalProcessor
}

func NewDataPoller(cfg *config.Config, proc *processor.SignalProcessor) *DataPoller {
	return &DataPoller{
		cfg:       cfg,
		processor: proc,
	}
}

// Start begins polling Firebase Cloud Messaging and Identity Toolkit metrics.
func (p *DataPoller) Start(ctx context.Context) error {
	// Parse interval, default to 1 minute
	interval, err := time.ParseDuration(p.cfg.MetricsInterval)
	if err != nil {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p.pollFCM(ctx)
			p.pollIdentityToolkit(ctx)
		}
	}
}

func (p *DataPoller) pollFCM(ctx context.Context) {
	// TODO: Implement actual GCP Monitoring API call for FCM
	// e.g. fetching "firebasecloudmessaging.googleapis.com/..."
	log.Println("Polling FCM metrics...")

	// Dummy metric
	metric := models.RawMetric{
		API:       "FCM",
		Timestamp: time.Now(),
		LatencyMs: 150.0,
		ErrorRate: 0.05, // 5%
	}
	p.processor.ProcessMetric(metric)
}

func (p *DataPoller) pollIdentityToolkit(ctx context.Context) {
	// TODO: Implement actual GCP Monitoring API call for Identity Toolkit
	log.Println("Polling Identity Toolkit metrics...")

	// Dummy metric
	metric := models.RawMetric{
		API:       "IdentityToolkit",
		Timestamp: time.Now(),
		LatencyMs: 45.0,
		ErrorRate: 0.01, // 1%
	}
	p.processor.ProcessMetric(metric)
}
