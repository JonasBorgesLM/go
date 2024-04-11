package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", Root)
	http.ListenAndServe(":8080", nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("o/"))
}
