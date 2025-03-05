package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	forum "forum/func"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./databases/data.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStatements, err := os.ReadFile("./databases/my.sql")
	if err != nil {
		log.Fatal("Error reading SQL file:", err)
	}
	_, err = db.Exec(string(sqlStatements))
	if err != nil {
		log.Fatal("Error executing SQL statements:", err)
	}

	forum.InitHandlers(db)
}

func main() {
	defer db.Close()

	http.HandleFunc("/", forum.HomePage)
	http.HandleFunc("/login-page", forum.LoginPage)
	http.HandleFunc("/singup-page", forum.SingupPage)
	http.HandleFunc("/singup", forum.Singup)
	http.HandleFunc("/login", forum.Login)
	http.HandleFunc("/logout", forum.Logout)
	http.HandleFunc("/posts", forum.Posts)
	http.HandleFunc("/filter", forum.Filter)
	http.HandleFunc("/like", forum.LikeHandler)
	http.HandleFunc("/create-Post", forum.CreatePost)
	http.HandleFunc("/liked-posts", forum.LikedPostsHandler)
	http.HandleFunc("/comment", forum.CommentHandler)

	http.HandleFunc("/static/", forum.StaticHandle)


	fmt.Println("http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
