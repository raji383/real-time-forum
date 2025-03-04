package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func init() {
	
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
}

func IsAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session:", err)
		return false
	}

	auth, ok := session.Values["authenticated"].(bool)
	log.Println("Session Data:", session.Values) 

	return ok && auth
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
