package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "real-time-forum/databases"
)

type LikeData struct {
    Reaction string `json:"reaction_type"`
    PostID   int    `json:"post_id"` // Changed from string to int
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

    // Validate reaction type
    if pd.Reaction != "like" && pd.Reaction != "dislike" {
        http.Error(w, "Invalid reaction type", http.StatusBadRequest)
        return
    }

    cookie, err := r.Cookie("session_token")
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Get user ID from session
    query1 := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > DATETIME('now')`
    var userID int
    err = databases.DB.QueryRow(query1, cookie.Value).Scan(&userID)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Check existing reaction
    checkQuery := `
        SELECT reaction_type FROM post_reactions 
        WHERE post_id = ? AND user_id = ?
    `
    var existingReaction string
    err = databases.DB.QueryRow(checkQuery, pd.PostID, userID).Scan(&existingReaction)
    if err != nil && err.Error() != "sql: no rows in result set" {
        log.Println(err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Begin transaction
    tx, err := databases.DB.Begin()
    if err != nil {
        log.Println(err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    if existingReaction == pd.Reaction {
        // Same reaction: remove it
        deleteQuery := `
            DELETE FROM post_reactions 
            WHERE post_id = ? AND user_id = ? AND reaction_type = ?
        `
        _, err = tx.Exec(deleteQuery, pd.PostID, userID, pd.Reaction)
        if err != nil {
            log.Println(err)
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
    } else if existingReaction != "" {
        // Different reaction: update it
        updateQuery := `
            UPDATE post_reactions 
            SET reaction_type = ?, created_at = CURRENT_TIMESTAMP 
            WHERE post_id = ? AND user_id = ?
        `
        _, err = tx.Exec(updateQuery, pd.Reaction, pd.PostID, userID)
        if err != nil {
            log.Println(err)
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
    } else {
        // No reaction: insert new
        insertQuery := `
            INSERT INTO post_reactions (post_id, user_id, reaction_type)
            VALUES (?, ?, ?)
        `
        _, err = tx.Exec(insertQuery, pd.PostID, userID, pd.Reaction)
        if err != nil {
            log.Println(err)
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
    }

    // Get updated counts
    countQuery := `
        SELECT 
            (SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND reaction_type = 'like') as like_count,
            (SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND reaction_type = 'dislike') as dislike_count
    `
    var likeCount, dislikeCount int
    err = tx.QueryRow(countQuery, pd.PostID, pd.PostID).Scan(&likeCount, &dislikeCount)
    if err != nil {
        log.Println(err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Commit transaction
    if err = tx.Commit(); err != nil {
        log.Println(err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    // Return updated counts
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]int{
        "like_count":    likeCount,
        "dislike_count": dislikeCount,
    })
}