package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type ListQuotation struct {
	Quotation Quotation `json:"USDBRL"`
}

type Quotation struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Dollar struct {
	Dollar string `json:"dollar"`
}

func main() {
	http.HandleFunc("/cotacao", cotacao)
	http.ListenAndServe(":8080", nil)
}

func cotacao(w http.ResponseWriter, r *http.Request) {
	const file string = "./server/db/server.db"

	db, err := sql.Open("sqlite", file)
	if err != nil {
		log.Println("Error: Establishing a Database Connection")
	}
	defer db.Close()

	ctx := r.Context()
	ctxRequest, cancelRequest := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancelRequest()

	quotation := getQuotationUSDBRL(ctxRequest)

	ctxDB, cancelDB := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancelDB()

	dollar := Dollar{Dollar: quotation.Bid}
	insertQuotation(ctxDB, db, &dollar)

	err = json.NewEncoder(w).Encode(dollar)
	if err != nil {
		log.Println("Error: Failed to encode")
	}
}

func getQuotationUSDBRL(ctx context.Context) *Quotation {
	select {
	case <-ctx.Done():
		log.Println("Error: Request for dollar quotation canceled")

		return nil

	case <-time.After(200 * time.Millisecond):
		req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			log.Println("Error: Request to API failed")

			return nil
		}
		defer req.Body.Close()

		var quotations ListQuotation
		err = json.NewDecoder(req.Body).Decode(&quotations)
		if err != nil {
			log.Println("Error: Failed to decode")

			return nil
		}

		return &quotations.Quotation
	}
}

func insertQuotation(ctx context.Context, db *sql.DB, dollar *Dollar) {
	select {
	case <-ctx.Done():
		log.Println("Error: Request for insert canceled")

	case <-time.After(10 * time.Millisecond):
		stmt, err := db.Prepare("INSERT INTO quotations(dollar) VALUES($1)")
		if err != nil {
			log.Println("Error: Failed to stmt")
		}
		defer stmt.Close()

		_, err = stmt.Exec(dollar.Dollar)
		if err != nil {
			log.Println("Error: Failed to execute")
		}
	}
}
