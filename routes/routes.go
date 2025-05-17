package routes

import (
    "net/http"
    "real-time-forum/handlers"
    "github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

    // Static files
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Pages
    r.HandleFunc("/", handlers.HomeHandler)
  
    // Auth
    r.HandleFunc("/login", handlers.Login)
    r.HandleFunc("/Signup", handlers.Signup)
    r.HandleFunc("/logout", handlers.Logout)
    r.HandleFunc("/check-session", handlers.CheckSession)
    r.HandleFunc("/posts", handlers.PostsHandler).Methods("POST")


    return r
}
