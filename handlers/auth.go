package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", "./databases/forum.db")
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
	);
`
	db.Exec(creat)
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
	r.ParseForm()
	nickname := r.FormValue("user")
	password := r.FormValue("password")
	

	query := `SELECT password FROM users WHERE nickname = ?`

	var hashedPassword string
	err := db.QueryRow(query, nickname).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			http.Error(w, "User not found", http.StatusNotFound)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	_, err = db.Exec(`
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
