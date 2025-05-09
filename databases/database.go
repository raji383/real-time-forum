package databases

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3" // Import the driver anonymously

)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "./databases/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	creat := `CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
	    nickname TEXT NOT NULL,
    	age INTEGER,
	    gender TEXT,
    	first_name TEXT,
    	last_name TEXT,
    	email TEXT UNIQUE,
    	password TEXT
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP

	);
	CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			session_token TEXT,
			expires_at DATETIME,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);

	CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	content TEXT,
	title TEXT,
	interest TEXT,
	photo TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user_id) REFERENCES users(id)
);


		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS post_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			reaction_type TEXT NOT NULL CHECK(reaction_type IN ('like', 'dislike')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(post_id, user_id)
		);

		CREATE TABLE IF NOT EXISTS comment_reactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			comment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			reaction_type TEXT NOT NULL CHECK(reaction_type IN ('like', 'dislike')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(comment_id, user_id)
		);
`
	_, err = DB.Exec(creat)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database initialized successfully")
}

/*func CloseDB() {
	err := DB.Close()
	fmt.Println("closing the databse")
	if err != nil {
		log.Fatal("Error closing database:", err)
	}
}*/
