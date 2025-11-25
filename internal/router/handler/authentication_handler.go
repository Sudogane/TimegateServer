package handler

import (
	"errors"

	"github.com/sudogane/project_timegate/internal/crypt"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/internal/services"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type AuthenticationHandler struct {
	BaseHandler
	userService *services.UserService
}

func NewAuthenticationHandler(server server.GameServerInterface) *AuthenticationHandler {
	return &AuthenticationHandler{
		BaseHandler: *NewBaseHandler(server),
		userService: services.NewUserService(server),
	}
}

func (h *AuthenticationHandler) Handle(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	packet := msg.GetAuthenticationRequest()

	switch packet.GetType() {
	case packets.AuthenticationType_LOGIN:
		onUserLogin(h, session, packet)
	case packets.AuthenticationType_REGISTER:
		onUserRegister(h, session, packet)
	}

	return nil
}

func onUserLogin(h *AuthenticationHandler, session *server.PlayerSession, loginRequestData *packets.AuthenticationRequest) {
	if err := validateLoginData(loginRequestData, h.userService); err != nil {
		h.SendError(session, packets.ErrorCode_INVALID_CREDENTIALS)
		return
	}

	user, _ := h.userService.GetByUsername(loginRequestData.Username)
	resources, _ := h.userService.GetUserWithResources(user.ID)

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
	responsePacket := packets.NewAuthenticationResponse("user access token", userDataPacket)
	h.Send(session, responsePacket)
}

func validateLoginData(loginData *packets.AuthenticationRequest, userService *services.UserService) error {
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

func onUserRegister(h *AuthenticationHandler, session *server.PlayerSession, registerRequestData *packets.AuthenticationRequest) {
	if code, err := validateRegisterData(registerRequestData, h.userService); err != nil {
		h.SendError(session, code)
		return
	}

	hashedPassword, _ := crypt.HashPassword(registerRequestData.GetPassword())
	user, err := h.userService.CreateUserWithResources(registerRequestData.GetUsername(), hashedPassword)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
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

	responsePacket := packets.NewAuthenticationResponse("user access token", userData)
	h.Send(session, responsePacket)
}

func validateRegisterData(registerData *packets.AuthenticationRequest, userService *services.UserService) (packets.ErrorCode, error) {
	username := registerData.GetUsername()
	password := registerData.GetPassword()

	if len(username) < 3 || len(username) > 10 {
		return packets.ErrorCode_USERNAME_TOO_SHORT_OR_LONG, errors.New("username too long or too short")
	}

	if len(password) < 3 {
		return packets.ErrorCode_PASSWORD_TOO_SHORT, errors.New("username too long or too short")
	}

	_, usernameWasTaken := userService.GetByUsername(username)

	if usernameWasTaken == nil {
		return packets.ErrorCode_USERNAME_TAKEN, errors.New("username too long or too short")
	}

	return packets.ErrorCode_UNKOWN_ERROR, nil
}
