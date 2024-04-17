package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/JonasBorgesLM/go/stress/internal/stress"
	"github.com/spf13/cobra"
)

var url, method string
var concurrency, requests int
var timeout time.Duration
var verifyTls, verbose bool

var rootCmd = &cobra.Command{
	Use:   "stress",
	Short: "Um testador de estresse CLI para simular cargas em serviços web",
	Long: `O stress é uma ferramenta CLI desenvolvida em Go para realizar testes de estresse em serviços web.
Ele permite simular cargas em um serviço especificado, fornecendo estatísticas detalhadas sobre o desempenho do serviço sob carga.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Verifica se a URL é fornecida
		if url == "" {
			fmt.Println("Erro: a URL é obrigatória. Use a flag --url para fornecer a URL a ser testada.")
			os.Exit(1)
		}

		// Verifica se o número de solicitações e a concorrência são valores positivos
		if requests <= 0 || concurrency <= 0 {
			fmt.Println("Erro: tanto o número de solicitações quanto a concorrência devem ser valores positivos.")
			os.Exit(1)
		}

		// Cria e executa o teste de estresse
		s := stress.NewStress(url, method, concurrency, requests, timeout, verifyTls, verbose)
		err := s.Run()
		if err != nil {
			fmt.Println("Erro ao executar o teste de estresse:", err)
			os.Exit(1)
		}

		// Imprime o relatório do teste de estresse
		s.PrintStressReport()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Erro ao executar o comando raiz:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "A URL a ser testada (obrigatório)")
	rootCmd.Flags().StringVarP(&method, "method", "m", "GET", "O método HTTP a ser usado")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "O número de solicitações simultâneas a serem feitas")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1, "O número de solicitações a serem feitas")
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "O tempo limite em segundos")
	rootCmd.Flags().BoolVar(&verifyTls, "verify", false, "Verificar a autenticidade do certificado TLS")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Saída detalhada")
}
