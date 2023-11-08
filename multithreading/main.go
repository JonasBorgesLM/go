package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Cep struct {
	Code string `json:"cep,omitempty"`
}

func main() {
	var code string
	fmt.Scan(&code)

	if IsCEPValid(code) {
		cep := NewCep(code)

		ctx := context.Background()
		ctxRequest, cancelRequest := context.WithTimeout(ctx, 1*time.Second)
		defer cancelRequest()

		getCep(ctxRequest, *cep)

	} else {
		log.Println("CEP com formato inválido")
	}
}

// Expressão regular para verificar se a string corresponde ao formato de CEP no Brasil
// A expressão regular permite tanto "99999-999" como "99999999"
func IsCEPValid(cep string) bool {
	re := regexp.MustCompile(`^\d{5}-?\d{3}$`)

	return re.MatchString(cep)
}

// Gera um cep com formato valido "99999999"
func NewCep(code string) *Cep {
	var cep Cep
	re := regexp.MustCompile(`^\d{5}-\d{3}$`)

	if re.MatchString(code) {
		cep.Code = code

		return &cep
	} else {
		cep.Code = fmt.Sprintf("%s-%s", code[:5], code[5:])

		return &cep
	}
}

// Faz a requisição a api externa e armazena em um canal
func getInfoCep(api, url string, ch chan string) {
	req, err := http.Get(url)
	if err != nil {
		log.Printf("%s: Error - %v", url, err)
		return
	}
	defer req.Body.Close()

	resp, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Error: Failed to decode")
		return
	}

	ch <- api + "\n" + string(resp)
}

// Imprime as informações do cep com timeout de 1 second
func getCep(ctx context.Context, cep Cep) {
	urlCDN := "https://cdn.apicep.com/file/apicep/" + cep.Code + ".json"
	urlVIACEP := "http://viacep.com.br/ws/" + cep.Code + "/json/"

	ch := make(chan string)

	go getInfoCep("CDN", urlCDN, ch)
	go getInfoCep("VIACEP", urlVIACEP, ch)

	select {
	case <-ctx.Done():
		log.Println("Error: Request for cep canceled")

	case <-time.After(1 * time.Second):
		response := <-ch
		fmt.Println(response)
	}
}
