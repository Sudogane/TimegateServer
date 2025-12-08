package services

import (
	"context"

	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/server"
)

type BaseService struct {
	db  *models.Queries
	rdb *cache.RedisClient
	ctx context.Context
}

func NewBaseService(gs server.GameServerInterface) *BaseService {
	return &BaseService{
		db:  gs.GetDB(),
		rdb: gs.GetRDB(),
		ctx: gs.Ctx(),
	}
}
