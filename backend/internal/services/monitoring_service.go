package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MonitoringService struct {
	QueueDepth     prometheus.GaugeVec
	MessagesTotal  prometheus.CounterVec
	ProcessingTime prometheus.HistogramVec
}

func NewMonitoringService() *MonitoringService {
	return &MonitoringService{
		QueueDepth: *promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "qafka_queue_depth",
				Help: "The current number of messages in each queue",
			},
			[]string{"queue_name", "status"},
		),
		MessagesTotal: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "qafka_messages_total",
				Help: "The total number of messages processed",
			},
			[]string{"queue_name", "status"},
		),
		ProcessingTime: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "qafka_message_processing_duration_seconds",
				Help: "Time spent processing messages",
			},
			[]string{"queue_name"},
		),
	}
}

func (m *MonitoringService) IncrementMessageCounter(queueName, status string) {
	m.MessagesTotal.WithLabelValues(queueName, status).Inc()
}

func (m *MonitoringService) SetQueueDepth(queueName, status string, depth float64) {
	m.QueueDepth.WithLabelValues(queueName, status).Set(depth)
}

func (m *MonitoringService) ObserveProcessingTime(queueName string, duration float64) {
	m.ProcessingTime.WithLabelValues(queueName).Observe(duration)
}
