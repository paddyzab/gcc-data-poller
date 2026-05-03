package processor

import (
	"log"
	"sync"

	"github.com/paddyzab/gcc-data-poller/internal/models"
)

// SignalProcessor processes raw metrics and converts them to actionable signals
type SignalProcessor struct {
	mu          sync.RWMutex
	subscribers []chan models.ActionSignal
}

func NewSignalProcessor() *SignalProcessor {
	return &SignalProcessor{
		subscribers: make([]chan models.ActionSignal, 0),
	}
}

// ProcessMetric evaluates latency and error rates to generate signals
func (sp *SignalProcessor) ProcessMetric(metric models.RawMetric) {
	var level models.SignalLevel
	var msg string

	if metric.ErrorRate > 0.05 || metric.LatencyMs > 500 {
		level = models.LevelCritical
		msg = "High error rate or latency detected"
	} else if metric.ErrorRate > 0.01 || metric.LatencyMs > 200 {
		level = models.LevelWarning
		msg = "Elevated error rate or latency"
	} else {
		// Just info, or skip. For PoC, let's emit INFO.
		level = models.LevelInfo
		msg = "API operating normally"
	}

	signal := models.ActionSignal{
		API:       metric.API,
		Timestamp: metric.Timestamp,
		Level:     level,
		Message:   msg,
	}

	log.Printf("Generated signal: [%s] %s - %s", signal.API, signal.Level, signal.Message)
	sp.broadcast(signal)
}

// Subscribe returns a channel that receives newly generated signals
func (sp *SignalProcessor) Subscribe() <-chan models.ActionSignal {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	ch := make(chan models.ActionSignal, 100)
	sp.subscribers = append(sp.subscribers, ch)
	return ch
}

func (sp *SignalProcessor) broadcast(sig models.ActionSignal) {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	for _, ch := range sp.subscribers {
		select {
		case ch <- sig:
		default:
			// Dropping signal if subscriber is too slow
			log.Println("Subscriber channel full, dropping signal")
		}
	}
}
