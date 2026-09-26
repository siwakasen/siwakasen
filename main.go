package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/siwakasen/siwakasen/handlers"
)

func main() {
	mux := http.NewServeMux()
	port := 8080
	fmt.Printf("Listen to port %v", port)

	mux.HandleFunc("/addmoji", handlers.AddMoji)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), mux))
}
