package handlers

import (
	"net/http"
	"regexp"

	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/service"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/internal/usecase"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/pkg/errors"
	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/pkg/utils"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Handler representa um manipulador HTTP.
type Handler struct {
	viaCepService     service.ViaCepServiceInterface
	weatherApiService service.WeatherApiServiceInterface
}

// NewHandler cria uma nova instância do manipulador.
func NewHandler(
	viaCepService service.ViaCepServiceInterface,
	weatherApiService service.WeatherApiServiceInterface,
) *Handler {
	return &Handler{
		viaCepService:     viaCepService,
		weatherApiService: weatherApiService,
	}
}

// GetTemperatures manipula a solicitação para obter temperaturas com base no CEP fornecido.
func (h *Handler) GetTemperatures(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	tracer := otel.Tracer(viper.GetString("SERVICE_NAME"))

	ctx, span := tracer.Start(ctx, "GetTemperaturesHandler")
	defer span.End()

	cepParam := r.URL.Query().Get("cep")
	cep, err := h.formatCEP(cepParam)
	if err != nil {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Success:    false,
		})
		return
	}

	getTemperaturesUseCase := usecase.NewGetTemperatureUseCase(h.viaCepService, h.weatherApiService)
	data, err := getTemperaturesUseCase.Execute(ctx, cep)
	if err != nil {
		statusCode := http.StatusNotFound

		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: statusCode,
			Message:    err.Error(),
			Success:    false,
		})
		return
	}

	utils.JsonResponse(w, utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Message:    http.StatusText(http.StatusOK),
		Success:    true,
		Data:       data,
	})
}

// formatCEP formata o CEP fornecido para o formato esperado (XXXXX-XXX).
func (h *Handler) formatCEP(cep string) (string, error) {
	cepRegEx := `^\d{5}-\d{3}$`

	if regexp.MustCompile(cepRegEx).MatchString(cep) {
		return cep, nil
	}

	if len(cep) != 8 {
		return "", errors.ErrInvalidCEP
	}

	return cep[:5] + "-" + cep[5:], nil
}
