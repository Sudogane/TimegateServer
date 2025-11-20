package models

import "github.com/gorilla/websocket"

type User struct {
	ID          string
	Email       string
	Username    string
	Password    string
	AccessToken string
	Conn        *websocket.Conn
}
