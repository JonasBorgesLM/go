package stress

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintStressReport imprime o relatório de teste de estresse
func (s *Stress) PrintStressReport() {
	// Define um tabwriter com configurações personalizadas
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Imprime cabeçalho geral
	fmt.Fprintln(w, "------------------------------------------------------")

	fmt.Fprintln(w, "               Relatório de Stress Test")
	fmt.Fprintln(w, "------------------------------------------------------")
	fmt.Fprintf(w, " Requisições:\t%d\n", s.Report.Requests)
	fmt.Fprintf(w, " Falhas:\t%d\n", s.Report.Failed)
	fmt.Fprintf(w, " Sucesso:\t%d\n", s.Report.Succeeded)
	fmt.Fprintf(w, " Tempo Limite Excedido:\t%d\n", s.Report.TimedOut)
	fmt.Fprintf(w, " Tempo Total:\t%.0f ms\n", s.Report.TotalTime)
	fmt.Fprintf(w, " Tempo Médio por Requisição:\t%.1f ms\n", s.Report.AverageTime)
	fmt.Fprintf(w, " Tempo mais Rápido:\t%d ms\n", s.Report.FastestTime)
	fmt.Fprintf(w, " Tempo mais Lento:\t%d ms\n", s.Report.SlowestTime)
	fmt.Fprintf(w, " Porcentagem de Sucesso:\t%.0f %%\n", s.Report.PercentageSucceeded)
	fmt.Fprintf(w, " Porcentagem de Falha:\t%.0f %%\n", s.Report.PercentageFailed)
	fmt.Fprintf(w, " Porcentagem de Tempo Limite Excedido:\t%.0f %%\n", s.Report.PercentageTimedOut)

	// Imprime separador entre seções
	fmt.Fprintln(w, "------------------------------------------------------")

	// Imprime cabeçalho para detalhes por código de status
	fmt.Fprintln(w, "             Requisições por Código de Status")
	fmt.Fprintln(w, "------------------------------------------------------")

	// Imprime detalhes por código de status
	for status, requests := range s.Report.StatusRequests {
		fmt.Fprintf(w, " Status %d:\t%d requisições\n", status, requests)
	}

	fmt.Fprintln(w, "------------------------------------------------------")

	// Escreve os dados no tabwriter
	w.Flush()
}
