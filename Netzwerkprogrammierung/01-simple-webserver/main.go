package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {
	http.HandleFunc("/", handleFunction)
	log.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", nil)
}
