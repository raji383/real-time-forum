package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login-page", http.StatusSeeOther)
		return
	}

	ex := CheckCookie(cookie.Value)
	if !ex {
		http.Redirect(w, r, "/login-page", http.StatusSeeOther)
		return
	}

	fmt.Println(cookie.Value)

	posts, err := GetLikedPosts(cookie.Value, r)
	if err != nil {
		http.Error(w, "Error fetching liked posts", http.StatusInternalServerError)
		return
	}
	fmt.Println(len(posts))

	tmpl, err := template.ParseFiles("template/liked_posts.html")
	if err != nil {
		http.Error(w, "Error parsing !", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, posts)
}

func GetLikedPosts(cookieValue string, r *http.Request) ([]Post, error) {
	var userID int

	err := db.QueryRow("SELECT id FROM users WHERE session_token = ?", cookieValue).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rows, err := db.Query(`
		SELECT posts.id, posts.username, posts.title, posts.descriptions, posts.time, posts.topic, posts.likes, posts.dislikes 
		FROM posts 
		JOIN post_likes ON posts.id = post_likes.post_id 
		WHERE post_likes.user_id = ?
		ORDER BY posts.time DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Username, &p.Title, &p.Descriptions, &p.Time, &p.Topic, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}

		is, username := IsLoggedIn(r)
		p.User.Username = username
		p.User.IsLoggedIn = is

		R := GetUserReaction(r, p.ID)

		if R == 1 {
			p.Reaction.Like = true
		} else if R == -1 {
			p.Reaction.Dislike = true
		} else {
			p.Reaction.Like = false
			p.Reaction.Dislike = false
		}

		posts = append(posts, p)
	}

	return posts, nil
}
