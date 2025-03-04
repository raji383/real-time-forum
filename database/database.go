package database

import (
	"database/sql"
	"log"
)

var Db *sql.DB

func Initdb() {
	var err error
	Db, err = sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    username TEXT NOT NULL UNIQUE,
	    email TEXT NOT NULL UNIQUE,
	    password TEXT NOT NULL
	);
	`
	_, err = Db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}


	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id INTEGER NOT NULL,
	    title TEXT NOT NULL,
	    content TEXT NOT NULL,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`
	_, err = Db.Exec(createPostsTable)
	if err != nil {
		log.Fatal("Error creating posts table:", err)
	}

	createLikesTable := `
	CREATE TABLE IF NOT EXISTS likes (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id INTEGER NOT NULL,
	    post_id INTEGER NOT NULL,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);
	`
	_, err = Db.Exec(createLikesTable)
	if err != nil {
		log.Fatal("Error creating likes table:", err)
	}

	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    user_id INTEGER NOT NULL,
	    post_id INTEGER NOT NULL,
	    content TEXT NOT NULL,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);
	`
	_, err = Db.Exec(createCommentsTable)
	if err != nil {
		log.Fatal("Error creating comments table:", err)
	}
}
