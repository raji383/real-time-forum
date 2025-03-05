package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LikeRequest struct {
	PostID int    `json:"post_id"`
	Action string `json:"action"`
	Change int    `json:"change"`
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed !", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		sendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
		})
		return
	}

	ex := CheckCookie(cookie.Value)
	if !ex {
		sendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request !", http.StatusBadRequest)
		return
	}

	var req LikeRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Bad Request !", http.StatusBadRequest)
		return
	}

	var column string
	switch req.Action {
	case "like":
		column = "likes"
	case "dislike":
		column = "dislikes"
	default:
		http.Error(w, "Bad Request !!", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE posts SET "+column+" = "+column+" + ? WHERE id = ?", req.Change, req.PostID)
	if err != nil {
		fmt.Println("errrrrr", err)
		http.Error(w, "Internal Server Error !!", http.StatusInternalServerError)
		return
	}

	if req.Change != -1 {
		err = LikedPost(cookie.Value, req.PostID, req.Action)
		if err != nil {
			fmt.Println("err", err)
			http.Error(w, "Internal Server Error !!", http.StatusInternalServerError)
			return
		}
	}
	
	sendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})

}

func CheckCookie(cookieValue string) bool {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE session_token = ?", cookieValue).Scan(&userID)
	if err != nil {
		fmt.Println("Error fetching user ID:", err)
		return false
	}
	return true
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func LikedPost(cookieValue string, postID int, likeType string) error {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE session_token = ?", cookieValue).Scan(&userID)
	if err != nil {
		fmt.Println("Error fetching user ID:", err)
		return err
	}

	var n int
	if likeType == "like" {
		n = 1
	} else {
		n = -1
	}

	fmt.Println(n, likeType)

	_, err = db.Exec(`
        INSERT INTO post_likes (user_id, post_id, like_type) 
        VALUES (?, ?, ?) 
        ON CONFLICT(user_id, post_id) 
        DO UPDATE SET like_type = excluded.like_type;
    `, userID, postID, n)

	if err != nil {
		fmt.Println("Error inserting/updating post like:", err)
		return err
	}

	return nil
}

func GetUserReaction(r *http.Request, postID int) int {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return 0 
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err == sql.ErrNoRows || err != nil {
		return 0
	}

	
	var likeType int
	err = db.QueryRow("SELECT like_type FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeType)
	if err == sql.ErrNoRows {
		return 0 
	}
	if err != nil {
		return 0
	}

	return likeType
}

