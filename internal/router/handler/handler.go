package handler

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type Handler interface {
	Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error
}

type BaseHandler struct {
	server server.GameServerInterface
}

func NewBaseHandler(server server.GameServerInterface) *BaseHandler {
	return &BaseHandler{
		server: server,
	}
}

func (h *BaseHandler) Send(session *server.PlayerSession, packet packets.ServerPayload) {
	h.server.SendMessage(session.ID, packet)
}

func (h *BaseHandler) SendError(session *server.PlayerSession, code packets.ErrorCode) {
	errorPacket := packets.NewErrorMessage(code)
	h.Send(session, errorPacket)
}
