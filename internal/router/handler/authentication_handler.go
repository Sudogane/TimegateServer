package handler

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sudogane/project_timegate/internal/crypt"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type AuthenticationHandler struct {
	BaseHandler
	userService  *services.UserService
	flagsService *services.UserFlagsService
}

func NewAuthenticationHandler(server server.GameServerInterface) *AuthenticationHandler {
	return &AuthenticationHandler{
		BaseHandler:  *NewBaseHandler(server),
		userService:  services.NewUserService(server),
		flagsService: services.NewUserFlagsService(server),
	}
}

func (h *AuthenticationHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	packet := msg.GetAuthenticationRequest()

	switch packet.GetType() {
	case packets.AuthenticationType_LOGIN:
		h.handleUserLogin(session, packet)
	case packets.AuthenticationType_REGISTER:
		h.handleUserRegister(session, packet)
	}

	return nil
}

func (h *AuthenticationHandler) handleUserLogin(session *server.PlayerSession, loginRequestData *packets.AuthenticationRequest) {
	if err := h.validateLoginData(loginRequestData, h.userService); err != nil {
		h.SendError(session, packets.ErrorCode_INVALID_CREDENTIALS)
		return
	}

	user, _ := h.userService.GetByUsername(loginRequestData.Username)
	resources, _ := h.userService.GetUserWithResources(user.ID)
	redirectDialogueId := ""
	starterFlag, err := h.flagsService.GetUserFlag(user.ID, "has_selected_starter")

	if err == pgx.ErrNoRows || err == nil && !starterFlag.IsActive.Bool { {
		redirectDialogueId = "DEVELOPMENT"
	}

	session.PlayerId = user.ID
	userDataPacket := &packets.UserData{
		Username:       user.Username,
		Level:          resources.Level.Int32,
		Exp:            resources.Exp.Int32,
		Bits:           resources.Bits.Int64,
		Yen:            resources.Yen.Int64,
		StaminaCurrent: resources.StaminaCurrent.Int32,
		StaminaMax:     resources.StaminaMax.Int32,
	}

	responsePacket := packets.NewAuthenticationResponse("user access token", userDataPacket, redirectDialogueId)
	h.Send(session, responsePacket)

}

func (h *AuthenticationHandler) validateLoginData(loginData *packets.AuthenticationRequest, userService *services.UserService) error {
	username := loginData.GetUsername()
	password := loginData.GetPassword()

	if username == "" || password == "" {
		return errors.New("empty data")
	}

	if len(username) < 3 || len(username) > 10 {
		return errors.New("length of username is too short or too long")
	}

	invalidCredentialsError := errors.New("invalid credentials")

	user, doesUsernameExist := userService.GetByUsername(username)
	if doesUsernameExist != nil {
		return invalidCredentialsError
	}

	if isPasswordCorrect := crypt.ComparePassword(password, user.PasswordHash); !isPasswordCorrect {
		return invalidCredentialsError
	}

	return nil
}

func (h *AuthenticationHandler) handleUserRegister(session *server.PlayerSession, registerRequestData *packets.AuthenticationRequest) {
	if code, err := h.validateRegisterData(registerRequestData, h.userService); err != nil {
		h.SendError(session, code)
		return
	}

	hashedPassword, _ := crypt.HashPassword(registerRequestData.GetPassword())
	user, err := h.userService.CreateUserWithResources(registerRequestData.GetUsername(), hashedPassword)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Log("ERROR", err.Error())
		return
	}

	resources, _ := h.userService.GetUserWithResources(user.ID)

	userData := &packets.UserData{
		Username:       user.Username,
		Level:          resources.Level.Int32,
		Exp:            resources.Exp.Int32,
		Bits:           resources.Bits.Int64,
		Yen:            resources.Yen.Int64,
		StaminaCurrent: resources.StaminaCurrent.Int32,
		StaminaMax:     resources.StaminaMax.Int32,
	}
	session.PlayerId = user.ID

	responsePacket := packets.NewAuthenticationResponse("user access token", userData, "DEVELOPMENT")
	h.Send(session, responsePacket)
}

func (h *AuthenticationHandler) validateRegisterData(registerData *packets.AuthenticationRequest, userService *services.UserService) (packets.ErrorCode, error) {
	username := registerData.GetUsername()
	password := registerData.GetPassword()

	if len(username) < 3 || len(username) > 10 {
		return packets.ErrorCode_USERNAME_TOO_SHORT_OR_LONG, errors.New("username too long or too short")
	}

	if len(password) < 3 {
		return packets.ErrorCode_PASSWORD_TOO_SHORT, errors.New("username too long or too short")
	}

	isTaken, err := userService.CheckIfUsernameIsTaken(username)

	if err != nil || isTaken {
		return packets.ErrorCode_USERNAME_TAKEN, errors.New("username was taken")
	}

	return packets.ErrorCode_UNKOWN_ERROR, nil
}
