package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"RescisaoDesk/internal/pdf"
	"RescisaoDesk/internal/rescisao"

	"github.com/pkg/browser"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ValidarDados validates form data in the domain layer.
func (a *App) ValidarDados(entrada rescisao.EntradaRescisao) (rescisao.DadosRescisao, error) {
	return rescisao.ValidarEntrada(entrada)
}

// CalcularResultado calcula as verbas da rescisão a partir dos dados validados do formulário.
func (a *App) CalcularResultado(entrada rescisao.EntradaRescisao) (rescisao.ResultadoRescisao, error) {
	dados, err := rescisao.ValidarEntrada(entrada)
	if err != nil {
		return rescisao.ResultadoRescisao{}, err
	}

	return rescisao.CalcularValoresRescisao(
		dados.SalarioBase,
		dados.AvosDecimoDisponiveis,
		dados.AvosDecimo,
		dados.SalarioBase,
		dados.AvosFeriasDisponiveis,
		dados.AvosFeriasPagos,
		dados.SalarioBase,
		dados.DiasTrabalhados,
	)
}

// GerarPDF gera o termo de rescisão em PDF usando o resultado já calculado pelo domínio.
func (a *App) GerarPDF(entrada rescisao.EntradaRescisao) (string, error) {
	dados, err := rescisao.ValidarEntrada(entrada)
	if err != nil {
		return "", err
	}

	resultado, err := a.CalcularResultado(entrada)
	if err != nil {
		return "", err
	}

	dir := filepath.Join(os.TempDir(), "RescisaoDesk")
	filename := pdf.NomeArquivo(dados.Servidor.Nome, dados.Servidor.Matricula) + ".pdf"
	destino := filepath.Join(dir, filename)

	path, err := pdf.GerarArquivo(pdf.Documento{
		Nome:                 dados.Servidor.Nome,
		Matricula:            dados.Servidor.Matricula,
		Cargo:                dados.Servidor.Cargo,
		Banco:                dados.Banco,
		Agencia:              dados.Agencia,
		Conta:                dados.Conta,
		Pix:                  dados.Pix,
		DataAdmissao:         dados.DataAdmissao,
		DataDesligamento:     dados.DataDesligamento,
		SalarioBase:          fmtMoney(dados.SalarioBase),
		DiasTrabalhados:      dados.DiasTrabalhados,
		ValorDiasTrabalhados: int64(resultado.DiasTrabalhados),
		DecimoTerceiro:       int64(resultado.DecimoTerceiro),
		Ferias:               int64(resultado.Ferias),
		TercoFerias:          int64(resultado.TercoFerias),
		Total:                int64(resultado.Total),
		AvosFeriasDevidos:    resultado.AvosFeriasDevidos,
	}, destino)
	if err != nil {
		return "", err
	}
	if err := browser.OpenFile(path); err != nil {
		return path, fmt.Errorf("PDF gerado, mas não foi possível abri-lo: %w", err)
	}
	return path, nil
}

func fmtMoney(value rescisao.Money) string {
	if value < 0 {
		return fmt.Sprintf("-R$ %.2f", float64(value)/100.0)
	}
	return fmt.Sprintf("R$ %.2f", float64(value)/100.0)
}
