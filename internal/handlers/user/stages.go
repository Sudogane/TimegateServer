package user

import (
	"fmt"

	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type UserHandler struct {
	gameServer server.GameServerInterface
}

func NewUserHandler(gs server.GameServerInterface) *UserHandler {
	return &UserHandler{
		gameServer: gs,
	}
}

func (uh *UserHandler) GetUserChapters(session *server.PlayerSession, userService *services.UserService) {
	chapters, err := userService.GetUnlockedChapters(session.PlayerId)
	if err != nil {
		fmt.Println(err)
		return
	}

	var chapterData []*packets.ChapterData
	for _, chapter := range chapters {
		chapterData = append(chapterData, &packets.ChapterData{
			ChapterId:          chapter.ID,
			ChapterName:        chapter.ChapterName,
			ChapterDescription: chapter.Description.String,
			IsUnlocked:         chapter.IsUnlocked.Bool,
			IsBeaten:           chapter.IsBeaten,
		})
	}

	response := &packets.FromServerToClient_ChapterDataResponse{
		ChapterDataResponse: &packets.AccessibleChapterResponse{
			Chapters: chapterData,
		},
	}

	uh.gameServer.SendMessage(session.ID, response)
}

func (uh *UserHandler) GetEpisodesByChapter(session *server.PlayerSession, chapterData *packets.EpisodesByChapterRequest, userService *services.UserService) {
	chapterId := chapterData.GetChapterId()
	if chapterId == 0 {
		return
	}

	episodes, err := userService.GetAvailableEpisodesByChapterId(chapterId, session.PlayerId)
	var episodesData []*packets.EpisodeData
	for _, episode := range episodes {
		episodesData = append(episodesData, &packets.EpisodeData{
			EpisodeId:     episode.ID,
			EpisodeNumber: episode.EpisodeNumber,
			EpisodeName:   episode.EpisodeName,
		})
	}

	if err != nil {
		return
	}

	response := &packets.FromServerToClient_EpisodeDataResponse{
		EpisodeDataResponse: &packets.AccessibleEpisodesResponse{
			Episodes: episodesData,
		},
	}

	uh.gameServer.SendMessage(session.ID, response)
}
