package stress

// StressInterface é uma interface que define métodos para executar testes de stress e imprimir relatórios.
type StressInterface interface {
	// Run executa o teste de stress.
	Run() error
	// PrintReport imprime o relatório do teste de stress.
	PrintReport()
}
