package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-fuego/fuego"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/shravan20/qafka/internal/models"
	"github.com/shravan20/qafka/internal/services"
)

func SetupRoutes(app *fuego.Server, queueService *services.QueueService, monitoringService *services.MonitoringService) {
	// Health check
	fuego.Get(app, "/health", func(c fuego.ContextNoBody) (any, error) {
		return map[string]string{"status": "healthy"}, nil
	})

	// API v1 group
	v1 := fuego.Group(app, "/api/v1")

	// Queue routes
	setupQueueRoutes(v1, queueService, monitoringService)

	// Message routes
	setupMessageRoutes(v1, queueService, monitoringService)

	// Worker routes
	setupWorkerRoutes(v1, queueService, monitoringService)

	// Metrics endpoint
	app.Handle(http.MethodGet, "/metrics", promhttp.Handler().ServeHTTP)
}

func SetupSwagger(app *fuego.Server) {
	app.Handle(http.MethodGet, "/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}

func setupQueueRoutes(group *fuego.Group, queueService *services.QueueService, monitoringService *services.MonitoringService) {
	// Get all queues
	// @Summary Get all queues
	// @Description Get a list of all queues
	// @Tags queues
	// @Accept json
	// @Produce json
	// @Success 200 {array} models.Queue
	// @Router /api/v1/queues [get]
	fuego.Get(group, "/queues", func(c fuego.ContextNoBody) (any, error) {
		queues, err := queueService.GetQueues(context.Background())
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get queues",
			}
		}
		return queues, nil
	})

	// Create queue
	// @Summary Create a new queue
	// @Description Create a new message queue
	// @Tags queues
	// @Accept json
	// @Produce json
	// @Param queue body models.CreateQueueRequest true "Queue creation request"
	// @Success 201 {object} models.Queue
	// @Router /api/v1/queues [post]
	fuego.Post(group, "/queues", func(c fuego.ContextWithBody[models.CreateQueueRequest]) (*models.Queue, error) {
		body, err := c.Body()
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
			}
		}

		queue, err := queueService.CreateQueue(context.Background(), &body)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to create queue",
			}
		}

		return queue, nil
	})

	// Get specific queue
	// @Summary Get a queue by ID
	// @Description Get a specific queue by its ID
	// @Tags queues
	// @Accept json
	// @Produce json
	// @Param id path int true "Queue ID"
	// @Success 200 {object} models.Queue
	// @Router /api/v1/queues/{id} [get]
	fuego.Get(group, "/queues/{id}", func(c fuego.ContextNoBody) (any, error) {
		idStr := c.PathParam("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid queue ID",
			}
		}

		queue, err := queueService.GetQueue(context.Background(), id)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusNotFound,
				Message:    "Queue not found",
			}
		}

		return queue, nil
	})

	// Delete queue
	// @Summary Delete a queue
	// @Description Delete a queue by ID
	// @Tags queues
	// @Accept json
	// @Produce json
	// @Param id path int true "Queue ID"
	// @Success 204
	// @Router /api/v1/queues/{id} [delete]
	fuego.Delete(group, "/queues/{id}", func(c fuego.ContextNoBody) (any, error) {
		idStr := c.PathParam("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid queue ID",
			}
		}

		err = queueService.DeleteQueue(context.Background(), id)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to delete queue",
			}
		}

		return nil, nil
	})
}

func setupMessageRoutes(group *fuego.Group, queueService *services.QueueService, monitoringService *services.MonitoringService) {
	// Get messages
	// @Summary Get messages
	// @Description Get messages from a queue
	// @Tags messages
	// @Accept json
	// @Produce json
	// @Param queue_id query int false "Queue ID filter"
	// @Param limit query int false "Limit number of results"
	// @Success 200 {array} models.Message
	// @Router /api/v1/messages [get]
	fuego.Get(group, "/messages", func(c fuego.ContextNoBody) (any, error) {
		queueIDStr := c.QueryParam("queue_id")
		limitStr := c.QueryParam("limit")

		var queueID int64
		var limit int
		var err error

		if queueIDStr != "" {
			queueID, err = strconv.ParseInt(queueIDStr, 10, 64)
			if err != nil {
				return nil, fuego.HTTPError{
					StatusCode: http.StatusBadRequest,
					Message:    "Invalid queue_id parameter",
				}
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				return nil, fuego.HTTPError{
					StatusCode: http.StatusBadRequest,
					Message:    "Invalid limit parameter",
				}
			}
		}

		messages, err := queueService.GetMessages(context.Background(), queueID, limit)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get messages",
			}
		}

		return messages, nil
	})

	// Create message
	// @Summary Create a new message
	// @Description Add a new message to a queue
	// @Tags messages
	// @Accept json
	// @Produce json
	// @Param message body models.CreateMessageRequest true "Message creation request"
	// @Success 201 {object} models.Message
	// @Router /api/v1/messages [post]
	fuego.Post(group, "/messages", func(c fuego.ContextWithBody[models.CreateMessageRequest]) (*models.Message, error) {
		body, err := c.Body()
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid request body",
			}
		}

		message, err := queueService.CreateMessage(context.Background(), &body)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to create message",
			}
		}

		// Update monitoring metrics
		monitoringService.IncrementMessageCounter("queue_"+strconv.FormatInt(body.QueueID, 10), "created")

		return message, nil
	})
}

func setupWorkerRoutes(group *fuego.Group, queueService *services.QueueService, monitoringService *services.MonitoringService) {
	// Get workers
	// @Summary Get workers
	// @Description Get workers for a queue
	// @Tags workers
	// @Accept json
	// @Produce json
	// @Param queue_id query int false "Queue ID filter"
	// @Success 200 {array} models.Worker
	// @Router /api/v1/workers [get]
	fuego.Get(group, "/workers", func(c fuego.ContextNoBody) (any, error) {
		queueIDStr := c.QueryParam("queue_id")

		var queueID int64
		var err error

		if queueIDStr != "" {
			queueID, err = strconv.ParseInt(queueIDStr, 10, 64)
			if err != nil {
				return nil, fuego.HTTPError{
					StatusCode: http.StatusBadRequest,
					Message:    "Invalid queue_id parameter",
				}
			}
		}

		workers, err := queueService.GetWorkers(context.Background(), queueID)
		if err != nil {
			return nil, fuego.HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to get workers",
			}
		}

		return workers, nil
	})
}
