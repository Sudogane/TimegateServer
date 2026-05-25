package handler

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sudogane/project_timegate/internal/crypt"
	"github.com/sudogane/project_timegate/internal/database/models"
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
	if packet == nil {
		packet := msg.GetAuthenticationReconnectRequest()
		if packet == nil {
			return errors.New("invalid packet")
		}
	}

	gameVersion := msg.GetGameVersion()
	if gameVersion != "0.0.1TB" {
		h.SendError(session, packets.ErrorCode_INVALID_VERSION)
		return nil
	}

	if msg.GetPacketType() == packets.PacketType_AUTHENTICATION_RECONNECT_REQUEST {
		h.handleUserReconnect(session, msg.GetAuthenticationReconnectRequest())
		return nil
	}

	switch packet.GetType() {
	case packets.AuthenticationType_LOGIN:
		h.handleUserLogin(session, packet)
		return nil
	case packets.AuthenticationType_REGISTER:
		h.handleUserRegister(session, packet)
		return nil
	}

	return nil
}

func (h *AuthenticationHandler) handleUserLogin(session *server.PlayerSession, loginRequestData *packets.AuthenticationRequest) {
	user, err := h.validateLoginData(loginRequestData)
	if err != nil {
		h.SendError(session, packets.ErrorCode_INVALID_CREDENTIALS)
		return
	}

	resources, err := h.userService.GetUserWithResources(user.ID)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
		return
	}
	redirectDialogueId := ""
	starterFlag, err := h.flagsService.GetUserFlag(user.ID, "has_selected_starter")

	if err == pgx.ErrNoRows || err == nil && !starterFlag.IsActive.Bool {
		redirectDialogueId = "DEVELOPMENT"
	}

	session.PlayerId = user.ID
	session.Logger = h.server.GetLogger().WithSession(session.ID, session.PlayerId)
	userDataPacket := &packets.UserData{
		Username:       user.Username,
		Level:          resources.Level.Int32,
		Exp:            resources.Exp.Int32,
		Bits:           resources.Bits.Int64,
		Yen:            resources.Yen.Int64,
		StaminaCurrent: resources.StaminaCurrent.Int32,
		StaminaMax:     resources.StaminaMax.Int32,
	}

	token, err := crypt.CreateToken(user.ID.String())
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
		return
	}

	responsePacket := packets.NewAuthenticationResponse(token, userDataPacket, redirectDialogueId, user.ID.String())
	h.Send(session, responsePacket)

}

func (h *AuthenticationHandler) validateLoginData(loginData *packets.AuthenticationRequest) (*models.User, error) {
	username := loginData.GetUsername()
	password := loginData.GetPassword()

	if username == "" || password == "" {
		return nil, errors.New("empty data")
	}

	if len(username) < 3 || len(username) > 10 {
		return nil, errors.New("length of username is too short or too long")
	}

	if len(password) < 4 {
		return nil, errors.New("length of password is too short or too long")
	}

	invalidCredentialsError := errors.New("invalid credentials")

	user, doesUsernameExist := h.userService.GetByUsername(username)
	if doesUsernameExist != nil {
		return nil, invalidCredentialsError
	}

	if isPasswordCorrect := crypt.ComparePassword(password, user.PasswordHash); !isPasswordCorrect {
		return nil, invalidCredentialsError
	}

	return user, nil
}

func (h *AuthenticationHandler) handleUserRegister(session *server.PlayerSession, registerRequestData *packets.AuthenticationRequest) {
	if code, err := h.validateRegisterData(registerRequestData); err != nil {
		h.SendError(session, code)
		return
	}

	hashedPassword, err := crypt.HashPassword(registerRequestData.GetPassword())
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
		return
	}
	user, err := h.userService.CreateUserWithResources(registerRequestData.GetUsername(), hashedPassword)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
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
	session.Logger = h.server.GetLogger().WithSession(session.ID, user.ID)

	token, err := crypt.CreateToken(user.ID.String())
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
		return
	}

	responsePacket := packets.NewAuthenticationResponse(token, userData, "DEVELOPMENT", user.ID.String())
	h.Send(session, responsePacket)
}

func (h *AuthenticationHandler) validateRegisterData(registerData *packets.AuthenticationRequest) (packets.ErrorCode, error) {
	username := registerData.GetUsername()
	password := registerData.GetPassword()

	if len(username) < 3 || len(username) > 10 {
		return packets.ErrorCode_USERNAME_TOO_SHORT_OR_LONG, errors.New("username too long or too short")
	}

	if len(password) < 4 {
		return packets.ErrorCode_PASSWORD_TOO_SHORT, errors.New("password is too short")
	}

	isTaken, err := h.userService.CheckIfUsernameIsTaken(username)

	if err != nil || isTaken {
		return packets.ErrorCode_USERNAME_TAKEN, errors.New("username was taken")
	}

	return packets.ErrorCode_UNKOWN_ERROR, nil
}

func (h *AuthenticationHandler) handleUserReconnect(session *server.PlayerSession, reconnectRequestData *packets.AuthenticationReconnectRequest) {
	token := reconnectRequestData.GetToken()
	if token == "" {
		h.SendError(session, packets.ErrorCode_INVALID_CREDENTIALS)
		return
	}

	claims, err := crypt.VerifyToken(token)
	if err != nil {
		h.SendError(session, packets.ErrorCode_INVALID_CREDENTIALS)
		return
	}

	playerId, err := uuid.Parse(claims.UserId)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		return
	}

	user, err := h.userService.GetById(playerId)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		return
	}

	resources, err := h.userService.GetUserWithResources(user.ID)
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
		return
	}
	redirectDialogueId := ""
	starterFlag, err := h.flagsService.GetUserFlag(user.ID, "has_selected_starter")
	if err == pgx.ErrNoRows || err == nil && !starterFlag.IsActive.Bool {
		redirectDialogueId = "DEVELOPMENT"
	}

	session.PlayerId = user.ID
	session.Logger = h.server.GetLogger().WithSession(session.ID, session.PlayerId)
	userDataPacket := &packets.UserData{
		Username:       user.Username,
		Level:          resources.Level.Int32,
		Exp:            resources.Exp.Int32,
		Bits:           resources.Bits.Int64,
		Yen:            resources.Yen.Int64,
		StaminaCurrent: resources.StaminaCurrent.Int32,
		StaminaMax:     resources.StaminaMax.Int32,
	}

	token, err = crypt.CreateToken(user.ID.String())
	if err != nil {
		h.SendError(session, packets.ErrorCode_UNKOWN_ERROR)
		session.Logger.Errorw(err.Error())
		return
	}

	responsePacket := packets.NewAuthenticationResponse(token, userDataPacket, redirectDialogueId, user.ID.String())
	h.Send(session, responsePacket)
}
