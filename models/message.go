package models

import "time"

type Message struct {
    ID          string    `json:"id"`
    SenderID    string    `json:"senderId"`
    ReceiverID  string    `json:"receiverId"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"createdAt"`
}

type ChatUser struct {
    User         User      `json:"user"`
    LastMessage  *Message  `json:"lastMessage,omitempty"`
    UnreadCount  int       `json:"unreadCount"`
}