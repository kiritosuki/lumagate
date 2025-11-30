package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "v1 C") })
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	log.Fatal(http.ListenAndServe(":9003", nil))
}
