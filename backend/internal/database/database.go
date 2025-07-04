package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/shravan20/qafka/internal/models"
)

func Connect(databaseURL string) (*bun.DB, error) {
	// Open a PostgreSQL database.
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseURL)))

	// Create a Bun db on top of it.
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add query hook for debugging in development.
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(ctx context.Context, db *bun.DB) error {
	// Create tables
	models := []interface{}{
		(*models.Queue)(nil),
		(*models.Message)(nil),
		(*models.Worker)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to create table for %T: %w", model, err)
		}
	}

	// Create indexes for better performance
	if err := createIndexes(ctx, db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

func createIndexes(ctx context.Context, db *bun.DB) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_messages_queue_id ON messages(queue_id)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_scheduled_at ON messages(scheduled_at)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_priority ON messages(priority DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_workers_queue_id ON workers(queue_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workers_status ON workers(status)`,
	}

	for _, indexSQL := range indexes {
		if _, err := db.ExecContext(ctx, indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
