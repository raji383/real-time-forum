package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"real-time-forum/databases"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	username := r.FormValue("nickname")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	age := r.FormValue("age")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	password := r.FormValue("password")
	password2 := r.FormValue("confirm_password")
	fmt.Println(username)
	if password != password2 {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	// INSERT into database
	_, err = databases.DB.Exec(`
		INSERT INTO users (nickname, age, gender, first_name, last_name, email, password)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		username, age, gender, firstName, lastName, email, hashedPassword)
	if err != nil {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUserInfo(user_id int) User {
	var user User

	query := `SELECT nickname FROM users WHERE id = ?`
	err := databases.DB.QueryRow(query, user_id).Scan(&user.Username)
	if err != nil {
		log.Printf("Error retrieving user info: %v", err)
		return User{} // Return an empty User struct on error
	}

	return user
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	hasSession := DeleteUserBySession(sessionCookie.Value)

	// Remove session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	if !hasSession {
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/index.html")
}

func GetUserBySession(session_id string) (int, bool) {
	var user_id int
	query := `SELECT user_id FROM sessions WHERE session_token = ?`
	err := databases.DB.QueryRow(query, session_id).Scan(&user_id)
	if err != nil {
		log.Println(err)
		return user_id, false
	}
	return user_id, true
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	var userID int
	query := `SELECT user_id FROM sessions WHERE session_token = ? AND expires_at > DATETIME('now')`
	err = databases.DB.QueryRow(query, cookie.Value).Scan(&userID)
	if err != nil {
		w.Write([]byte(`{"loggedIn": false}`))
		return
	}

	// Get user info
	username := GetUserInfo(userID)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"loggedIn": true, "username": "%s"}`, username.Username)))
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	nickname := r.FormValue("user")
	password := r.FormValue("password")

	// Query for hashed password
	var hashedPassword string
	err := databases.DB.QueryRow(`SELECT password FROM users WHERE nickname = ?`, nickname).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Set session token
	user_id, _ := GetUserHash(nickname)
	session := GenerateToken(32)

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})

	// Set session in database
	SetSessionToken(user_id, session)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Login successful"}`))
}

func DeleteUserBySession(session_id string) bool {
	query := `UPDATE sessions SET session_token = NULL WHERE session_token = ?`
	_, err := databases.DB.Exec(query, session_id)
	if err != nil {
		log.Printf("Error updating session: %v", err)
		return false
	}
	return true
}

func SetSessionToken(id int, token string) {
	update_query := `UPDATE sessions SET expires_at = DATETIME('now', '+1 hour'), session_token = ? WHERE user_id = ?`
	insert_query := `INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, DATETIME('now', '+1 hour'))`

	// Try to update session
	res, err := databases.DB.Exec(update_query, token, id)
	if err != nil {
		log.Println(err)
	}

	// If no rows affected, insert new session
	if count, _ := res.RowsAffected(); count == 0 {
		_, err = databases.DB.Exec(insert_query, id, token)
		if err != nil {
			log.Println(err)
		}
	}
}

func GetUserHash(username string) (int, string) {
	var hash string
	var id int
	query := `SELECT id, password FROM users WHERE nickname = ? OR email = ?`
	err := databases.DB.QueryRow(query, username, username).Scan(&id, &hash)
	if err != nil {
		log.Println(err)
	}
	return id, hash
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
