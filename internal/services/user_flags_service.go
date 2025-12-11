package services

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/server"
)

type UserFlagsService struct {
	BaseService
}

func NewUserFlagsService(gs server.GameServerInterface) *UserFlagsService {
	return &UserFlagsService{
		BaseService: *NewBaseService(gs),
	}
}

func (s *UserFlagsService) GetAllUserFlags(userId uuid.UUID) ([]models.UserFlag, error) {
	// TODO: Automatize this, try creating a fuction that does s.rdb.GetOne(key, modelType) and s.rdb.GetMany(key, modelType)
	key := cache.GetUserGetAllFlagsKey(userId.String())
	if cachedFlags, err := s.rdb.Get(key); err == nil {
		var flags []models.UserFlag
		if err := json.Unmarshal([]byte(cachedFlags), &flags); err != nil {
			return nil, fmt.Errorf("[User Flags Service] failed to unmarshal user: %w", err)
		}
		return flags, nil
	}

	allFlags, err := s.db.GetAllUserFlags(s.ctx, userId)
	if err != nil {
		return nil, err
	}

	s.rdb.Set(key, allFlags)
	return allFlags, nil
}

func (s *UserFlagsService) GetUserFlag(userId uuid.UUID, name string) (models.UserFlag, error) {
	key := cache.GetUserGetFlagByNameKey(userId.String(), name)
	if cachedFlag, err := s.rdb.Get(key); err == nil {
		var flag models.UserFlag
		if err := json.Unmarshal([]byte(cachedFlag), &flag); err != nil {
			return models.UserFlag{}, fmt.Errorf("[User Flags Service] failed to unmarshal user: %w", err)
		}
		return flag, nil
	}

	flag, err := s.db.GetUserFlagByName(s.ctx, models.GetUserFlagByNameParams{
		UserID:   userId,
		FlagName: name,
	})

	if err != nil {
		return models.UserFlag{}, err
	}

	s.rdb.Set(key, flag)
	return flag, nil
}

func (s *UserFlagsService) UpdateUserFlag(userId uuid.UUID, name string, isActive bool) error {
	_, err := s.db.GetUserFlagByName(s.ctx, models.GetUserFlagByNameParams{
		UserID:   userId,
		FlagName: name,
	})

	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	if err == pgx.ErrNoRows {
		_, err := s.db.CreateUserFlag(s.ctx, models.CreateUserFlagParams{
			UserID:   userId,
			FlagName: name,
			IsActive: pgtype.Bool{Bool: isActive, Valid: true},
		})

		if err != nil {
			return err
		}

		return nil
	}

	_, err = s.db.UpdateUserFlag(s.ctx, models.UpdateUserFlagParams{
		UserID:   userId,
		FlagName: name,
		IsActive: pgtype.Bool{Bool: isActive, Valid: true},
	})

	return err
}
