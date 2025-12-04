package handler

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type DevelopmentHandler struct {
	BaseHandler
	userService *services.UserService
}

func NewDevelopmentHandle(userService *services.UserService) *DevelopmentHandler {
	return &DevelopmentHandler{
		BaseHandler: *NewBaseHandler(nil),
		userService: userService,
	}
}

func (h *DevelopmentHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	session.Log("INFO", "Received development packet")

	err := h.userService.GiveDigimonToUser(session.PlayerId, 1, true, false)
	if err != nil {
		session.Log("ERROR", "Error giving digimon to user: "+err.Error())
	} else {
		session.Log("INFO", "Gave digimon to user")
	}

	return nil
}
