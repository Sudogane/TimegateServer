package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sudogane/project_timegate/internal/database"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type GameServerInterface interface {
	GetSession(sessionID string) *PlayerSession
	SendMessage(sessionId string, message packets.ServerPayload)
	AddSession(session *PlayerSession)
	RemoveSession(sessionID string)
	GetDB() *database.Repository
	Ctx() context.Context
	SendErrorMessage(sessionId string, code packets.ErrorCode)
}

type RouterInterface interface {
	Route(session *PlayerSession, msg *packets.FromClientToServer)
}

type PlayerSession struct {
	ID       string
	Conn     *websocket.Conn
	SendChan chan *packets.FromServerToClient

	PlayerId uuid.UUID
}
