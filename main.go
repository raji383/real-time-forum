package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
"forum/models"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.Initdb()
	http.HandleFunc("/", handlers.ForumHandler)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/regist", handlers.Regist)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/like", models.LikeHandler)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/unlike", models.UnlikeHandler)
	http.HandleFunc("/comment", models.CommentHandler)
	http.HandleFunc("/likes/count", models.LikesCountHandler)
	http.HandleFunc("/comments", models.CommentsHandler)
	http.HandleFunc("/static/", handlers.StaticHnadler)
	
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
