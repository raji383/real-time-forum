package models

import (
	"encoding/json"
	"net/http"
	"forum/database"

)
func addLike(userID, postID int) error {
	_, err := database.Db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	return err
}
func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := addLike(1, 1)
	if err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}
	
}
func getLikesCount() int {
	var count int
	err := database.Db.QueryRow("SELECT COUNT(*) FROM likes").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}
func LikesCountHandler(w http.ResponseWriter, r *http.Request) {
	count := getLikesCount()
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}
func removeLike(userID, postID int) error {
	_, err := database.Db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
	return err
}
func UnlikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := removeLike(1, 1)
	if err != nil {
		http.Error(w, "Failed to unlike post", http.StatusInternalServerError)
		return
	}

}