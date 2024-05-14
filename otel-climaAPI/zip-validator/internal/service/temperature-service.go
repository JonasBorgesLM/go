package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// DataResponse representa a estrutura dos dados de resposta de temperatura
type DataResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// GetTemperatureServiceResponse representa a estrutura da resposta do serviço de temperatura
type GetTemperatureServiceResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    DataResponse `json:"data,omitempty"`
}

// GetTemperatureServiceInterface define a interface para o serviço de obtenção de temperatura
type GetTemperatureServiceInterface interface {
	GetTemperatureService(ctx context.Context, cep string) (GetTemperatureServiceResponse, error)
}

// GetTemperatureService implementa o serviço de obtenção de temperatura
type GetTemperatureService struct {
	client *http.Client
}

// NewGetTemperatureService cria uma nova instância de GetTemperatureService
func NewGetTemperatureService() *GetTemperatureService {
	return &GetTemperatureService{
		client: &http.Client{},
	}
}

// GetTemperatureService obtém as informações de temperatura com base no CEP fornecido
func (s *GetTemperatureService) GetTemperatureService(ctx context.Context, cep string) (GetTemperatureServiceResponse, error) {
	// Constrói a URL do serviço de clima
	weatherServiceURL := viper.GetString("WEATHER_SERVICE")
	url := weatherServiceURL + "?cep=" + cep

	// Cria uma nova requisição HTTP com o contexto fornecido
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	// Injeta o contexto de propagação do OpenTelemetry
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Executa a requisição HTTP
	res, err := s.client.Do(req)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}
	defer res.Body.Close()

	// Decodifica a resposta JSON
	var response GetTemperatureServiceResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return GetTemperatureServiceResponse{}, err
	}

	// Verifica se a resposta indica sucesso
	if !response.Success {
		return GetTemperatureServiceResponse{}, errors.New(response.Message)
	}

	return response, nil
}
