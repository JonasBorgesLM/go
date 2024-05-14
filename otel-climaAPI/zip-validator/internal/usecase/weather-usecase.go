package usecase

import (
	"context"

	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/service"
)

// GetTemperaturesUseCase é o caso de uso responsável por obter as temperaturas
type GetTemperaturesUseCase struct {
	weatherApiService service.GetTemperatureServiceInterface
}

// Response representa a estrutura da resposta contendo as temperaturas
type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// NewGetTemperatureUseCase cria uma nova instância de GetTemperaturesUseCase
func NewGetTemperatureUseCase(weatherApiService service.GetTemperatureServiceInterface) *GetTemperaturesUseCase {
	return &GetTemperaturesUseCase{
		weatherApiService: weatherApiService,
	}
}

// Execute executa o caso de uso para obter as temperaturas baseadas no CEP fornecido
func (u *GetTemperaturesUseCase) Execute(ctx context.Context, cep string) (Response, error) {
	// Chama o serviço para obter os dados de temperatura
	weatherData, err := u.weatherApiService.GetTemperatureService(ctx, cep)
	if err != nil {
		return Response{}, err
	}

	// Retorna os dados de temperatura formatados na estrutura Response
	return Response{
		City:  weatherData.Data.City,
		TempC: weatherData.Data.TempC,
		TempF: weatherData.Data.TempF,
		TempK: weatherData.Data.TempK,
	}, nil
}
