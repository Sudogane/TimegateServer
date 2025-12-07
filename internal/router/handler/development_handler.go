package handler

import (
	"fmt"
	"slices"

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

	allowedSpecies := []string{"Morphomon", "Alphamon", "Chronomon: Holy Mode"}
	speciesIdMap := map[string]int32{
		"Morphomon":            1,
		"Alphamon":             2,
		"Chronomon: Holy Mode": 3,
	}

	speciesSelected := msg.GetDev().GetSTARTER_SPECIES()

	if !slices.Contains(allowedSpecies, speciesSelected) {
		session.Log("ERROR", "Invalid starter selected")
		return fmt.Errorf("invalid starter selected")
	}

	selectedId, ok := speciesIdMap[speciesSelected]
	if !ok {
		session.Log("ERROR", "Invalid starter selected")
		return fmt.Errorf("invalid starter selected")
	}

	err := h.userService.GiveDigimonToUser(session.PlayerId, selectedId, true, false)
	if err != nil {
		session.Log("ERROR", "Error giving digimon to user: "+err.Error())
	} else {
		session.Log("INFO", "Gave digimon to user")
	}

	return nil
}
