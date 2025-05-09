package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"real-time-forum/databases"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

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

	if password != password2 {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// hashedPassword
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to encrypt password", http.StatusInternalServerError)
		return
	}

	// INSERT
	_, err = databases.DB.Exec(`
        INSERT INTO users (nickname, age, gender, first_name, last_name, email, password)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		username, age, gender, firstName, lastName, email, hashedPassword)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	err1 := r.ParseForm()
	if err1 != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	nickname := r.FormValue("user")
	password := r.FormValue("password")
	fmt.Println("nickname", nickname, "password", password)

	query := `SELECT password FROM users WHERE nickname = ?`

	var hashedPassword string
	err := databases.DB.QueryRow(query, nickname).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			
			http.Error(w, "User not found1", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}
	user_id,_ := GetUserHash(nickname)
	session := GenerateToken(32) // TODO: UUID bonus csrf implementation genrate csrf read it in front end js and match it with server go

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session,
		Expires:  time.Now().Add(time.Hour * 1),
		HttpOnly: true,
	})

	 SetSessionToken(user_id, session)
	 

	 w.Header().Set("Content-Type", "application/json")
	 w.WriteHeader(http.StatusOK)
	 w.Write([]byte(`{"success": true, "message": "Login successful"}`))
	 
}

func SetSessionToken(id int, token string) {
	update_query := `UPDATE sessions SET expires_at = DATETIME('now', '+1 hour') , session_token = ? WHERE user_id = ?`
	insert_query := `INSERT INTO sessions (user_id,session_token,expires_at) VALUES (? , ? , DATETIME('now','+1 hour'))`
	res, err := databases.DB.Exec(update_query, token, id)
	if err != nil {
		log.Println(err)
	}
	if count, _ := res.RowsAffected(); count > 0 {
		return
	}
	_, err = databases.DB.Exec(insert_query, id, token)
	if err != nil {
		log.Println(err)
	}
}

func GetUserHash(username string) (int, string) {
	var hash string
	var id int
	query := `SELECT id,password_hash FROM users WHERE username = ? OR email = ?`
	err := databases.DB.QueryRow(query, username, username).Scan(&id, &hash)
	if err != nil {
		log.Println(err)
	}
	return id, hash
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("failed to generat token %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}



