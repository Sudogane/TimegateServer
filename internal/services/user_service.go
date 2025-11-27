package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sudogane/project_timegate/internal/database/models"
)

type UserService struct {
	querier models.Querier
	ctx     context.Context
}

func NewUserService(querier models.Querier) *UserService {
	return &UserService{
		querier: querier,
		ctx:     context.Background(),
	}
}

func (us *UserService) GetById(id uuid.UUID) (*models.User, error) {
	user, err := us.querier.GetUserById(us.ctx, id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get by id: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetByUsername(username string) (*models.User, error) {
	user, err := us.querier.GetUserByUsername(us.ctx, username)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get by username: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetUserWithResources(id uuid.UUID) (*models.GetUserWithResourcesRow, error) {
	user, err := us.querier.GetUserWithResources(us.ctx, id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get user resources: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetUnlockedChapters(id uuid.UUID) ([]models.GetUserUnlockedChaptersRow, error) {
	user, err := us.querier.GetUserUnlockedChapters(us.ctx, id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get unlocked chapters: %w", err)
	}

	return user, nil
}

func (us *UserService) CreateUserWithResources(username, hashedPassword string) (*models.User, error) {
	user, err := us.querier.CreateUser(us.ctx, models.CreateUserParams{Username: username, PasswordHash: hashedPassword})

	if err != nil {
		return nil, fmt.Errorf("[User Service] create user: %w", err)
	}

	err = us.querier.CreatePlayerResources(us.ctx, user.ID)

	if err != nil {
		return nil, fmt.Errorf("[User Service] create player resources: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetAvailableStages(id uuid.UUID) ([]models.GetAvailableStagesRow, error) {
	user, err := us.querier.GetAvailableStages(us.ctx, id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] getting stages: %w", err)
	}

	return user, nil
}

func (us *UserService) GetAvailableEpisodesByChapterId(chapterId int32, userId uuid.UUID) ([]models.Episode, error) {
	user, err := us.querier.GetAvailableEpisodesByChapterId(us.ctx, models.GetAvailableEpisodesByChapterIdParams{
		ChapterID: pgtype.Int4{Int32: chapterId, Valid: true},
		UserID:    userId,
	})

	if err != nil {
		return nil, fmt.Errorf("[User Service] getting episodes: %w", err)
	}

	return user, nil
}

func (us *UserService) GetAvailableStagesByEpisodeId(episodeId int32, userId uuid.UUID) ([]models.Stage, error) {
	user, err := us.querier.GetAvailableStagesByEpisodeId(us.ctx, models.GetAvailableStagesByEpisodeIdParams{
		EpisodeID: pgtype.Int4{Int32: episodeId, Valid: true},
		UserID:    userId,
	})

	if err != nil {
		return nil, fmt.Errorf("[User Service] getting stages: %w", err)
	}

	return user, nil
}
