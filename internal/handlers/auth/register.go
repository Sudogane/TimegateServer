package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/sudogane/project_timegate/internal/crypt"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type RegisterHandler struct {
	gameServer server.GameServerInterface
}

func NewRegisterHandler(gs server.GameServerInterface) *RegisterHandler {
	return &RegisterHandler{
		gameServer: gs,
	}
}

func validateRegisterData(username, password string, rh *RegisterHandler, sessionId string) error {
	if len(username) < 3 || len(username) > 10 {
		rh.gameServer.SendErrorMessage(sessionId, packets.ErrorCode_USERNAME_TOO_SHORT_OR_LONG)
		return errors.New("username too long or too short")
	}

	if len(password) < 3 {
		rh.gameServer.SendErrorMessage(sessionId, packets.ErrorCode_PASSWORD_TOO_SHORT)
		return errors.New("password too short")
	}

	_, usernameWasTaken := rh.gameServer.GetDB().Queries.GetUserByUsername(context.Background(), username)

	if usernameWasTaken == nil {
		rh.gameServer.SendErrorMessage(sessionId, packets.ErrorCode_USERNAME_TAKEN)
		return errors.New("username taken")
	}

	return nil
}

func (rh *RegisterHandler) Handle(session *server.PlayerSession, registerRequestData *packets.RegisterRequest, userService *services.UserService) {
	username := registerRequestData.Username
	password := registerRequestData.Password
	if err := validateRegisterData(username, password, rh, session.ID); err != nil {
		return
	}

	hashedPassword, _ := crypt.HashPassword(password)

	if _, err := userService.CreateUserWithResources(username, hashedPassword); err != nil {
		fmt.Println("Ocorreu algum erro ao criar o usuario")
		rh.gameServer.SendErrorMessage(session.ID, packets.ErrorCode_UNKOWN_ERROR)
		return
	}

	registerPacket := packets.NewRegisterResponse("user access token")
	rh.gameServer.SendMessage(session.ID, registerPacket)
}
