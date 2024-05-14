package utils

import (
	"net/http"

	"github.com/goccy/go-json"
)

// Response representa a estrutura da resposta JSON
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseDTO é uma estrutura para transferir dados de resposta
type ResponseDTO struct {
	StatusCode int
	Success    bool
	Message    string
	Data       interface{}
}

// JsonResponse envia uma resposta JSON HTTP
func JsonResponse(w http.ResponseWriter, response ResponseDTO) {
	res := createResponse(response)

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, response.StatusCode, jsonResponse)
}

// createResponse cria uma estrutura de resposta JSON a partir de ResponseDTO
func createResponse(response ResponseDTO) Response {
	return Response{
		Message: response.Message,
		Data:    response.Data,
		Success: response.Success,
	}
}

// writeJSONResponse escreve a resposta JSON no ResponseWriter
func writeJSONResponse(w http.ResponseWriter, statusCode int, jsonResponse []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
