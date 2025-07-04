package services

import (
	"context"
	"fmt"
	"time"

	"github.com/shravan20/qafka/internal/models"
	"github.com/uptrace/bun"
)

type QueueService struct {
	db *bun.DB
}

func NewQueueService(db *bun.DB) *QueueService {
	return &QueueService{db: db}
}

// Queue operations
func (s *QueueService) CreateQueue(ctx context.Context, req *models.CreateQueueRequest) (*models.Queue, error) {
	queue := &models.Queue{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Config:      req.Config,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.NewInsert().Model(queue).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	return queue, nil
}

func (s *QueueService) GetQueues(ctx context.Context) ([]*models.Queue, error) {
	var queues []*models.Queue
	err := s.db.NewSelect().Model(&queues).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queues: %w", err)
	}
	return queues, nil
}

func (s *QueueService) GetQueue(ctx context.Context, id int64) (*models.Queue, error) {
	queue := &models.Queue{}
	err := s.db.NewSelect().Model(queue).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}
	return queue, nil
}

func (s *QueueService) UpdateQueue(ctx context.Context, id int64, updates map[string]interface{}) (*models.Queue, error) {
	updates["updated_at"] = time.Now()

	_, err := s.db.NewUpdate().Model((*models.Queue)(nil)).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update queue: %w", err)
	}

	return s.GetQueue(ctx, id)
}

func (s *QueueService) DeleteQueue(ctx context.Context, id int64) error {
	_, err := s.db.NewDelete().Model((*models.Queue)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}
	return nil
}

// Message operations
func (s *QueueService) CreateMessage(ctx context.Context, req *models.CreateMessageRequest) (*models.Message, error) {
	message := &models.Message{
		QueueID:     req.QueueID,
		Payload:     req.Payload,
		Priority:    req.Priority,
		Status:      "pending",
		ScheduledAt: req.ScheduledAt,
		MaxRetries:  req.MaxRetries,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if message.MaxRetries == 0 {
		message.MaxRetries = 3
	}

	_, err := s.db.NewInsert().Model(message).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

func (s *QueueService) GetMessages(ctx context.Context, queueID int64, limit int) ([]*models.Message, error) {
	var messages []*models.Message
	query := s.db.NewSelect().Model(&messages).Relation("Queue")

	if queueID > 0 {
		query = query.Where("queue_id = ?", queueID)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	return messages, nil
}

func (s *QueueService) GetNextMessage(ctx context.Context, queueID int64) (*models.Message, error) {
	message := &models.Message{}
	err := s.db.NewSelect().Model(message).
		Where("queue_id = ? AND status = 'pending'", queueID).
		Where("(scheduled_at IS NULL OR scheduled_at <= ?)", time.Now()).
		Order("priority DESC", "created_at ASC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get next message: %w", err)
	}

	return message, nil
}

func (s *QueueService) UpdateMessageStatus(ctx context.Context, messageID int64, status string, errorMessage string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	switch status {
	case "processing":
		updates["processed_at"] = time.Now()
	case "failed":
		updates["failed_at"] = time.Now()
		if errorMessage != "" {
			updates["error_message"] = errorMessage
		}
	}

	_, err := s.db.NewUpdate().Model((*models.Message)(nil)).
		Where("id = ?", messageID).
		Set("status = ?", status).
		Set("updated_at = ?", time.Now()).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	return nil
}

// Worker operations
func (s *QueueService) RegisterWorker(ctx context.Context, name string, queueID int64) (*models.Worker, error) {
	worker := &models.Worker{
		Name:      name,
		QueueID:   queueID,
		Status:    "idle",
		LastPing:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.NewInsert().Model(worker).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to register worker: %w", err)
	}

	return worker, nil
}

func (s *QueueService) GetWorkers(ctx context.Context, queueID int64) ([]*models.Worker, error) {
	var workers []*models.Worker
	query := s.db.NewSelect().Model(&workers).Relation("Queue")

	if queueID > 0 {
		query = query.Where("queue_id = ?", queueID)
	}

	err := query.Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get workers: %w", err)
	}
	return workers, nil
}

func (s *QueueService) UpdateWorkerPing(ctx context.Context, workerID int64, status string) error {
	_, err := s.db.NewUpdate().Model((*models.Worker)(nil)).
		Set("last_ping = ?", time.Now()).
		Set("status = ?", status).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", workerID).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update worker ping: %w", err)
	}

	return nil
}
