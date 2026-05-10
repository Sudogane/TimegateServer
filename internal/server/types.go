package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/logger"
	"github.com/sudogane/project_timegate/pkg/packets"
	"go.uber.org/zap"
)

type GameServerInterface interface {
	GetSession(sessionID string) *PlayerSession
	SendMessage(sessionId string, message packets.ServerPayload)
	AddSession(session *PlayerSession)
	RemoveSession(sessionID string)
	GetDB() *models.Queries
	GetRDB() *cache.RedisClient
	Ctx() context.Context
	SendErrorMessage(sessionId string, code packets.ErrorCode)
	GetLogger() *logger.Logger
}

type RouterInterface interface {
	Route(session *PlayerSession, msg *packets.FromClientToServer)
}

type PlayerSession struct {
	ID       string
	Conn     *websocket.Conn
	SendChan chan *packets.FromServerToClient

	PlayerId uuid.UUID
	Logger   *zap.SugaredLogger
}
