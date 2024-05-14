package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/service"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/internal/usecase"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/pkg/errors"
	"github.com/JonasBorgesLM/go/otel-climaAPI/zip-validator/pkg/utils"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Handler representa o manipulador HTTP
type Handler struct {
	weatherApiService service.GetTemperatureServiceInterface
}

// NewHandler cria uma nova instância de Handler
func NewHandler(weatherApiService service.GetTemperatureServiceInterface) *Handler {
	return &Handler{
		weatherApiService: weatherApiService,
	}
}

// GetTemperatures é o manipulador para obter as temperaturas com base no CEP fornecido
func (h *Handler) GetTemperatures(w http.ResponseWriter, r *http.Request) {
	// Extrai o contexto do cabeçalho da requisição
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	tracer := otel.Tracer(viper.GetString("SERVICE_NAME"))

	// Inicia o span para monitoramento
	ctx, span := tracer.Start(ctx, "GetTemperaturesHandler")
	defer span.End()

	// Decodifica o corpo da requisição
	var input InputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Success:    false,
		})
		return
	}

	// Valida o formato do CEP
	if !h.validateCEP(input.Cep) {
		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    errors.ErrInvalidCEP.Error(),
			Success:    false,
		})
		return
	}

	// Executa o caso de uso para obter as temperaturas
	getTemperaturesUseCase := usecase.NewGetTemperatureUseCase(h.weatherApiService)
	data, err := getTemperaturesUseCase.Execute(ctx, input.Cep)
	if err != nil {
		if err.Error() == errors.ErrInvalidCEP.Error() {
			utils.JsonResponse(w, utils.ResponseDTO{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    err.Error(),
				Success:    false,
			})
			return
		}

		utils.JsonResponse(w, utils.ResponseDTO{
			StatusCode: http.StatusNotFound,
			Message:    errors.ErrCannotFindZipCode.Error(),
			Success:    false,
		})
		return
	}

	// Responde com os dados das temperaturas
	utils.JsonResponse(w, utils.ResponseDTO{
		StatusCode: http.StatusOK,
		Message:    http.StatusText(http.StatusOK),
		Success:    true,
		Data:       data,
	})
}

// InputDTO representa a estrutura de entrada da requisição
type InputDTO struct {
	Cep string `json:"cep"`
}

// validateCEP valida o formato do CEP
func (h *Handler) validateCEP(cep string) bool {
	regex := regexp.MustCompile(`^\d{8}$`)
	return len(cep) == 8 && regex.MatchString(cep)
}
