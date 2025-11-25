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
		fmt.Println("No handler for packet type: ", msg.PacketType)
		return
	}

	err := handler.Handle(session, msg)
	if err != nil {
		fmt.Println("Error handling packet: ", err, " of packet type: ", msg.PacketType)
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
