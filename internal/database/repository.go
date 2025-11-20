package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sudogane/project_timegate/internal/database/models"
)

type Repository struct {
	Queries *models.Queries
	DB      *pgxpool.Pool
}

func NewRepository(connString string) (*Repository, error) {
	config, err := pgxpool.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = "10000"

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, fmt.Errorf("unable to create new pool with config: %v", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	} else {
		fmt.Println("Database pinged")
	}

	queries := models.New(pool)

	return &Repository{
		Queries: queries,
		DB:      pool,
	}, nil
}

func (r *Repository) Close() {
	r.DB.Close()
}

func (r *Repository) CreateUserWithResources(ctx context.Context, username, passwordHash string) (*models.User, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.Queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, models.CreateUserParams{Username: username, PasswordHash: passwordHash})
	if err != nil {
		return nil, err
	}

	err = qtx.CreatePlayerResources(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	err = qtx.UnlockUserChapter(ctx, models.UnlockUserChapterParams{
		UserID:     user.ID,
		ChapterID:  1,
		IsUnlocked: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &user, nil
}
