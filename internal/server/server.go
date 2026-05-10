package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/sudogane/project_timegate/internal/database"
	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/database/models"
	"github.com/sudogane/project_timegate/internal/logger"
	"github.com/sudogane/project_timegate/pkg/packets"
	"google.golang.org/protobuf/proto"
)

type GameServer struct {
	sessions map[string]*PlayerSession
	mutex    sync.RWMutex
	db       *database.Repository
	rdb      *cache.RedisClient
	logger   *logger.Logger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewGameServer(db *database.Repository, rdb *cache.RedisClient, l *logger.Logger) *GameServer {
	gs := &GameServer{
		sessions: make(map[string]*PlayerSession),
		db:       db,
		rdb:      rdb,
		logger:   l,
	}

	return gs
}

func (gs *GameServer) GetDB() *models.Queries {
	return gs.db.Queries
}

func (gs *GameServer) GetRDB() *cache.RedisClient {
	return gs.rdb
}

func (gs *GameServer) Ctx() context.Context {
	return context.Background()
}

func (gs *GameServer) AddSession(session *PlayerSession) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	if exists := gs.sessions[session.ID]; exists != nil {
		return
	}

	gs.sessions[session.ID] = session
}

func (gs *GameServer) RemoveSession(sessionId string) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()

	delete(gs.sessions, sessionId)
}

func (gs *GameServer) GetSession(sessionId string) *PlayerSession {
	gs.mutex.RLock()
	defer gs.mutex.RUnlock()

	session := gs.sessions[sessionId]
	return session
}

func (gs *GameServer) GetLogger() *logger.Logger {
	return gs.logger
}

func (gs *GameServer) HandleWebsocket(w http.ResponseWriter, r *http.Request, router RouterInterface) {
	gs.logger.Infow("New WebSocket connection Received")
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		gs.logger.Errorw("Error while upgrading WebSocket: " + err.Error())
		return
	}

	sessionId := uuid.NewString()
	session := &PlayerSession{
		ID:       sessionId,
		Conn:     conn,
		SendChan: make(chan *packets.FromServerToClient),
		Logger:   gs.logger.WithSession(sessionId, uuid.Nil),
	}

	gs.AddSession(session)
	go gs.ReadLoop(session, router)
	go gs.WriteLoop(session)

	idPacket := packets.NewWebsocketIdResponse(sessionId)
	gs.SendMessage(sessionId, idPacket)
}

func (gs *GameServer) ReadLoop(session *PlayerSession, router RouterInterface) {
	defer gs.RemoveSession(session.ID)

	for {
		messageType, data, err := session.Conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.BinaryMessage {
			fromClient := &packets.FromClientToServer{}
			if err := proto.Unmarshal(data, fromClient); err != nil {
				session.Logger.Errorw("Error unmarshaling proto" + err.Error())
				continue
			}

			router.Route(session, fromClient)

		}
	}

	session.Conn.Close()
}

func (gs *GameServer) WriteLoop(session *PlayerSession) {
	defer gs.RemoveSession(session.ID)

	for packet := range session.SendChan {
		writer, err := session.Conn.NextWriter(websocket.BinaryMessage)

		if err != nil {
			session.Logger.Errorw("Error creating next writer")
			return
		}

		data, err := proto.Marshal(packet)
		if err != nil {
			session.Logger.Errorw("Error marshaling proto")
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			session.Logger.Errorw("Error writing data")
			return
		}

		writer.Write([]byte{'\n'})

		if err := writer.Close(); err != nil {
			session.Logger.Errorw("Error closing writer")
			return
		}
	}
}

func (gs *GameServer) SendMessage(sessionId string, message packets.ServerPayload) {
	session := gs.sessions[sessionId]
	if session == nil {
		return
	}

	session.SendChan <- &packets.FromServerToClient{Payload: message}
}

func (gs *GameServer) SendErrorMessage(sessionId string, code packets.ErrorCode) {
	errorPacket := packets.NewErrorMessage(code)
	gs.SendMessage(sessionId, errorPacket)
}
