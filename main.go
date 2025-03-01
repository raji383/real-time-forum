package main

import (
	//"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/imag/", Test)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/regist.html", regist)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/posts", postsHandler)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
/*func createTables(db *sql.DB) {
}*/
func regist(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	templates, err := template.ParseFiles("./temp/regist.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template error:", err)
		return
	}

	err = templates.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	templates, err := template.ParseFiles("./temp/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template error:", err)
		return
	}

	err = templates.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/imag/" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	imagFile := "temp/imag/" + r.URL.Path[len("/imag/"):]
	http.ServeFile(w, r, imagFile)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Register Endpoint")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login Endpoint")
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Posts Endpoint")
}
