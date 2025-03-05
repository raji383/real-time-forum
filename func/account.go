package forum

import (
	"html/template"
	"net/http"
)

func SingupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/singup-page" {
		http.Error(w, "page not found", 404)
		return
	}

	tmp, err := template.ParseFiles("template/signup.html")
	if err != nil {
		http.Error(w, "htppp", 500)
		return
	}

	tmp.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/login-page" {
		http.Error(w, "page not found", 404)
		return
	}

	tmp, err := template.ParseFiles("template/login.html")
	if err != nil {
		http.Error(w, "htppp", 500)
		return
	}

	tmp.Execute(w, nil)
}
