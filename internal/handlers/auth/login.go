package auth

import (
	"context"
	"errors"

	"github.com/sudogane/project_timegate/internal/crypt"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type LoginHandler struct {
	gameServer server.GameServerInterface
}

func NewLoginHandler(gs server.GameServerInterface) *LoginHandler {
	return &LoginHandler{
		gameServer: gs,
	}
}

func validateLoginData(loginData *packets.LoginRequest, lh *LoginHandler) error {
	username := loginData.Username
	password := loginData.Password

	if username == "" || password == "" {
		return errors.New("empty data")
	}

	if len(username) < 3 || len(username) > 10 {
		return errors.New("length of username is too short or too long")
	}

	invalidCredentials := errors.New("invalid credentials")

	user, doesUsernameExist := lh.gameServer.GetDB().Queries.GetUserByUsername(context.Background(), username)
	if doesUsernameExist != nil {
		return invalidCredentials
	}

	if isPasswordCorrect := crypt.ComparePassword(password, user.PasswordHash); !isPasswordCorrect {
		return invalidCredentials
	}

	return nil
}

func (lh *LoginHandler) Handle(session *server.PlayerSession, loginRequestData *packets.LoginRequest, userService *services.UserService) {
	if err := validateLoginData(loginRequestData, lh); err != nil {
		errorPacket := packets.NewErrorMessage(packets.ErrorCode_INVALID_CREDENTIALS)
		lh.gameServer.SendMessage(session.ID, errorPacket)
		return
	}

	user, err := userService.GetByUsername(loginRequestData.Username)
	if err != nil {
		errorPacket := packets.NewErrorMessage(packets.ErrorCode_UNKOWN_ERROR)
		lh.gameServer.SendMessage(session.ID, errorPacket)
		return
	}

	userResources, err := userService.GetUserWithResources(user.ID)
	if err != nil {
		errorPacket := packets.NewErrorMessage(packets.ErrorCode_UNKOWN_ERROR)
		lh.gameServer.SendMessage(session.ID, errorPacket)
		return
	}

	session.PlayerId = user.ID

	loginPacket := packets.NewLoginResponse("user access token", user.Username, userResources.Level.Int32, userResources.Exp.Int32, userResources.Bits.Int64, userResources.Yen.Int64, userResources.StaminaCurrent.Int32, userResources.StaminaMax.Int32)
	lh.gameServer.SendMessage(session.ID, loginPacket)
}
