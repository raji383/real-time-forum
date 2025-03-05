package forum

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ex := CheckCookie(cookie.Value)
	if !ex {
		http.Redirect(w, r, "/login-page", http.StatusSeeOther)
		return
	}

	_, err = db.Exec("UPDATE users SET session_token = NULL WHERE session_token = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Error logging out!", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
