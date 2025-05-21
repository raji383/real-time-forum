package models

import "time"

type Post struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Category  string    `json:"category"`
    CreatedAt time.Time `json:"createdAt"`
}

type Comment struct {
    ID        string    `json:"id"`
    PostID    string    `json:"postId"`
    UserID    string    `json:"userId"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"createdAt"`
}