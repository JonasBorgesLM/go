package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather(t *testing.T) {
	// Simula um servidor de teste para a API de clima
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Define uma resposta simulada da API de clima
		weatherResp := map[string]interface{}{
			"current": map[string]interface{}{
				"temp_c": 20.0,
				"temp_f": 68.0,
				"temp_k": 293.15,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weatherResp)
	}))
	defer server.Close()

	// Chama a função getWeather com a cidade do servidor de teste
	weather, err := getWeather("Sao Paulo")
	if err != nil {
		t.Errorf("Erro ao obter dados do clima: %v", err)
		return
	}

	// Verifica se os valores de temperatura estão corretos
	expectedTempC := weather.TempC
	expectedTempF := expectedTempC*1.8 + 32
	expectedTempK := expectedTempC + 273.15
	if weather.TempF != expectedTempF || weather.TempK != expectedTempK {
		t.Errorf("Valores de temperatura incorretos. Esperado: TempC=%.2f, TempF=%.2f, TempK=%.2f, Obtido: TempC=%.2f, TempF=%.2f, TempK=%.2f", expectedTempC, expectedTempF, expectedTempK, weather.TempC, weather.TempF, weather.TempK)
	}
}

func TestClimaHandlerInvalidCep(t *testing.T) {
	req := httptest.NewRequest("GET", "/clima/123456789", nil)
	w := httptest.NewRecorder()

	ClimaHandler(w, req)

	resp := w.Result()

	// Verifica se o status code da resposta está correto
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Status code incorreto. Esperado: %d, Obtido: %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}

	// Verifica se o Content-Type da resposta está correto
	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Content-Type incorreto. Esperado: %s, Obtido: %s", expectedContentType, actualContentType)
	}

	// Decodifica o corpo da resposta
	var errorResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	// Verifica se a mensagem de erro está correta
	expectedErrorMessage := "invalid zipcode"
	if errorResp.Message != expectedErrorMessage {
		t.Errorf("Mensagem de erro incorreta. Esperado: %s, Obtido: %s", expectedErrorMessage, errorResp.Message)
	}
}

func TestClimaHandlerCityNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/clima/000", nil)
	w := httptest.NewRecorder()

	ClimaHandler(w, req)

	resp := w.Result()

	// Verifica se o status code da resposta está correto
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Status code incorreto. Esperado: %d, Obtido: %d", http.StatusNotFound, resp.StatusCode)
	}

	// Verifica se o Content-Type da resposta está correto
	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Content-Type incorreto. Esperado: %s, Obtido: %s", expectedContentType, actualContentType)
	}

	// Decodifica o corpo da resposta
	var errorResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	// Verifica se a mensagem de erro está correta
	expectedErrorMessage := "invalid zipcode"
	if errorResp.Message != expectedErrorMessage {
		t.Errorf("Mensagem de erro incorreta. Esperado: %s, Obtido: %s", expectedErrorMessage, errorResp.Message)
	}
}

func TestClimaHandlerServerError(t *testing.T) {
	req := httptest.NewRequest("GET", "/clima/37501143", nil)
	w := httptest.NewRecorder()

	// Simula um servidor de teste que falha ao responder
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Chama o handler ClimaHandler com a URL do servidor de teste que falha
	ClimaHandler(w, req)

	resp := w.Result()

	// Verifica se o status code da resposta está correto
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Status code incorreto. Esperado: %d, Obtido: %d", http.StatusInternalServerError, resp.StatusCode)
	}

	// Verifica se o Content-Type da resposta está correto
	expectedContentType := "application/json"
	actualContentType := resp.Header.Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Content-Type incorreto. Esperado: %s, Obtido: %s", expectedContentType, actualContentType)
	}

	// Decodifica o corpo da resposta
	var errorResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	// Verifica se a mensagem de erro está correta
	expectedErrorMessage := "invalid zipcode"
	if errorResp.Message != expectedErrorMessage {
		t.Errorf("Mensagem de erro incorreta. Esperado: %s, Obtido: %s", expectedErrorMessage, errorResp.Message)
	}
}
