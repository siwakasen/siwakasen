package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/siwakasen/siwakasen/handlers"
)

func main() {
	port := 80
	http.HandleFunc("/addmoji", handlers.AddMoji)
	fmt.Printf("Listen to port %v", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
