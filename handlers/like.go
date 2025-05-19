package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/databases"
)

type LikeData struct {
	Reaction string `json:"reaction_type"`
	PostID   string `json:"post_id"`
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var pd LikeData
	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query1 := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > DATETIME('now')`
	var userID int
	err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the user already liked the post
	checkQuery := `
		SELECT COUNT(*) FROM post_reactions 
		WHERE post_id = ? AND user_id = ? AND reaction_type = ?
	`
	var count int
	err = databases.DB.QueryRow(checkQuery, pd.PostID, userID, pd.Reaction).Scan(&count)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Already liked this post",
		})
		return
	}

	// Insert the like
	insertQuery := `
		INSERT INTO post_reactions (post_id, user_id, reaction_type)
		VALUES (?, ?, ?)
	`
	_, err = databases.DB.Exec(insertQuery, pd.PostID, userID, pd.Reaction)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Get the updated like count
	countQuery := `
		SELECT COUNT(*) FROM post_reactions 
		WHERE post_id = ? AND reaction_type = ?
	`
	var likeCount int
	err = databases.DB.QueryRow(countQuery, pd.PostID, pd.Reaction).Scan(&likeCount)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"like_count": likeCount,
	})
}
