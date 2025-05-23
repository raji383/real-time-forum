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
        SELECT 
            p.id,
            p.title, 
            p.content, 
            p.interest, 
            p.user_id, 
            p.created_at, 
            u.nickname,
            (SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'like') as likes,
            (SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'dislike') as dislikes
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
		var id, userID, likes, dislikes int
		var createdAt string

		if err := rows.Scan(&id, &title, &content, &interest, &userID, &createdAt, &nickname, &likes, &dislikes); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		post := map[string]interface{}{
			"id":         id,
			"title":      title,
			"content":    content,
			"topics":     strings.Split(interest, ","),
			"author":     nickname,
			"created_at": createdAt,
			"likes":      likes,
			"dislikes":   dislikes,
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// handlers/comments.go

type CommentData struct {
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var cd CommentData
	if err := json.NewDecoder(r.Body).Decode(&cd); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	var userID int
	err = databases.DB.QueryRow(`SELECT user_id FROM sessions WHERE session_token = ?`, cookie.Value).Scan(&userID)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	// Insert comment
	_, err = databases.DB.Exec(`
		INSERT INTO comments (post_id, user_id, content)
		VALUES (?, ?, ?)
	`, cd.PostID, userID, cd.Content)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment created successfully",
	})
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
		return
	}

	rows, err := databases.DB.Query(`
        SELECT 
            c.id,
            c.content,
            c.created_at,
            u.nickname,
            (SELECT COUNT(*) FROM comment_reactions WHERE comment_id = c.id AND reaction_type = 'like') as likes,
            (SELECT COUNT(*) FROM comment_reactions WHERE comment_id = c.id AND reaction_type = 'dislike') as dislikes
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC
    `, postID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
		return
	}
	defer rows.Close()

	comments := make([]map[string]interface{}, 0) // Initialize as empty array
	for rows.Next() {
		var id, likes, dislikes int
		var content, createdAt, nickname string

		if err := rows.Scan(&id, &content, &createdAt, &nickname, &likes, &dislikes); err != nil {
			continue
		}

		comment := map[string]interface{}{
			"id":         id,
			"content":    content,
			"author":     nickname,
			"created_at": createdAt,
			"likes":      likes,
			"dislikes":   dislikes,
		}
		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(comments); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]interface{}{})
	}
}

func CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type ReactionData struct {
		CommentID    int    `json:"comment_id"`
		ReactionType string `json:"reaction_type"`
	}

	var rd ReactionData
	if err := json.NewDecoder(r.Body).Decode(&rd); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get user ID from session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	var userID int
	err = databases.DB.QueryRow(`SELECT user_id FROM sessions WHERE session_token = ?`, cookie.Value).Scan(&userID)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	// Check if reaction exists
	var existingReaction string
	err = databases.DB.QueryRow(`
		SELECT reaction_type FROM comment_reactions 
		WHERE comment_id = ? AND user_id = ?
	`, rd.CommentID, userID).Scan(&existingReaction)

	if err == nil {
		// Reaction exists - update or remove
		if existingReaction == rd.ReactionType {
			// Remove reaction
			_, err = databases.DB.Exec(`
				DELETE FROM comment_reactions 
				WHERE comment_id = ? AND user_id = ?
			`, rd.CommentID, userID)
		} else {
			// Update reaction
			_, err = databases.DB.Exec(`
				UPDATE comment_reactions 
				SET reaction_type = ?
				WHERE comment_id = ? AND user_id = ?
			`, rd.ReactionType, rd.CommentID, userID)
		}
	} else {
		// No reaction exists - insert new
		_, err = databases.DB.Exec(`
			INSERT INTO comment_reactions (comment_id, user_id, reaction_type)
			VALUES (?, ?, ?)
		`, rd.CommentID, userID, rd.ReactionType)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get updated counts
	var likes, dislikes int
	err = databases.DB.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND reaction_type = 'like') as likes,
			(SELECT COUNT(*) FROM comment_reactions WHERE comment_id = ? AND reaction_type = 'dislike') as dislikes
	`, rd.CommentID, rd.CommentID).Scan(&likes, &dislikes)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"like_count":    likes,
		"dislike_count": dislikes,
	})
}
