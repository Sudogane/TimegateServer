package router

import (
	"fmt"
	"sync"

	"github.com/sudogane/project_timegate/internal/router/handler"
	"github.com/sudogane/project_timegate/internal/server"
	"github.com/sudogane/project_timegate/pkg/packets"
)

type Router struct {
	handlers map[packets.PacketType]handler.Handler
	mutex    sync.RWMutex
	server   server.GameServerInterface
}

func NewRouter(server server.GameServerInterface) *Router {
	router := &Router{
		handlers: make(map[packets.PacketType]handler.Handler),
		server:   server,
	}

	router.RegisterRoutes()
	return router
}

func (r *Router) RegisterRouter(packetType packets.PacketType, handler handler.Handler) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.handlers[packetType] = handler
}

func (r *Router) Route(session *server.PlayerSession, msg *packets.FromClientToServer) {
	r.mutex.RLock()
	handler, exists := r.handlers[msg.PacketType]
	r.mutex.RUnlock()

	if !exists {
		session.Log("ERROR", "No handler found for packet type: "+msg.PacketType.String())
		return
	}

	err := handler.Handle(session, msg)
	if err != nil {
		session.Log("ERROR", fmt.Sprintf("Error handling packet type: %s, error: %s", msg.PacketType.String(), err.Error()))
		return
	}
}

func (r *Router) RegisterRoutes() {
	r.RegisterRouter(packets.PacketType_AUTHENTICATION_REQUEST, handler.NewAuthenticationHandler(r.server))

	// --> Stages ::
	stagesHandler := handler.NewStagesHandler(r.server)
	r.RegisterRouter(packets.PacketType_CHAPTER_DATA_REQUEST, stagesHandler)
	r.RegisterRouter(packets.PacketType_EPISODES_BY_CHAPTER_REQUEST, stagesHandler)
}
