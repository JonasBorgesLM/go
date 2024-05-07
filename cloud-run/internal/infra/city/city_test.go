package city

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocation(t *testing.T) {
	// Cria um servidor HTTP de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica se o método e o caminho da solicitação estão corretos
		if r.Method != http.MethodGet || r.URL.Path != "/ws/37501143/json/" {
			t.Errorf("Solicitação incorreta. Método: %s, Caminho: %s", r.Method, r.URL.Path)
			http.Error(w, "Solicitação incorreta", http.StatusBadRequest)
			return
		}

		// Define uma resposta simulada
		responseJSON := `{
			"cep": "37501143",
			"logradouro": "",
			"complemento": "",
			"bairro": "",
			"localidade": "Itajuba",
			"uf": "",
			"ibge": "",
			"gia": "",
			"ddd": "",
			"siafi": ""
		}`

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, responseJSON)
	}))
	defer server.Close()

	// Chama a função GetLocation com a URL do servidor de teste
	city, err := GetLocation("37501143")
	if err != nil {
		t.Errorf("Erro ao obter localização: %v", err)
		return
	}

	// Verifica se a cidade retornada está correta
	expectedCity := "Itajuba"
	if city != expectedCity {
		t.Errorf("Cidade incorreta. Esperado: %s, Obtido: %s", expectedCity, city)
	}
}

func TestRemoveAccents(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"São Paulo", "Sao Paulo"},
		{"Belo Horizonte", "Belo Horizonte"},
		{"Cidade com acentuação çãõ", "Cidade com acentuacao cao"},
		{"Árvore de maçã", "Arvore de maca"},
		{"Über", "Uber"},
	}

	for _, test := range tests {
		result := removeAccents(test.input)
		if result != test.expected {
			t.Errorf("Erro: para entrada '%s', esperado '%s', mas obtido '%s'", test.input, test.expected, result)
		}
	}
}

func TestGetLocationInvalidCep(t *testing.T) {
	// Chama a função GetLocation com um CEP inválido
	_, err := GetLocation("123456789")
	if err == nil {
		t.Error("Esperava-se um erro para um CEP inválido, mas nenhum erro foi retornado")
	}
}

func TestGetLocationServerFailure(t *testing.T) {
	// Simula um servidor de teste que falha ao responder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Chama a função GetLocation com a URL do servidor de teste que falha
	_, err := GetLocation(server.URL)
	if err == nil {
		t.Error("Esperava-se um erro devido a falha do servidor, mas nenhum erro foi retornado")
	}
}

func TestRemoveAccentsEmptyString(t *testing.T) {
	// Testa a função removeAccents com uma string vazia
	result := removeAccents("")
	if result != "" {
		t.Error("A função removeAccents não deve alterar uma string vazia")
	}
}

func TestRemoveAccentsNoAccents(t *testing.T) {
	// Testa a função removeAccents com uma string sem acentos
	input := "Hello World"
	result := removeAccents(input)
	if result != input {
		t.Error("A função removeAccents não deve alterar uma string sem acentos")
	}
}
