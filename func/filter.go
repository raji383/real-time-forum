package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func Filter(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/filter" {
		http.Error(w, "page not found", 404)
		return
	}

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

	category := r.FormValue("category")


	rows, err := db.Query("SELECT id, username, title, descriptions, time, topic, likes, dislikes FROM posts WHERE topic LIKE ? ORDER BY time DESC", "%"+category+"%")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "h", 500)
		return
	}
	defer rows.Close()
	
	var arrFilter []Post
	var newFilter Post

	is, username := IsLoggedIn(r)
	newFilter.User.Username = username
	newFilter.User.IsLoggedIn = is

	arrFilter = append(arrFilter, newFilter)

	for rows.Next() {
		var username, title, descriptions, t, topic string
		var id, like, dislike int
		er := rows.Scan(&id, &username, &title, &descriptions, &t, &topic, &like, &dislike)

		if er != nil {
			http.Error(w, "htp", 500)
			return
		}

		R := GetUserReaction(r, id)

		if R == 1 {
			newFilter.Reaction.Like = true
		} else if R == -1 {
			newFilter.Reaction.Dislike = true
		} else {
			newFilter.Reaction.Like = false
			newFilter.Reaction.Dislike = false
		}

		newFilter.ID = id
		newFilter.Username = username
		newFilter.Title = title
		newFilter.Descriptions = descriptions
		newFilter.Time = t
		newFilter.Topic = topic
		newFilter.Likes = like
		newFilter.Dislikes = dislike

		arrFilter = append(arrFilter, newFilter)
		newFilter = Post{}
	}
	tmp, err:= template.ParseFiles("template/filter.html")
	if err != nil {
		http.Error(w, "ht", 500)
		return
	}
	err = tmp.Execute(w, arrFilter)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "htppp", 500)
			return
	}
}
