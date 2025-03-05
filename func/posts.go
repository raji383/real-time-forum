package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func Posts(w http.ResponseWriter, r *http.Request) {
	
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login-page", http.StatusSeeOther)
		return
	}

	ex := CheckCookie(cookie.Value)
	if !ex{
		http.Redirect(w, r, "/login-page", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	s := r.Form["topic"]
	hh := time.Now()
	t := hh.Format("2006-01-02 15:04:05")

	topic := ""

	for _, v := range s {
		topic += v + ", "
	}
	fmt.Println(cookie.Value)

	var userName string
	err = db.QueryRow("SELECT username FROM users WHERE session_token = ?", cookie.Value).Scan(&userName)
	if err == sql.ErrNoRows || err != nil {
		fmt.Println(err)
		return 
	}

	_, err = db.Exec("INSERT INTO posts (username, title, descriptions, time, topic) VALUES (?, ?, ?, ?, ?);", userName, title, description, t, topic)
	if err != nil {
		fmt.Println(err)
		http.ServeFile(w, r, "template/emailCheck.html")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
