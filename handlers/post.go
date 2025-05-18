package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"real-time-forum/databases"
)

type PostData struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Topics      []string `json:"topics"`
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var pd PostData
	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}
	query1 := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	query := `
    INSERT INTO posts (title, content, interest, user_id)
    VALUES (?, ?, ?, ?)
`
	_, err = databases.DB.Exec(query, pd.Title, pd.Description, strings.Join(pd.Topics, ","), userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Data received successfully",
		"title":    pd.Title,
		"content":  pd.Description,
		"interest": strings.Join(pd.Topics, ","),
	})
}

func ApiPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := databases.DB.Query(`
        SELECT p.title, p.content, p.interest, p.user_id, p.created_at, u.nickname 
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
    `)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var title, content, interest, nickname string
		var userID int
		var createdAt string

		if err := rows.Scan(&title, &content, &interest, &userID, &createdAt, &nickname); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		post := map[string]interface{}{
			"title":      title,
			"content":    content,
			"topics":     strings.Split(interest, ","),
			"author":     nickname,
			"created_at": createdAt,
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
