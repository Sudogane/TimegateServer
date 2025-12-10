package handler

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type DialogueHandler struct {
	BaseHandler
	userService  *services.UserService
	flagsService *services.UserFlagsService
}

func NewDialogueHandler(userService *services.UserService, flagsService *services.UserFlagsService) *DialogueHandler {
	return &DialogueHandler{
		BaseHandler:  *NewBaseHandler(nil),
		userService:  userService,
		flagsService: flagsService,
	}
}

func (h *DialogueHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	packetType := msg.GetPacketType()

	switch packetType {
	case packets.PacketType_DIALOGUE_CHOICE_SELECTED:
		h.onDialogueChoiceSelected(session, msg.GetDialogueChoiceSelected())
	}

	return nil
}

func (h *DialogueHandler) onDialogueChoiceSelected(session *server.PlayerSession, choice *packets.DialogueChoiceSelected) {
	dialogueId := choice.GetDialogueId()
	choiceId := choice.GetDialogueChoiceId()

	if dialogueId == "DEVELOPMENT" {
		speciesIdMap := map[string]int32{
			"Morphomon":            1,
			"Alphamon":             2,
			"Chronomon: Holy Mode": 3,
		}
		starterDigimons := []string{"Alphamon", "Morphomon", "Chronomon: Holy Mode"}
		digimonSelected := starterDigimons[choiceId]
		if digimonSelected == "" {
			digimonSelected = starterDigimons[0]
		}

		err := h.userService.GiveDigimonToUser(session.PlayerId, speciesIdMap[digimonSelected], true, true)
		if err != nil {
			session.Log("ERROR", err.Error())
			h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
			return
		}

		err = h.flagsService.SetUserFlag(session.PlayerId, "has_selected_starter", "{\"active\": true}")
		if err != nil {
			session.Log("ERROR", err.Error())
			h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
			return
		}
	}
}
