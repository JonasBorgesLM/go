package stress

// MapStatusRequests mapeia códigos de status HTTP para o número de requisições correspondentes.
type MapStatusRequests map[int]int

// StressReport é uma estrutura que representa um relatório de testes de estresse.
type StressReport struct {
	Requests            int               // O número total de requisições feitas
	Failed              int               // O número de requisições que falharam
	Succeeded           int               // O número de requisições que tiveram sucesso
	TimedOut            int               // O número de requisições que excederam o tempo limite
	TotalTime           float64           // O tempo total para completar todas as requisições (em milissegundos)
	AverageTime         float64           // O tempo médio para completar uma requisição (em milissegundos)
	FastestTime         int64             // O tempo mais rápido para completar uma requisição (em milissegundos)
	SlowestTime         int64             // O tempo mais lento para completar uma requisição (em milissegundos)
	PercentageSucceeded float64           // A porcentagem de requisições que tiveram sucesso
	PercentageFailed    float64           // A porcentagem de requisições que falharam
	PercentageTimedOut  float64           // A porcentagem de requisições que excederam o tempo limite
	StatusRequests      MapStatusRequests // O mapa de códigos de status de requisições para o número de requisições
}

// NewStressReport cria um novo relatório de stress com valores iniciais.
func NewStressReport() *StressReport {
	return &StressReport{
		Requests:            0,
		Failed:              0,
		Succeeded:           0,
		TimedOut:            0,
		TotalTime:           0,
		AverageTime:         0,
		FastestTime:         0,
		SlowestTime:         0,
		PercentageSucceeded: 0,
		PercentageFailed:    0,
		PercentageTimedOut:  0,
		StatusRequests:      make(MapStatusRequests),
	}
}
