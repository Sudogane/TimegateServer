package handler

import (
	"github.com/google/uuid"
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
	if choice == nil || session.PlayerId == uuid.Nil {
		return
	}

	dialogueId := choice.GetDialogueId()
	choiceId := choice.GetDialogueChoiceId()

	if dialogueId == "" || choiceId < 0 {
		return
	}

	if dialogueId == "DEVELOPMENT" {
		speciesIdMap := map[string]int32{
			"Morphomon":            1,
			"Alphamon":             2,
			"Chronomon: Holy Mode": 3,
		}
		starterDigimons := []string{"Alphamon", "Morphomon", "Chronomon: Holy Mode"}

		var digimonSelected string
		if choiceId < 0 || choiceId >= int32(len(starterDigimons)) {
			digimonSelected = starterDigimons[0]
		} else {
			digimonSelected = starterDigimons[choiceId]
		}

		err := h.userService.GiveDigimonToUser(session.PlayerId, speciesIdMap[digimonSelected], true, true)
		if err != nil {
			session.Logger.Errorw(err.Error())
			h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
			return
		}

		err = h.flagsService.UpdateUserFlag(session.PlayerId, "has_selected_starter", true)
		if err != nil {
			session.Logger.Errorw(err.Error())
			h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
			return
		}
	}
}
