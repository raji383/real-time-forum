package forum

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type gg struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request !", http.StatusBadRequest)
		return
	}

	var req gg
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Bad Request !", http.StatusBadRequest)
		return
	}

	var userID int
	var hashedPassword string
	err = db.QueryRow("SELECT id, password FROM users WHERE username = ?", req.Username).Scan(&userID, &hashedPassword)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Username is incorrect!"})
		return
	} else if err != nil {
		http.Error(w, "Error in Database!", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password is incorrect!"})
		return
	}

	sessionToken := uuid.New().String()
	expiration := time.Now().Add(1 * time.Hour)

	_, err = db.Exec("UPDATE users SET session_token = ? WHERE id = ?", sessionToken, userID)
	if err != nil {
		http.Error(w, "Error saving session!", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiration,
		HttpOnly: true,
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]string{"success": "Login successful!"})
}

func IsLoggedIn(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return false, ""
	}

	var userName string
	err = db.QueryRow("SELECT username FROM users WHERE session_token = ?", cookie.Value).Scan(&userName)
	if err == sql.ErrNoRows || err != nil {
		return false, ""
	}

	return true, userName
}
