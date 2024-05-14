package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/JonasBorgesLM/go/otel-climaAPI/location-weather/pkg/errors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

// ViaCEPResponse é a estrutura que representa a resposta da API ViaCEP.
type ViaCEPResponse struct {
	Erro        string `json:"erro"`
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

// ViaCepServiceInterface define a interface para o serviço ViaCep.
type ViaCepServiceInterface interface {
	GetCEPData(ctx context.Context, cep string) (*ViaCEPResponse, error)
}

// ViaCepService é a implementação do serviço ViaCep.
type ViaCepService struct {
	client *http.Client
}

// NewViaCepService cria uma nova instância de ViaCepService.
func NewViaCepService() *ViaCepService {
	return &ViaCepService{client: &http.Client{}}
}

// GetCEPData obtém os dados de um CEP usando a API ViaCEP.
func (s *ViaCepService) GetCEPData(ctx context.Context, cep string) (*ViaCEPResponse, error) {
	// Inicia um span para rastreamento de desempenho com OpenTelemetry.
	tracer := otel.Tracer(viper.GetString("SERVICE"))
	ctx, span := tracer.Start(ctx, "ViaCEPService.GetCEPData")
	defer span.End()

	// Monta a URL de solicitação para a API ViaCEP.
	url := "http://viacep.com.br/ws/" + cep + "/json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Faz a solicitação HTTP.
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Verifica se a resposta é bem-sucedida.
	if res.StatusCode == http.StatusOK {
		var viaCEPResponse ViaCEPResponse
		// Decodifica o corpo da resposta JSON.
		err = json.NewDecoder(res.Body).Decode(&viaCEPResponse)
		if err != nil {
			return nil, err
		}

		// Verifica se o CEP foi encontrado.
		if viaCEPResponse.Erro == "true" {
			return nil, errors.ErrCannotFindZipCode
		}

		return &viaCEPResponse, nil
	}

	// Retorna nil se o status da resposta não for 200 (OK).
	return nil, nil
}
