package handler

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type DigimonHandler struct {
	BaseHandler
	digimonService *services.DigimonService
}

func NewDigimonHandler(server server.GameServerInterface, digimonService *services.DigimonService) *DigimonHandler {
	return &DigimonHandler{
		BaseHandler: *NewBaseHandler(server),

		digimonService: digimonService,
	}
}

func (h *DigimonHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	packetType := msg.GetPacketType()

	switch packetType {
	case packets.PacketType_VIEW_DIGIMON_REQUEST:
		h.onShowDigimon(session)
	}

	return nil
}

func (h *DigimonHandler) onShowDigimon(session *server.PlayerSession) {
	userDigimon, err := h.digimonService.GetUserDigimonByStarterFlag(session.PlayerId)
	if err != nil {
		session.Log("ERROR", err.Error())
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		return
	}

	digimonResponseData := &packets.DigimonData{
		Species:   userDigimon.Species.String,
		Level:     userDigimon.Level.Int32,
		Exp:       userDigimon.Exp.Int32,
		CurrentHp: userDigimon.CurrentHealth.Int32,
		CurrentMp: userDigimon.CurrentMana.Int32,
		Health:    userDigimon.Health.Int32,
		Mana:      userDigimon.Mana.Int32,
		Attack:    userDigimon.Attack.Int32,
		Defense:   userDigimon.Defense.Int32,
		Speed:     userDigimon.Speed.Int32,
	}
	digimonResponse := packets.NewDigimonTeamViewResponse(digimonResponseData)
	h.Send(session, digimonResponse)
}
