package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Dollar struct {
	Dollar string `json:"dollar"`
}

func main() {
	quotation, err := getQuotation()
	if err != nil {
		log.Println("Error: Request to server failed")
	}

	var dollar Dollar
	err = json.Unmarshal(quotation, &dollar)
	if err != nil {
		log.Println("Error: Failed to Unmarshal")
	}

	writeToFile(dollar)
}

func getQuotation() ([]byte, error) {
	client := http.Client{Timeout: 300 * time.Millisecond}

	req, err := client.Get("http://localhost:8080/cotacao")
	if err != nil {
		log.Println("Error: Request to server canceled")
		return nil, err
	}
	defer req.Body.Close()

	resp, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Error: Failed to read")
		return nil, err
	}

	return resp, nil
}

func writeToFile(dollar Dollar) {
	arq, err := os.Create("./server/db/cotacao.txt")
	if err != nil {
		log.Println("Error: Text file")
	}
	defer arq.Close()

	_, err = arq.Write([]byte("Dólar: " + dollar.Dollar))
	if err != nil {
		log.Println("Error: Failed to write")
	}
}
