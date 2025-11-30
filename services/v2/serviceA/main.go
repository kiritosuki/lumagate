package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "v2 A") })
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	log.Fatal(http.ListenAndServe(":9011", nil))
}
