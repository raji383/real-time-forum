package handlers

import (
	"log"
	"net/http"

	"real-time-forum/common"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub = common.NewHub()

func init() {
	go hub.Run()
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Error getting session: %v", err)
		conn.Close()
		return
	}

	userID, exists := GetUserBySession(cookie.Value)
	if !exists {
		conn.Close()
		return
	}

	user := GetUserInfo(userID)

	client := common.NewClient(string(userID), user.Username, conn, hub)
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
