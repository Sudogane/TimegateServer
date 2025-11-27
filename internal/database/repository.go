package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sudogane/project_timegate/internal/database/models"
)

type Repository struct {
	Queries *models.Queries
	DB      *pgxpool.Pool
}

// AddDigimonToTeam implements models.Querier.
func (r *Repository) AddDigimonToTeam(ctx context.Context, arg models.AddDigimonToTeamParams) error {
	panic("unimplemented")
}

// AddUserBits implements models.Querier.
func (r *Repository) AddUserBits(ctx context.Context, arg models.AddUserBitsParams) error {
	panic("unimplemented")
}

// CheckUserChapterCompletion implements models.Querier.
func (r *Repository) CheckUserChapterCompletion(ctx context.Context, arg models.CheckUserChapterCompletionParams) (models.CheckUserChapterCompletionRow, error) {
	panic("unimplemented")
}

// ClearTeamSlot implements models.Querier.
func (r *Repository) ClearTeamSlot(ctx context.Context, arg models.ClearTeamSlotParams) error {
	panic("unimplemented")
}

// CompleteStage implements models.Querier.
func (r *Repository) CompleteStage(ctx context.Context, arg models.CompleteStageParams) error {
	panic("unimplemented")
}

// CompleteUserChapter implements models.Querier.
func (r *Repository) CompleteUserChapter(ctx context.Context, arg models.CompleteUserChapterParams) error {
	panic("unimplemented")
}

// CreatePlayerResources implements models.Querier.
func (r *Repository) CreatePlayerResources(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

// CreateUser implements models.Querier.
func (r *Repository) CreateUser(ctx context.Context, arg models.CreateUserParams) (models.User, error) {
	panic("unimplemented")
}

// CreateUserDigimon implements models.Querier.
func (r *Repository) CreateUserDigimon(ctx context.Context, arg models.CreateUserDigimonParams) (models.UserDigimon, error) {
	panic("unimplemented")
}

// GetAllBaseDigimon implements models.Querier.
func (r *Repository) GetAllBaseDigimon(ctx context.Context) ([]models.BaseDigimon, error) {
	panic("unimplemented")
}

// GetAllUserDigimon implements models.Querier.
func (r *Repository) GetAllUserDigimon(ctx context.Context, id uuid.UUID) ([]models.UserDigimon, error) {
	panic("unimplemented")
}

// GetAvailableEpisodesByChapterId implements models.Querier.
func (r *Repository) GetAvailableEpisodesByChapterId(ctx context.Context, arg models.GetAvailableEpisodesByChapterIdParams) ([]models.Episode, error) {
	panic("unimplemented")
}

// GetAvailableStages implements models.Querier.
func (r *Repository) GetAvailableStages(ctx context.Context, userID uuid.UUID) ([]models.GetAvailableStagesRow, error) {
	panic("unimplemented")
}

// GetAvailableStagesByEpisodeId implements models.Querier.
func (r *Repository) GetAvailableStagesByEpisodeId(ctx context.Context, arg models.GetAvailableStagesByEpisodeIdParams) ([]models.Stage, error) {
	panic("unimplemented")
}

// GetBaseDigimon implements models.Querier.
func (r *Repository) GetBaseDigimon(ctx context.Context, id int32) (models.BaseDigimon, error) {
	panic("unimplemented")
}

// GetBaseDigimonBySpecies implements models.Querier.
func (r *Repository) GetBaseDigimonBySpecies(ctx context.Context, species string) (models.BaseDigimon, error) {
	panic("unimplemented")
}

// GetCompletedStages implements models.Querier.
func (r *Repository) GetCompletedStages(ctx context.Context, userID uuid.UUID) ([]models.UserCompletedStage, error) {
	panic("unimplemented")
}

// GetStageById implements models.Querier.
func (r *Repository) GetStageById(ctx context.Context, id int32) (models.GetStageByIdRow, error) {
	panic("unimplemented")
}

// GetStageWaves implements models.Querier.
func (r *Repository) GetStageWaves(ctx context.Context, id int32) ([]models.GetStageWavesRow, error) {
	panic("unimplemented")
}

// GetUserById implements models.Querier.
func (r *Repository) GetUserById(ctx context.Context, id uuid.UUID) (models.User, error) {
	panic("unimplemented")
}

// GetUserByUsername implements models.Querier.
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	panic("unimplemented")
}

// GetUserDigibank implements models.Querier.
func (r *Repository) GetUserDigibank(ctx context.Context, userID uuid.UUID) ([]models.GetUserDigibankRow, error) {
	panic("unimplemented")
}

// GetUserDigimon implements models.Querier.
func (r *Repository) GetUserDigimon(ctx context.Context, arg models.GetUserDigimonParams) (models.UserDigimon, error) {
	panic("unimplemented")
}

// GetUserTeam implements models.Querier.
func (r *Repository) GetUserTeam(ctx context.Context, userID uuid.UUID) ([]models.GetUserTeamRow, error) {
	panic("unimplemented")
}

// GetUserUnlockedChapters implements models.Querier.
func (r *Repository) GetUserUnlockedChapters(ctx context.Context, userID uuid.UUID) ([]models.GetUserUnlockedChaptersRow, error) {
	panic("unimplemented")
}

// GetUserWithResources implements models.Querier.
func (r *Repository) GetUserWithResources(ctx context.Context, id uuid.UUID) (models.GetUserWithResourcesRow, error) {
	panic("unimplemented")
}

// RemoveDigimonFromTeam implements models.Querier.
func (r *Repository) RemoveDigimonFromTeam(ctx context.Context, arg models.RemoveDigimonFromTeamParams) error {
	panic("unimplemented")
}

// UnlockUserChapter implements models.Querier.
func (r *Repository) UnlockUserChapter(ctx context.Context, arg models.UnlockUserChapterParams) error {
	panic("unimplemented")
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
