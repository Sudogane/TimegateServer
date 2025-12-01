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

const (
	MAX_CONNECTIONS             = 50
	MIN_CONNECTIONS             = 10
	MAX_CONN_LIFETIME           = 30 * time.Minute
	MAX_CONN_IDLE_TIME          = 5 * time.Minute
	HEALTH_CHECK_PERIOD         = 1 * time.Minute
	CONNECT_TIMEOUT             = 5 * time.Second
	IDLE_IN_TRANSACTION_TIMEOUT = "10000"
)

func NewRepository(connString string) (*Repository, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	config.MaxConns = MAX_CONNECTIONS
	config.MinConns = MIN_CONNECTIONS
	config.MaxConnLifetime = MAX_CONN_LIFETIME
	config.MaxConnIdleTime = MAX_CONN_IDLE_TIME
	config.HealthCheckPeriod = HEALTH_CHECK_PERIOD
	config.ConnConfig.ConnectTimeout = CONNECT_TIMEOUT
	config.ConnConfig.RuntimeParams["idle_in_transaction_session_timeout"] = IDLE_IN_TRANSACTION_TIMEOUT

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create new pool with config: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
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
