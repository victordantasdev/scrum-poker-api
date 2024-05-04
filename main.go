package main

import (
	"log"
	"net/http"
	"os"

	"github.com/victordantasdev/scrum-poker-api/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	http.HandleFunc("GET /healtz", api.HealtzHandler)

	log.Println("Server running on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
