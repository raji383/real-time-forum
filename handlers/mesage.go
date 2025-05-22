package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	// Get sender from session cookie
	sender := "anonymous"
	cookie, err := r.Cookie("session_user")
	if err == nil {
		sender = cookie.Value
	}

	_, err = databases.DB.Exec(`INSERT INTO messages (sender, recipient, content) VALUES (?, ?, ?)`, sender, req.Recipient, req.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false, "error":"Database error"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	recipient := r.URL.Query().Get("recipient")
	sender := "anonymous"
	cookie, err := r.Cookie("session_user")
	if err == nil {
		sender = cookie.Value
	}
	rows, err := databases.DB.Query(`SELECT sender, recipient, content, timestamp FROM messages WHERE (sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?) ORDER BY timestamp ASC`, sender, recipient, recipient, sender)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false, "error":"Database error"}`))
		return
	}
	defer rows.Close()
	var chat []map[string]interface{}
	for rows.Next() {
		var senderVal, recipientVal, content string
		var timestamp time.Time
		if err := rows.Scan(&senderVal, &recipientVal, &content, &timestamp); err == nil {
			chat = append(chat, map[string]interface{}{
				"sender":    senderVal,
				"recipient": recipientVal,
				"content":   content,
				"timestamp": timestamp,
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}
