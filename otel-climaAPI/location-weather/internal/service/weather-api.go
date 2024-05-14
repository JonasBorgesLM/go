package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

// WeatherAPIResponse representa a estrutura da resposta da API de clima.
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// WeatherApiServiceInterface define a interface para o serviço de API de clima.
type WeatherApiServiceInterface interface {
	GetWeatherData(ctx context.Context, location string) (*WeatherAPIResponse, error)
}

// WeatherApiService é a implementação do serviço de API de clima.
type WeatherApiService struct {
	client *http.Client
}

// NewWeatherApiService cria uma nova instância do serviço de API de clima.
func NewWeatherApiService() *WeatherApiService {
	return &WeatherApiService{
		client: &http.Client{},
	}
}

// GetWeatherData obtém os dados do clima para uma determinada localização.
func (s *WeatherApiService) GetWeatherData(ctx context.Context, location string) (*WeatherAPIResponse, error) {
	// Inicia um span para rastreamento de desempenho com OpenTelemetry.
	tracer := otel.Tracer(viper.GetString("SERVICE"))
	ctx, span := tracer.Start(ctx, "WeatherAPI.GetWeatherData")
	defer span.End()

	// Obtém a chave da API do ambiente.
	WEATHER_API_KEY := viper.GetString("WEATHER_API_KEY")

	// Constrói a URL da solicitação para a API de clima.
	urlString := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", WEATHER_API_KEY, url.QueryEscape(location))

	// Cria a solicitação HTTP.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlString, nil)
	if err != nil {
		return nil, err
	}

	// Envia a solicitação HTTP e obtém a resposta.
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Verifica o status da resposta HTTP.
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("cannot find weather data")
	}

	// Decodifica o corpo da resposta JSON.
	var weatherAPIResponse WeatherAPIResponse
	err = json.NewDecoder(res.Body).Decode(&weatherAPIResponse)
	if err != nil {
		return nil, err
	}

	return &weatherAPIResponse, nil
}
