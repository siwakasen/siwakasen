package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/siwakasen/siwakasen/handlers"
)

func main() {
	port := 80
	http.HandleFunc("/addmoji", handlers.AddMoji)
	fmt.Printf("Listen to port %v", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
