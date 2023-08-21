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
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	quotation, err := getQuotationUSDBRL(ctx)
	if err != nil {
		log.Println("Error: Request for dollar quotation canceled")
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	dollar := Dollar{Dollar: quotation.Bid}
	err = insertQuotation(ctx2, db, &dollar)
	if err != nil {
		log.Println("Error: Request for insert canceled")
	}

	err = json.NewEncoder(w).Encode(dollar)
	if err != nil {
		log.Println("Error: Failed to encode")
	}
}

func getQuotationUSDBRL(ctx context.Context) (*Quotation, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		log.Println("Error: Request to API failed")
		return nil, err
	}
	defer req.Body.Close()

	var quotations ListQuotation
	err = json.NewDecoder(req.Body).Decode(&quotations)
	if err != nil {
		log.Println("Error: Failed to decode")
		return nil, err
	}

	return &quotations.Quotation, nil
}

func insertQuotation(ctx context.Context, db *sql.DB, dollar *Dollar) error {
	stmt, err := db.Prepare("INSERT INTO quotations(dollar) VALUES($1)")
	if err != nil {
		println("Error: Failed to stmt ")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dollar.Dollar)
	if err != nil {
		println("Failed to execute")
		return err
	}
	return nil
}
