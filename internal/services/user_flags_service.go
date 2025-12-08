package services

import (
	"encoding/json"

	"github.com/google/uuid"
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

func (s *UserFlagsService) GetUserFlags(userId uuid.UUID) ([]models.UserFlag, error) {
	return s.db.GetAllUserFlags(s.ctx, userId)
}

func (s *UserFlagsService) SetUserFlag(userId uuid.UUID, key, value string) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.db.CreateUserFlag(s.ctx, models.CreateUserFlagParams{
		UserID:    userId,
		FlagKey:   key,
		FlagValue: jsonValue,
	})

	return err
}

type flagValue struct {
	Active bool `json:"active"`
}

func (s *UserFlagsService) CheckUserFlag(userId uuid.UUID, key string) (flagValue, error) {
	flag, err := s.db.GetUserFlagByKey(s.ctx, models.GetUserFlagByKeyParams{
		UserID:  userId,
		FlagKey: key,
	})
	fv := flagValue{}

	if err != nil {
		return fv, err
	}

	if err := json.Unmarshal(flag.FlagValue, &fv); err != nil {
		return flagValue{}, err
	}
	return fv, nil
}
