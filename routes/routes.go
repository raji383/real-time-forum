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
	r.HandleFunc("/api/posts", handlers.ApiPostsHandler).Methods("GET")
    r.HandleFunc("/like", handlers.LikeHandler).Methods("POST")
	// Add these to your SetupRoutes function
r.HandleFunc("/api/comments", handlers.GetCommentsHandler).Methods("GET")
r.HandleFunc("/comments", handlers.CreateCommentHandler).Methods("POST")
r.HandleFunc("/comment-reaction", handlers.CommentReactionHandler).Methods("POST")

	return r
}
