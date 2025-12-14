package services

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/server"
)

type DigimonService struct {
	BaseService
}

func NewDigimonService(gs server.GameServerInterface) *DigimonService {
	return &DigimonService{
		BaseService: *NewBaseService(gs),
	}
}

func (s *DigimonService) GetUserDigimonByStarterFlag(userId uuid.UUID) (models.UserDigimon, error) {
	key := cache.GetUserDigimonByFlagStarterKey(userId.String())
	if cachedStarter, err := s.rdb.Get(key); err == nil {
		var starter models.UserDigimon
		if err := json.Unmarshal([]byte(cachedStarter), &starter); err != nil {
			return models.UserDigimon{}, fmt.Errorf("[Digimon Service] failed to unmarshal user digimon: %w", err)
		}
		return starter, nil
	}

	starter, err := s.db.GetUserDigimonByStarterFlag(s.ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return models.UserDigimon{}, fmt.Errorf("[Digimon Service] failed to get by id: %w", err)
	}

	s.rdb.Set(key, starter)
	return starter, nil
}
