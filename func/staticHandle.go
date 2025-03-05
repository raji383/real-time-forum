package forum

import (
	"net/http"
	"os"
	"strings"
)

func StaticHandle(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	_, err := os.Stat("." + r.URL.Path)
	if strings.HasSuffix(r.URL.Path, "/") || err != nil {
		w.WriteHeader(http.StatusForbidden)
		http.Error(w, "static", http.StatusForbidden)
		// http.ServeFile(w, r, "Error/403.html")
		return
	}
	fs.ServeHTTP(w, r)
}
