package stress

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Stress representa um teste de stress para uma URL específica.
type Stress struct {
	URL         string        // A URL a ser testada
	Method      string        // O método HTTP a ser usado
	Concurrency int           // O número de requisições simultâneas a serem feitas
	Requests    int           // O número total de requisições a serem feitas
	Timeout     time.Duration // O tempo limite para cada requisição
	Verbose     bool          // Saída detalhada
	VerifyTls   bool          // Verificar TLS
	Report      *StressReport // O relatório de stress
	Mu          sync.Mutex    // Mutex para garantir exclusão mútua
}

// NewStress cria uma nova instância de Stress com os valores padrão.
func NewStress(url string, method string, concurrency int, requests int, timeout time.Duration, verbose bool, verifyTls bool) *Stress {
	return &Stress{
		URL:         url,
		Method:      method,
		Concurrency: concurrency,
		Requests:    requests,
		Timeout:     timeout,
		Verbose:     verbose,
		VerifyTls:   verifyTls,
		Report:      NewStressReport(),
	}
}

func (s *Stress) Run() error {
	s.run()

	return nil
}

// Run executa o teste de estresse
func (s *Stress) run() {
	start := time.Now()

	var wg sync.WaitGroup

	// Inicia as goroutines concorrentes para lidar com as requisições
	s.startWorkers(&wg)

	// Aguarda o término de todas as goroutines
	wg.Wait()

	// Calcula e define as estatísticas do teste de estresse
	s.calculateStatistics(start)
}

// startWorkers inicia as goroutines para lidar com as requisições
func (s *Stress) startWorkers(wg *sync.WaitGroup) {
	// Inicia o número especificado de goroutines concorrentes
	for i := 0; i < s.Concurrency; i++ {
		wg.Add(1)
		go s.worker(i+1, wg)
	}

	// Se o número de requisições não for divisível pelo número de requisições concorrentes
	// Executa as requisições restantes
	remainder := s.Requests % s.Concurrency
	for i := 0; i < remainder; i++ {
		wg.Add(1)
		go s.worker(i+1, wg)
	}
}

// worker é uma goroutine que executa uma parte das requisições
func (s *Stress) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Executa o número especificado de requisições
	for j := 0; j < s.Requests/s.Concurrency; j++ {
		s.runRequest(id)
	}
}

// calculateStatistics calcula e define as estatísticas do teste de estresse
func (s *Stress) calculateStatistics(start time.Time) {
	// Calcula o tempo decorrido em milissegundos
	elapsed := time.Since(start).Milliseconds()

	// Define o tempo total decorrido
	s.Report.TotalTime = float64(elapsed)

	// Calcula o AverageTime
	s.Report.AverageTime = s.Report.TotalTime / float64(s.Report.Requests)

	// Calcula a Porcentagem de Sucesso
	s.Report.PercentageSucceeded = float64(s.Report.Succeeded) / float64(s.Report.Requests) * 100

	// Calcula a Porcentagem de Falha
	s.Report.PercentageFailed = float64(s.Report.Failed) / float64(s.Report.Requests) * 100

	// Calcula a Porcentagem de Tempo Limite Excedido
	s.Report.PercentageTimedOut = float64(s.Report.TimedOut) / float64(s.Report.Requests) * 100
}

// runRequest executa uma requisição HTTP
func (s *Stress) runRequest(concurrencyGroup int) {
	start := time.Now()

	// Cria a requisição HTTP
	req, err := s.createRequest()
	if err != nil {
		panic(err)
	}

	// Executa a requisição HTTP
	res, err := s.executeRequest(req)

	// Calcula o tempo decorrido em milissegundos
	elapsed := time.Since(start).Milliseconds()

	// Se Verbose estiver habilitado, imprime a requisição
	if s.Verbose {
		s.printRequestInfo(concurrencyGroup, elapsed, res.StatusCode)
	}

	// Atualiza o relatório de estresse com os resultados da requisição
	s.updateReport(res, err, elapsed)
}

// createRequest cria uma nova requisição HTTP
func (s *Stress) createRequest() (*http.Request, error) {
	req, err := http.NewRequest(s.Method, s.URL, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// executeRequest executa a requisição HTTP e retorna a resposta
func (s *Stress) executeRequest(req *http.Request) (*http.Response, error) {
	// Cria um novo transporte HTTP
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !s.VerifyTls},
	}

	// Cria um novo cliente HTTP
	client := &http.Client{
		Timeout:   time.Duration(s.Timeout) * time.Second,
		Transport: tr,
	}

	// Executa a requisição HTTP
	res, err := client.Do(req)

	return res, err
}

// printRequestInfo imprime informações sobre a requisição
func (s *Stress) printRequestInfo(concurrencyGroup int, elapsed int64, statusCode int) {
	// Constrói a string formatada com as informações da requisição
	requestInfo := fmt.Sprintf("%d | %d %s %s Tempo: %d ms, Status: %d", concurrencyGroup, s.Report.Requests+1, s.Method, s.URL, elapsed, statusCode)
	// Imprime as informações da requisição
	fmt.Println(requestInfo)
}

// updateReport atualiza o relatório de estresse com os resultados da requisição
func (s *Stress) updateReport(res *http.Response, err error, elapsed int64) {
	// Trava o mutex para evitar acesso concorrente ao relatório
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Atualiza o relatório com base nos resultados da requisição
	if err != nil {
		s.handleRequestError(err)
	} else {
		s.handleRequestSuccess(res, elapsed)
	}

	// Atualiza o número total de requisições
	s.Report.Requests++
}

// handleRequestError atualiza o relatório quando ocorre um erro na requisição
func (s *Stress) handleRequestError(err error) {
	// Incrementa o contador de falhas
	s.Report.Failed++

	// Verifica se o erro é um timeout
	if err.Error() == http.ErrHandlerTimeout.Error() {
		// Incrementa o contador de timeouts
		s.Report.TimedOut++
	}

	// Imprime o erro
	fmt.Println(err)
}

// handleRequestSuccess atualiza o relatório quando a requisição é bem-sucedida
func (s *Stress) handleRequestSuccess(res *http.Response, elapsed int64) {
	// Incrementa o contador de requisições bem-sucedidas ou falhas, dependendo do status da resposta
	if res.StatusCode != http.StatusOK {
		s.Report.Failed++
	} else {
		s.Report.Succeeded++
	}

	// Incrementa o contador para o código de status da resposta
	s.Report.StatusRequests[res.StatusCode]++

	// Atualiza o tempo mais rápido e mais lento, se necessário
	s.updateFastestAndSlowestTime(elapsed)
}

// updateFastestAndSlowestTime atualiza os tempos mais rápido e mais lento das requisições
func (s *Stress) updateFastestAndSlowestTime(elapsed int64) {
	// Se o tempo decorrido for mais rápido do que o tempo mais rápido registrado ou se este for o primeiro tempo registrado, atualiza o tempo mais rápido
	if elapsed < s.Report.FastestTime || s.Report.FastestTime == 0 {
		s.Report.FastestTime = elapsed
	}

	// Se o tempo decorrido for mais lento do que o tempo mais lento registrado, atualiza o tempo mais lento
	if elapsed > s.Report.SlowestTime {
		s.Report.SlowestTime = elapsed
	}
}
