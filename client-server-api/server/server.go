package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type USDBRL struct {
	Code       string `json:"-"`
	Codein     string `json:"-"`
	Name       string `json:"-"`
	High       string `json:"-"`
	Low        string `json:"-"`
	VarBid     string `json:"-"`
	PctChange  string `json:"-"`
	Bid        string `json:"bid"`
	Ask        string `json:"-"`
	Timestamp  string `json:"-"`
	CreateDate string `json:"-"`
}

func main() {
	http.HandleFunc("/cotacao", cotacao)
	http.ListenAndServe(":8080", nil)
}

func cotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 2010*time.Millisecond)
	defer cancel()

	dolarCotation(ctx)
}

func dolarCotation(ctx context.Context) {
	req, err := getBidUSDBRL()
	if err != nil {
		log.Println("Get quote from api failed")
	}
	fmt.Println(req)

	select {
	case <-ctx.Done():
		log.Println("Request for dollar quotation canceled")
	case <-time.After(200 * time.Millisecond):
		log.Println("Segunda ação")
	}
}

func getBidUSDBRL() (*USDBRL, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		log.Println("Request to API failed")
		return nil, err
	}
	defer req.Body.Close()

	var bid USDBRL

	err = json.NewDecoder(req.Body).Decode(&bid)
	if err != nil {
		log.Println("Failed to decode")
		return nil, err
	}
	fmt.Println(bid)

	return &bid, nil
}
