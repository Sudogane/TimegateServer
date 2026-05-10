package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/server"
)

type UserService struct {
	BaseService
}

func NewUserService(gs server.GameServerInterface) *UserService {
	return &UserService{
		BaseService: *NewBaseService(gs),
	}
}

func (us *UserService) CheckIfUsernameIsTaken(username string) (bool, error) {
	return us.db.CheckIfUsernameIsTaken(us.ctx, username)
}

func (us *UserService) GetById(id uuid.UUID) (*models.User, error) {
	key := cache.GetUserByIdKey(id.String())

	if cachedUser, err := us.rdb.Get(key); err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
		}
		return &user, nil
	}

	user, err := us.db.GetUserById(us.ctx, id)
	if err != nil {
		return nil, errors.New("[User Service] failed to get user by ID: " + err.Error())
	}

	us.rdb.Set(key, user)
	return &user, nil
}

func (us *UserService) GetByUsername(username string) (*models.User, error) {
	key := cache.GetUserByUsernameKey(username)

	if cachedUser, err := us.rdb.Get(key); err == nil {
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
		}
		return &user, nil
	}

	user, err := us.db.GetUserByUsername(us.ctx, username)
	if err != nil {
		return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
	}

	us.rdb.Set(key, user)
	return &user, nil
}

func (us *UserService) GetUserWithResources(id uuid.UUID) (*models.GetUserWithResourcesRow, error) {
	key := cache.GetUserWithResourcesKey(id.String())

	if cachedUser, err := us.rdb.Get(key); err == nil {
		var user models.GetUserWithResourcesRow
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
		}
		return &user, nil
	}

	user, err := us.db.GetUserWithResources(us.ctx, id)
	if err != nil {
		return nil, errors.New("[User Service] failed to get user by id: " + err.Error())
	}

	us.rdb.Set(key, user)
	return &user, nil
}

func (us *UserService) GetUnlockedChapters(id uuid.UUID) ([]models.GetUserUnlockedChaptersRow, error) {
	key := cache.GetUserUnlockedChaptersKey(id.String())

	if cachedChapters, err := us.rdb.Get(key); err == nil {
		var chapters []models.GetUserUnlockedChaptersRow
		if err := json.Unmarshal([]byte(cachedChapters), &chapters); err != nil {
			return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
		}
		return chapters, nil
	}

	chapters, err := us.db.GetUserUnlockedChapters(us.ctx, id)
	if err != nil {
		return nil, errors.New("[User Service] failed to get user by id: " + err.Error())
	}

	us.rdb.Set(key, chapters)
	return chapters, nil
}

func (us *UserService) CreateUserWithResources(username, hashedPassword string) (*models.User, error) {
	user, err := us.db.CreateUser(us.ctx, models.CreateUserParams{Username: username, PasswordHash: hashedPassword})

	if err != nil {
		return nil, fmt.Errorf("[User Service] create user: %w", err)
	}

	err = us.db.CreatePlayerResources(us.ctx, user.ID)

	if err != nil {
		return nil, fmt.Errorf("[User Service] create player resources: %w", err)
	}

	return &user, nil
}

func (us *UserService) GetAvailableEpisodesByChapterId(chapterId int32, userId uuid.UUID) ([]models.Episode, error) {
	key := cache.GetUserAvailableEpisodesByChapterIdKey(userId.String(), chapterId)

	if cachedEpisodes, err := us.rdb.Get(key); err == nil {
		var episodes []models.Episode
		if err := json.Unmarshal([]byte(cachedEpisodes), &episodes); err != nil {
			return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
		}
		return episodes, nil
	}

	episodes, err := us.db.GetAvailableEpisodesByChapterId(us.ctx, models.GetAvailableEpisodesByChapterIdParams{
		ChapterID: pgtype.Int4{Int32: chapterId, Valid: true},
		UserID:    userId,
	})

	if err != nil {
		return nil, errors.New("[User Service] failed to get user episodes: " + err.Error())
	}

	us.rdb.Set(key, episodes)
	return episodes, nil
}

func (us *UserService) GetAvailableStagesByEpisodeId(episodeId int32, userId uuid.UUID) ([]models.Stage, error) {
	key := cache.GetUserAvailableStagesByEpisodeIdKey(userId.String(), episodeId)

	if cachedStages, err := us.rdb.Get(key); err == nil {
		var stages []models.Stage
		if err := json.Unmarshal([]byte(cachedStages), &stages); err != nil {
			return nil, errors.New("[User Service] failed to unmarshal user: " + err.Error())
		}
		return stages, nil
	}

	stages, err := us.db.GetAvailableStagesByEpisodeId(us.ctx, models.GetAvailableStagesByEpisodeIdParams{
		EpisodeID: pgtype.Int4{Int32: episodeId, Valid: true},
		UserID:    userId,
	})

	if err != nil {
		return nil, errors.New("[User Service] failed to get user stages: " + err.Error())
	}

	us.rdb.Set(key, stages)
	return stages, nil
}

func (us *UserService) GiveDigimonToUser(userId uuid.UUID, digimonId int32, isStarter, isLocked bool) error {
	cacheKey := cache.GetUserDigimonByFlagStarterKey(userId.String())
	if cachedStarter, err := us.rdb.Get(cacheKey); err == nil {
		var starter models.UserDigimon
		if err := json.Unmarshal([]byte(cachedStarter), &starter); err != nil {
			return errors.New("[User Service] failed unmarshal user digimon: " + err.Error())
		}
		return errors.New("[User Service] User already has a starter")
	}

	if doesUserHaveStarter, _ := us.db.GetUserDigimonByStarterFlag(us.ctx, pgtype.UUID{Bytes: userId, Valid: true}); doesUserHaveStarter.ID != uuid.Nil && isStarter {
		us.rdb.Set(cacheKey, doesUserHaveStarter)
		return errors.New("[User Service] User already has a starter")
	}

	_, err := us.db.CreateUserDigimon(us.ctx, models.CreateUserDigimonParams{
		UserID:    pgtype.UUID{Bytes: userId, Valid: true},
		BaseID:    pgtype.Int4{Int32: digimonId, Valid: true},
		IsStarter: pgtype.Bool{Bool: isStarter, Valid: true},
		IsLocked:  pgtype.Bool{Bool: isLocked, Valid: true},
	})

	if err != nil {
		return errors.New("[User Service] Error creating user starter: " + err.Error())
	}

	return nil
}
