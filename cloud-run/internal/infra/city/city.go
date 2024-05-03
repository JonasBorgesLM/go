package city

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Cep struct {
	Cep         string `json:"_"`
	Logradouro  string `json:"_"`
	Complemento string `json:"_"`
	Bairro      string `json:"_"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"_"`
	Ibge        string `json:"_"`
	Gia         string `json:"_"`
	Ddd         string `json:"_"`
	Siafi       string `json:"_"`
}

func GetLocation(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil || resp.Body == nil {
		return "", err
	}
	defer resp.Body.Close()

	var locationResp Cep
	if err = json.NewDecoder(resp.Body).Decode(&locationResp); err != nil {
		return "", err
	}

	city := removeAccents(locationResp.Localidade)

	return city, nil
}

func removeAccents(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)

	result, _, _ := transform.String(t, input)

	return result
}
