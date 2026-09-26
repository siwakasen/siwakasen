package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/siwakasen/siwakasen/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	fmt.Printf("Listen to port %v", port)

	mux.HandleFunc("/addmoji", handlers.AddMoji)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
