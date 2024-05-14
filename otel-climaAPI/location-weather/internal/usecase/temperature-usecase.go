package usecase

import (
	"context"

	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/service"
)

// GetTemperaturesUseCase é a estrutura que encapsula os serviços necessários para obter as temperaturas.
type GetTemperaturesUseCase struct {
	viaCepService     service.ViaCepServiceInterface
	weatherApiService service.WeatherApiServiceInterface
}

// Response é a estrutura que representa a resposta com os dados da cidade e temperaturas.
type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// NewGetTemperatureUseCase cria uma nova instância de GetTemperaturesUseCase.
func NewGetTemperatureUseCase(
	viaCepService service.ViaCepServiceInterface,
	weatherApiService service.WeatherApiServiceInterface,
) *GetTemperaturesUseCase {
	return &GetTemperaturesUseCase{
		viaCepService:     viaCepService,
		weatherApiService: weatherApiService,
	}
}

// Execute é o método que executa o caso de uso para obter as temperaturas baseado no CEP.
func (u *GetTemperaturesUseCase) Execute(ctx context.Context, cep string) (Response, error) {
	// Obtém os dados do CEP usando o serviço viaCepService.
	cepData, err := u.viaCepService.GetCEPData(ctx, cep)
	if err != nil {
		return Response{}, err
	}

	// Obtém os dados do clima para a localidade usando o serviço weatherApiService.
	weatherData, err := u.weatherApiService.GetWeatherData(ctx, cepData.Localidade)
	if err != nil {
		return Response{}, err
	}

	// Calcula a temperatura em Fahrenheit e Kelvin.
	tempF := weatherData.Current.TempC*1.8 + 32
	tempK := weatherData.Current.TempC + 273

	// Retorna a resposta com os dados da cidade e temperaturas.
	return Response{
		City:  cepData.Localidade,
		TempC: weatherData.Current.TempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
