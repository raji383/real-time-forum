package main

import (
	"fmt"
	"log"
	"net/http"


	"real-time-forum/routes"
)

func main() {


	router := routes.SetupRoutes()

	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
