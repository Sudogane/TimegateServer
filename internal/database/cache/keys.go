package cache

import "fmt"

const (
	UserById                         = "user:id:%s"
	UserByUsername                   = "user:username:%s"
	UserWithResourcesById            = "user:with-resources:%s"
	UserUnlockedChapters             = "user:unlocked-chapters:%s"
	UserAvailableEpisodesByChapterId = "user:available-episodes:chapter:%d:%s"
	UserAvailableStagesByEpisodeId   = "user:available-stages:episode:%d:%s"
)

func GetUserByIdKey(id string) string {
	return fmt.Sprintf(UserById, id)
}

func GetUserByUsernameKey(username string) string {
	return fmt.Sprintf(UserByUsername, username)
}

func GetUserWithResourcesKey(id string) string {
	return fmt.Sprintf(UserWithResourcesById, id)
}

func GetUserUnlockedChaptersKey(id string) string {
	return fmt.Sprintf(UserUnlockedChapters, id)
}

func GetUserAvailableEpisodesByChapterIdKey(userId string, chapterId int32) string {
	return fmt.Sprintf(UserAvailableEpisodesByChapterId, chapterId, userId)
}

func GetUserAvailableStagesByEpisodeIdKey(userId string, episodeId int32) string {
	return fmt.Sprintf(UserAvailableStagesByEpisodeId, episodeId, userId)
}
