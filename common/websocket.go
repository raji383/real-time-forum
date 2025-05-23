package common

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID   string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	hub      *Hub
}

type Message struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	From      string `json:"from"`
	To        string `json:"to"`
	Timestamp string `json:"timestamp"`
}

type Hub struct {
	Clients     map[*Client]bool
	UserMap     map[string]*Client
	Broadcast   chan []byte
	Register    chan *Client
	Unregister  chan *Client
	mu          sync.Mutex
	onlineUsers map[string]bool
}

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		UserMap:     make(map[string]*Client),
		Broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		onlineUsers: make(map[string]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.UserMap[client.Username] = client
			h.onlineUsers[client.Username] = true
			h.broadcastUserList()
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.UserMap, client.Username)
				delete(h.onlineUsers, client.Username)
				close(client.Send)
				h.broadcastUserList()
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
					delete(h.UserMap, client.Username)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) broadcastUserList() {
	users := make([]map[string]interface{}, 0)
	h.mu.Lock()
	for username := range h.UserMap {
		users = append(users, map[string]interface{}{
			"username": username,
			"online":   h.onlineUsers[username],
		})
	}
	h.mu.Unlock()

	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Printf("Error marshalling users list: %v", err)
		return
	}

	message := Message{
		Type:    "userList",
		Content: string(usersJSON),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	for client := range h.Clients {
		client.Send <- data
	}
}

func (h *Hub) SendPrivateMessage(from, to string, content string) {
	message := Message{
		Type:    "privateMessage",
		From:    from,
		To:      to,
		Content: content,
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling private message: %v", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if toClient, ok := h.UserMap[to]; ok {
		toClient.Send <- data
	}
	if fromClient, ok := h.UserMap[from]; ok {
		fromClient.Send <- data
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		switch msg.Type {
		case "privateMessage":
			c.hub.SendPrivateMessage(c.Username, msg.To, msg.Content)
		default:
			c.hub.Broadcast <- rawMessage
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewClient(userID string, username string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		UserID:   userID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		hub:      hub,
	}
}
