package handlers

import (
	"fmt"

	"github.com/sudogane/project_timegate/internal/handlers/auth"
	"github.com/sudogane/project_timegate/internal/handlers/development"
	"github.com/sudogane/project_timegate/internal/handlers/user"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type Router struct {
	gameServer server.GameServerInterface

	loginHandler    *auth.LoginHandler
	registerHandler *auth.RegisterHandler
	devHandler      *development.TestHandler
	userHandler     *user.UserHandler

	userService *services.UserService
}

func NewRouter(server server.GameServerInterface) *Router {
	return &Router{
		gameServer:      server,
		loginHandler:    auth.NewLoginHandler(server),
		registerHandler: auth.NewRegisterHandler(server),
		devHandler:      development.NewTestHandler(server),
		userHandler:     user.NewUserHandler(server),

		userService: services.NewUserService(server),
	}
}

func (r *Router) Route(session *server.PlayerSession, msg *packets.FromClientToServer) {
	if session.ID != msg.UserId {
		fmt.Println(session.ID + " Tried to do some shenanigans")
		return
	}

	switch msg.PacketType {
	case packets.PacketType_LOGIN_REQUEST:
		r.loginHandler.Handle(session, msg.GetLoginRequest(), r.userService)
	case packets.PacketType_REGISTRATION_REQUEST:
		r.registerHandler.Handle(session, msg.GetRegisterRequest(), r.userService)
	case packets.PacketType_DEVELOPMENT:
		r.devHandler.Handle(session, r.userService)
	// Stages
	case packets.PacketType_CHAPTER_DATA_REQUEST:
		r.userHandler.GetUserChapters(session, r.userService)
	case packets.PacketType_EPISODES_BY_CHAPTER_REQUEST:
		r.userHandler.GetEpisodesByChapter(session, msg.GetEpisodesByChapterRequest(), r.userService)
	}
}
