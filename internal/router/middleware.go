package router

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sudogane/project_timegate/internal/crypt"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/pkg/packets"
)

func (r *Router) authenticate(session *server.PlayerSession, msg *packets.FromClientToServer) error {
	if msg.PacketType == packets.PacketType_AUTHENTICATION_REQUEST || msg.PacketType == packets.PacketType_AUTHENTICATION_RECONNECT_REQUEST {
		return nil
	}

	token := msg.GetAccessToken()
	if token == "" {
		r.server.SendErrorMessage(session.ID, packets.ErrorCode_UNKOWN_ERROR)
		return errors.New("missing access token")
	}

	claims, err := crypt.VerifyToken(token)
	if err != nil {
		r.server.SendErrorMessage(session.ID, packets.ErrorCode_UNKOWN_ERROR)
		return errors.New("invalid or expired token")
	}

	if claims.UserId != msg.GetPlayerId() {
		r.server.SendErrorMessage(session.ID, packets.ErrorCode_UNKOWN_ERROR)
		return errors.New("user id mismatch")
	}

	session.PlayerId, err = uuid.Parse(claims.UserId)
	if err != nil {
		r.server.SendErrorMessage(session.ID, packets.ErrorCode_UNKOWN_ERROR)
		return errors.New("error parsing uuid")
	}
	return nil
}
