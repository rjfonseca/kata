package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", helloHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("FIXME")) //nolint: errcheck
}
