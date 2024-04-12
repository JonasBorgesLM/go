package main

import (
	"fmt"
	"net/http"

	"github.com/JonasBorgesLM/go/ratelimiter/configs"
)

func main() {
	configs, err := configs.LoadConfig("../../")
	if err != nil {
		panic(err)
	}

	fmt.Println(configs.DBDriver)

	http.HandleFunc("/", Root)
	http.ListenAndServe(":8080", nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("o/"))
}
