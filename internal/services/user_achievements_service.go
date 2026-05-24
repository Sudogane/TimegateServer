package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/server"
)

type UserAchievementsService struct {
	BaseService
}

func NewUserAchievementsService(gs server.GameServerInterface) *UserAchievementsService {
	return &UserAchievementsService{
		BaseService: *NewBaseService(gs),
	}
}

func (uas *UserAchievementsService) GetUserAchievements(userId uuid.UUID) ([]models.GetAllUserAchievementsRow, error) {
	achievements, err := uas.db.GetAllUserAchievements(uas.ctx, userId)
	if err != nil {
		return nil, errors.New("[User Achievements Service] failed to get achievements: " + err.Error())
	}

	return achievements, nil
}

func (uas *UserAchievementsService) GiveAchievementToUser(userId uuid.UUID, achievementId int) (*models.UserAchievement, error) {
	achievement, err := uas.db.GiveUserAchievement(uas.ctx, models.GiveUserAchievementParams{UserID: userId, ID: int32(achievementId)})
	if err != nil {
		return nil, errors.New("[User Achievements Service] failed to give achievement: " + err.Error())
	}

	return &achievement, nil
}

func (uas *UserAchievementsService) GetAchievement(achievementId int32) (*models.Achievement, error) {
	achievement, err := uas.db.GetAchievementById(uas.ctx, achievementId)
	if err != nil {
		return nil, errors.New("[User Achievements Service] failed to find base achievement: " + err.Error())
	}

	return &achievement, nil
}
