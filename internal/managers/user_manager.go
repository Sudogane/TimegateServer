package managers

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sudogane/project_timegate/internal/models"
)

type UserManager struct {
	users       map[string]*models.User
	connections map[*websocket.Conn]string
	mutex       sync.RWMutex
}

func NewUserManager() *UserManager {
	return &UserManager{
		users:       make(map[string]*models.User),
		connections: make(map[*websocket.Conn]string),
	}
}

func (um *UserManager) Authenticate(username, password string) (*models.User, error) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	for _, user := range um.users {
		if user.Username == "123" && user.Password == "123" {
			return user, nil
		}

		return nil, fmt.Errorf("wrong credentials")
	}

	user := &models.User{
		ID:          "123456",
		Username:    "123",
		Password:    "123",
		AccessToken: "123456",
	}

	um.users[user.ID] = user

	return user, nil
}

func (um *UserManager) GetUserByConnection(conn *websocket.Conn) *models.User {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	userID, exists := um.connections[conn]

	if !exists {
		return nil
	}

	return um.users[userID]
}

func (um *UserManager) AddConnection(user *models.User, conn *websocket.Conn) {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	if oldConn := user.Conn; oldConn != nil {
		delete(um.connections, oldConn)
		oldConn.Close()
	}

	user.Conn = conn
	um.connections[conn] = user.ID
}

func (um *UserManager) GetPlayerCount() int {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	return len(um.users)
}
