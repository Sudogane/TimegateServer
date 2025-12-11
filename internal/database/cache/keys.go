package cache

import "fmt"

const (
	UserById                         = "user:id:%s"
	UserByUsername                   = "user:username:%s"
	UserWithResourcesById            = "user:with-resources:%s"
	UserUnlockedChapters             = "user:unlocked-chapters:%s"
	UserAvailableEpisodesByChapterId = "user:available-episodes:chapter:%d:%s"
	UserAvailableStagesByEpisodeId   = "user:available-stages:episode:%d:%s"
	UserDigimonByFlagStarter         = "user:digimon:flag-starter:%s"

	UserGetAllFlags   = "user:get-all-flags:%s"
	UserGetFlagByName = "user:get-flag-by-name:%s:%s"

	DigimonById = "digimon:id:%s"
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

func GetUserDigimonByFlagStarterKey(userId string) string {
	return fmt.Sprintf(UserDigimonByFlagStarter, userId)
}

func GetDigimonByIdKey(id string) string {
	return fmt.Sprintf(DigimonById, id)
}

func GetUserGetAllFlagsKey(userId string) string {
	return fmt.Sprintf(UserGetAllFlags, userId)
}

func GetUserGetFlagByNameKey(userId string, name string) string {
	return fmt.Sprintf(UserGetFlagByName, userId, name)
}
