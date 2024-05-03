package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/JonasBorgesLM/go/clound-run/internal/infra/city"
	"github.com/go-chi/chi/v5"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func getWeather(city string) (WeatherResponse, error) {
	apiKey := "6b44a0c653be4aaa9c7130327240205"
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, strings.ReplaceAll(city, " ", "+"))
	resp, err := http.Get(url)
	if err != nil {
		return WeatherResponse{}, err
	}
	defer resp.Body.Close()

	var weatherResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return WeatherResponse{}, err
	}

	tempC, _ := weatherResp["current"].(map[string]interface{})["temp_c"].(float64)
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15

	return WeatherResponse{TempC: tempC, TempF: tempF, TempK: tempK}, nil
}

func ClimaHandler(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if len(cep) != 8 {
		errorResponse := ErrorResponse{Message: "invalid zipcode"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	city, err := city.GetLocation(cep)
	if err != nil || len(city) == 0 {
		errorResponse := ErrorResponse{Message: "can not find zipcode"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	weather, err := getWeather(city)
	if err != nil {
		errorResponse := ErrorResponse{Message: "Error fetching weather data"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}
