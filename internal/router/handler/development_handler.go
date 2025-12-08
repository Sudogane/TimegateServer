package handler

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type DevelopmentHandler struct {
	BaseHandler
	userService  *services.UserService
	flagsService *services.UserFlagsService
}

func NewDevelopmentHandle(userService *services.UserService, flagsService *services.UserFlagsService) *DevelopmentHandler {
	return &DevelopmentHandler{
		BaseHandler:  *NewBaseHandler(nil),
		userService:  userService,
		flagsService: flagsService,
	}
}

func (h *DevelopmentHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	session.Log("INFO", "Received development packet")

	//h.flagsService.SetUserFlag(session.PlayerId, "test_flag", "{\"active\": true}")

	//test, err := h.flagsService.GetUserFlags(session.PlayerId)
	//session.Log("DEBUG", fmt.Sprintf("test: %v", test)+" "+fmt.Sprintf("err: %v", err))

	//flagActive, _ := h.flagsService.CheckUserFlag(session.PlayerId, "test_flag")
	//fmt.Println(flagActive.Active)

	/*allowedSpecies := []string{"Morphomon", "Alphamon", "Chronomon: Holy Mode"}
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
	}*/

	return nil
}
