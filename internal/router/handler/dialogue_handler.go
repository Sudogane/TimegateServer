package handler

import (
	"github.com/google/uuid"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type DialogueHandler struct {
	BaseHandler
	userService            *services.UserService
	flagsService           *services.UserFlagsService
	userAchievementService *services.UserAchievementsService
}

func NewDialogueHandler(server server.GameServerInterface, userService *services.UserService, flagsService *services.UserFlagsService, userAchievementService *services.UserAchievementsService) *DialogueHandler {
	return &DialogueHandler{
		BaseHandler:            *NewBaseHandler(server),
		userService:            userService,
		flagsService:           flagsService,
		userAchievementService: userAchievementService,
	}
}

func (h *DialogueHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	packetType := msg.GetPacketType()

	switch packetType {
	case packets.PacketType_DIALOGUE_CHOICE_SELECTED:
		h.onDialogueChoiceSelected(session, msg.GetDialogueChoiceSelected())
	case packets.PacketType_DIALOGUE_FINISHED:
		h.onDialogueFinished(session, msg.GetDialogueFinished())
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

type DialogueAchievements struct {
	DialogueId    string
	AchievementId int
}

func (h *DialogueHandler) onDialogueFinished(session *server.PlayerSession, data *packets.DialogueFinished) {
	// TEMPORARY ZONE
	dialogueAchievements := []DialogueAchievements{
		{DialogueId: "DEVELOPMENT", AchievementId: 3},
		{DialogueId: "DEBUG", AchievementId: 1},
		{DialogueId: "DBG", AchievementId: 2},
	}

	for _, da := range dialogueAchievements {
		if da.DialogueId == data.GetDialogueId() {
			achievementFromDb, err := h.userAchievementService.GetAchievement(int32(da.AchievementId))
			if err != nil {
				session.Logger.Errorw(err.Error())
				h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
				return
			}

			uAchiev, err := h.userAchievementService.GiveAchievementToUser(session.PlayerId, da.AchievementId)
			if err != nil {
				session.Logger.Errorw(err.Error())
				h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
				return
			}
			responsePacket := &packets.FromServerToClient_AchievementObtained{
				AchievementObtained: &packets.AchievementObtained{
					Achievement: packets.NewAchievementResponse(achievementFromDb.Name, ""),
				},
			}

			session.Logger.Info(uAchiev)

			h.Send(session, responsePacket)
		}
	}
	// TEMPORARY ZONE
}
