package handler

import (
	"fmt"

	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type StagesHandler struct {
	BaseHandler
	userService *services.UserService
}

func NewStagesHandler(server server.GameServerInterface) *StagesHandler {
	return &StagesHandler{
		BaseHandler: *NewBaseHandler(server),
		userService: services.NewUserService(server.GetDB()),
	}
}

func (h *StagesHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	packetType := msg.GetPacketType()

	switch packetType {
	case packets.PacketType_CHAPTER_DATA_REQUEST:
		h.onGetUserChapters(session)
	case packets.PacketType_EPISODES_BY_CHAPTER_REQUEST:
		h.onGetUserEpisodesByChapter(session, msg.GetEpisodesByChapterRequest().ChapterId)
	}

	return nil
}

func (h *StagesHandler) onGetUserChapters(session *server.PlayerSession) {
	chapters, err := h.userService.GetUnlockedChapters(session.PlayerId)
	if err != nil {
		fmt.Println("failed to get unlocked chapters: ", err)
		return
	}

	chapterData := make([]*packets.ChapterData, len(chapters))
	for i, chapter := range chapters {
		chapterData[i] = &packets.ChapterData{
			ChapterId:          chapter.ID,
			ChapterName:        chapter.ChapterName,
			ChapterDescription: chapter.Description.String,
			IsUnlocked:         chapter.IsUnlocked.Bool,
			IsBeaten:           chapter.IsBeaten,
		}
	}

	response := &packets.FromServerToClient_ChapterDataResponse{
		ChapterDataResponse: &packets.AccessibleChapterResponse{
			Chapters: chapterData,
		},
	}

	h.Send(session, response)
}

func (h *StagesHandler) onGetUserEpisodesByChapter(session *server.PlayerSession, chapterId int32) {
	if chapterId == 0 {
		return
	}

	episodes, err := h.userService.GetAvailableEpisodesByChapterId(chapterId, session.PlayerId)
	if err != nil {
		fmt.Println("failed to get available episodes: ", err)
		return
	}

	episodeData := make([]*packets.EpisodeData, len(episodes))
	for i, episode := range episodes {
		episodeData[i] = &packets.EpisodeData{
			EpisodeId:     episode.ID,
			EpisodeNumber: episode.EpisodeNumber,
		}
	}

	response := &packets.FromServerToClient_EpisodeDataResponse{
		EpisodeDataResponse: &packets.AccessibleEpisodesResponse{
			Episodes: episodeData,
		},
	}

	h.Send(session, response)
}
