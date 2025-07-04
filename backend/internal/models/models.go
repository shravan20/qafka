package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Queue represents a message queue
type Queue struct {
	bun.BaseModel `bun:"table:queues"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Name        string    `bun:"name,notnull,unique" json:"name"`
	Description string    `bun:"description" json:"description"`
	Type        string    `bun:"type,notnull" json:"type"`        // fifo, priority, delay, etc.
	Config      string    `bun:"config,type:jsonb" json:"config"` // JSON configuration
	IsActive    bool      `bun:"is_active,notnull,default:true" json:"is_active"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

// Message represents a message in a queue
type Message struct {
	bun.BaseModel `bun:"table:messages"`

	ID           int64      `bun:"id,pk,autoincrement" json:"id"`
	QueueID      int64      `bun:"queue_id,notnull" json:"queue_id"`
	Queue        *Queue     `bun:"rel:belongs-to,join:queue_id=id" json:"queue,omitempty"`
	Payload      string     `bun:"payload,notnull" json:"payload"`
	Priority     int        `bun:"priority,notnull,default:0" json:"priority"`
	Status       string     `bun:"status,notnull,default:'pending'" json:"status"` // pending, processing, completed, failed
	ScheduledAt  *time.Time `bun:"scheduled_at" json:"scheduled_at,omitempty"`
	ProcessedAt  *time.Time `bun:"processed_at" json:"processed_at,omitempty"`
	FailedAt     *time.Time `bun:"failed_at" json:"failed_at,omitempty"`
	RetryCount   int        `bun:"retry_count,notnull,default:0" json:"retry_count"`
	MaxRetries   int        `bun:"max_retries,notnull,default:3" json:"max_retries"`
	ErrorMessage string     `bun:"error_message" json:"error_message,omitempty"`
	CreatedAt    time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

// Worker represents a queue worker/consumer
type Worker struct {
	bun.BaseModel `bun:"table:workers"`

	ID             int64     `bun:"id,pk,autoincrement" json:"id"`
	Name           string    `bun:"name,notnull" json:"name"`
	QueueID        int64     `bun:"queue_id,notnull" json:"queue_id"`
	Queue          *Queue    `bun:"rel:belongs-to,join:queue_id=id" json:"queue,omitempty"`
	Status         string    `bun:"status,notnull,default:'idle'" json:"status"` // idle, busy, stopped
	LastPing       time.Time `bun:"last_ping,nullzero,notnull,default:current_timestamp" json:"last_ping"`
	ProcessedCount int64     `bun:"processed_count,notnull,default:0" json:"processed_count"`
	FailedCount    int64     `bun:"failed_count,notnull,default:0" json:"failed_count"`
	CreatedAt      time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

// CreateQueueRequest represents the request to create a new queue
type CreateQueueRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Type        string `json:"type" validate:"required"`
	Config      string `json:"config"`
}

// CreateMessageRequest represents the request to create a new message
type CreateMessageRequest struct {
	QueueID     int64      `json:"queue_id" validate:"required"`
	Payload     string     `json:"payload" validate:"required"`
	Priority    int        `json:"priority"`
	ScheduledAt *time.Time `json:"scheduled_at"`
	MaxRetries  int        `json:"max_retries"`
}
