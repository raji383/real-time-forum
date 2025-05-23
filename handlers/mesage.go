package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"real-time-forum/databases"
)

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Recipient string `json:"recipient"`
		Content   string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Received message:", req)
	// Get sender from session token instead of cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	senderID, exists := GetUserBySession(cookie.Value)
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get receiver ID from nickname
	var receiverID int
	err = databases.DB.QueryRow("SELECT id FROM users WHERE nickname = ?", req.Recipient).Scan(&receiverID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false, "error":"Invalid recipient"}`))
		return
	}

	// Store message with proper user IDs
	_, err = databases.DB.Exec(
		`INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)`,
		senderID, receiverID, req.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false, "error":"Database error"}`))
		return
	}

	// After storing in database, broadcast through WebSocket
	hub.SendPrivateMessage(GetUserInfo(senderID).Username, req.Recipient, req.Content)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	recipient := r.URL.Query().Get("recipient")

	// Get sender from session token
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	senderID, exists := GetUserBySession(cookie.Value)
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get sender's username
	//sender := GetUserInfo(senderID).Username

	// Get recipient's ID
	var recipientID int
	err = databases.DB.QueryRow("SELECT id FROM users WHERE nickname = ?", recipient).Scan(&recipientID)
	if err != nil {
		http.Error(w, "Invalid recipient", http.StatusBadRequest)
		return
	}

	rows, err := databases.DB.Query(`
        SELECT m.sender_id, m.receiver_id, m.content, m.timestamp, 
               s.nickname as sender_name, r.nickname as receiver_name
        FROM messages m
        JOIN users s ON m.sender_id = s.id
        JOIN users r ON m.receiver_id = r.id
        WHERE (m.sender_id = ? AND m.receiver_id = ?) 
           OR (m.sender_id = ? AND m.receiver_id = ?)
        ORDER BY m.timestamp ASC
    `, senderID, recipientID, recipientID, senderID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var senderID, receiverID int
		var content, timestamp, senderName, receiverName string
		if err := rows.Scan(&senderID, &receiverID, &content, &timestamp, &senderName, &receiverName); err != nil {
			continue
		}
		messages = append(messages, map[string]interface{}{
			"sender":    senderName,
			"receiver":  receiverName,
			"content":   content,
			"timestamp": timestamp,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
