package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/siwakasen/siwakasen/handlers"
)

func main() {
	ghToken := os.Getenv("GH_TOKEN")
	if strings.TrimSpace(ghToken) == "" {
		log.Fatal(fmt.Errorf("GH_TOKEN is not set"))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	fmt.Printf("Listen to port %v", port)

	mux.HandleFunc("/addmoji", handlers.AddMoji)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
