package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/server"
)

type UserService struct {
	gs server.GameServerInterface
}

func NewUserService(gameserver server.GameServerInterface) *UserService {
	return &UserService{
		gs: gameserver,
	}
}

func (us *UserService) GetById(id uuid.UUID) (*models.User, error) {
	user, err := us.gs.GetDB().Queries.GetUserById(us.gs.Ctx(), id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get by id: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetByUsername(username string) (*models.User, error) {
	user, err := us.gs.GetDB().Queries.GetUserByUsername(us.gs.Ctx(), username)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get by username: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetUserWithResources(id uuid.UUID) (*models.GetUserWithResourcesRow, error) {
	user, err := us.gs.GetDB().Queries.GetUserWithResources(us.gs.Ctx(), id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get user resources: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetUnlockedChapters(id uuid.UUID) ([]models.GetUserUnlockedChaptersRow, error) {
	user, err := us.gs.GetDB().Queries.GetUserUnlockedChapters(us.gs.Ctx(), id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] failed to get unlocked chapters: %w", err)
	}

	return user, nil
}

func (us *UserService) CreateUserWithResources(username, hashedPassword string) (*models.User, error) {
	user, err := us.gs.GetDB().CreateUserWithResources(us.gs.Ctx(), username, hashedPassword)

	if err != nil {
		return nil, fmt.Errorf("[User Service] create user with resources: %w", err)
	}

	return user, nil
}

func (us *UserService) GetAvailableStages(id uuid.UUID) ([]models.GetAvailableStagesRow, error) {
	user, err := us.gs.GetDB().Queries.GetAvailableStages(us.gs.Ctx(), id)

	if err != nil {
		return nil, fmt.Errorf("[User Service] getting stages: %w", err)
	}

	return user, nil
}

func (us *UserService) GetAvailableEpisodesByChapterId(chapterId int32, userId uuid.UUID) ([]models.Episode, error) {
	user, err := us.gs.GetDB().Queries.GetAvailableEpisodesByChapterId(us.gs.Ctx(), models.GetAvailableEpisodesByChapterIdParams{
		ChapterID: pgtype.Int4{Int32: chapterId, Valid: true},
		UserID:    userId,
	})

	if err != nil {
		return nil, fmt.Errorf("[User Service] getting episodes: %w", err)
	}

	return user, nil
}

func (us *UserService) GetAvailableStagesByEpisodeId(episodeId int32, userId uuid.UUID) ([]models.Stage, error) {
	user, err := us.gs.GetDB().Queries.GetAvailableStagesByEpisodeId(us.gs.Ctx(), models.GetAvailableStagesByEpisodeIdParams{
		EpisodeID: pgtype.Int4{Int32: episodeId, Valid: true},
		UserID:    userId,
	})

	if err != nil {
		return nil, fmt.Errorf("[User Service] getting stages: %w", err)
	}

	return user, nil
}
