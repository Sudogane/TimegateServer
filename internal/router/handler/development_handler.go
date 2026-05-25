package handler

import (
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type DevelopmentHandler struct {
	BaseHandler
	newService *services.UserAchievementsService
}

func NewDevelopmentHandle(server server.GameServerInterface, ns *services.UserAchievementsService) *DevelopmentHandler {
	return &DevelopmentHandler{
		BaseHandler: *NewBaseHandler(server),
		newService:  ns,
	}
}

func (h *DevelopmentHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	session.Logger.Infow("Received Development Packet")
	achievements, err := h.newService.GetUserAchievements(session.PlayerId)
	if err != nil {
		session.Logger.Error(err)
	}

	achievementsData := make([]*packets.AchievementData, len(achievements))
	for i, achievement := range achievements {
		achievementsData[i] = &packets.AchievementData{
			AchievementGroupName: achievement.AchievementGroupName,
			AchievementName:      achievement.AchievementName,
			ObtainedDate:         achievement.UnlockedAt.Time.Format("02/01/2006"),
		}

		session.Logger.Info(achievement.IsComplete)
	}

	responsePacket := &packets.FromServerToClient_GetAllAchievements{
		GetAllAchievements: &packets.GetAllAchievementsResponse{
			Achievement: achievementsData,
		},
	}

	h.Send(session, responsePacket)

	return nil
}
