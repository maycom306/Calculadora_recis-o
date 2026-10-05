package rescisao

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Money stores a monetary amount as integer cents.
type Money int64

// EntradaRescisao contains raw values received from the form.
type EntradaRescisao struct {
	Nome             string `json:"nome"`
	Matricula        string `json:"matricula"`
	Cargo            string `json:"cargo"`
	Banco            string `json:"banco"`
	Agencia          string `json:"agencia"`
	Conta            string `json:"conta"`
	Pix              string `json:"pix"`
	DataAdmissao     string `json:"dataAdmissao"`
	DataDesligamento string `json:"dataDesligamento"`
	SalarioBase      string `json:"salarioBase"`
	AvosDecimo       string `json:"avosDecimo"`
	AvosFeriasPagos  string `json:"avosFeriasPagos"`
}

// Servidor contains identification fields used by the term.
type Servidor struct {
	Nome      string `json:"nome"`
	Matricula string `json:"matricula"`
	Cargo     string `json:"cargo"`
}

// DadosRescisao contains validated values used by domain calculations.
type DadosRescisao struct {
	Servidor              Servidor `json:"servidor"`
	Banco                 string   `json:"banco"`
	Agencia               string   `json:"agencia"`
	Conta                 string   `json:"conta"`
	Pix                   string   `json:"pix"`
	DataAdmissao          string   `json:"dataAdmissao"`
	DataDesligamento      string   `json:"dataDesligamento"`
	DiasTrabalhados       int      `json:"diasTrabalhados"`
	SalarioBase           Money    `json:"salarioBaseCentavos"`
	AvosDecimo            int      `json:"avosDecimo"`
	AvosDecimoDisponiveis int      `json:"avosDecimoDisponiveis"`
	AvosDecimoDevidos     int      `json:"avosDecimoDevidos"`
	AvosFeriasPagos       int      `json:"avosFeriasPagos"`
	AvosFeriasDisponiveis int      `json:"avosFeriasDisponiveis"`
	AvosFeriasDevidos     int      `json:"avosFeriasDevidos"`
}

// ValidarEntrada validates and normalizes values received from the form.
func ValidarEntrada(entrada EntradaRescisao) (DadosRescisao, error) {
	nome := strings.TrimSpace(entrada.Nome)
	if nome == "" {
		return DadosRescisao{}, fmt.Errorf("informe o nome do servidor")
	}

	dataAdmissao, err := parseDate(entrada.DataAdmissao, "admissão")
	if err != nil {
		return DadosRescisao{}, err
	}

	dataDesligamento, err := parseDate(entrada.DataDesligamento, "desligamento")
	if err != nil {
		return DadosRescisao{}, err
	}
	if dataDesligamento.Before(dataAdmissao) {
		return DadosRescisao{}, fmt.Errorf("a data de desligamento não pode ser anterior à admissão")
	}

	salarioBase, err := ParseMoney(entrada.SalarioBase)
	if err != nil {
		return DadosRescisao{}, fmt.Errorf("salário base: %w", err)
	}

	avosDecimo := 0
	if value := strings.TrimSpace(entrada.AvosDecimo); value != "" {
		avosDecimo, err = strconv.Atoi(value)
		if err != nil || avosDecimo < 0 {
			return DadosRescisao{}, fmt.Errorf("avos do 13º já recebidos devem ser um inteiro não negativo")
		}
	}
	avosDecimoDisponiveis, err := CalcularAvosDecimo(entrada.DataAdmissao, entrada.DataDesligamento)
	if err != nil {
		return DadosRescisao{}, err
	}
	avosDecimoDevidos := avosDecimoDisponiveis - avosDecimo
	if avosDecimoDevidos < 0 {
		avosDecimoDevidos = 0
	}

	avosFeriasPagos := 0
	if value := strings.TrimSpace(entrada.AvosFeriasPagos); value != "" {
		avosFeriasPagos, err = strconv.Atoi(value)
		if err != nil || avosFeriasPagos < 0 {
			return DadosRescisao{}, fmt.Errorf("avos de férias já gozados/pagos deve ser um inteiro não negativo")
		}
	}

	avosFeriasDisponiveis, err := CalcularAvosFerias(entrada.DataAdmissao, entrada.DataDesligamento)
	if err != nil {
		return DadosRescisao{}, err
	}
	avosFeriasDevidos, err := CalcularSaldoAvosFerias(avosFeriasDisponiveis, avosFeriasPagos)
	if err != nil {
		return DadosRescisao{}, err
	}

	return DadosRescisao{
		Servidor: Servidor{
			Nome:      nome,
			Matricula: strings.TrimSpace(entrada.Matricula),
			Cargo:     strings.TrimSpace(entrada.Cargo),
		},
		Banco:                 strings.TrimSpace(entrada.Banco),
		Agencia:               strings.TrimSpace(entrada.Agencia),
		Conta:                 strings.TrimSpace(entrada.Conta),
		Pix:                   strings.TrimSpace(entrada.Pix),
		DataAdmissao:          dataAdmissao.Format(time.DateOnly),
		DataDesligamento:      dataDesligamento.Format(time.DateOnly),
		DiasTrabalhados:       contarDiasMesDesligamento(dataDesligamento),
		SalarioBase:           salarioBase,
		AvosDecimo:            avosDecimo,
		AvosDecimoDisponiveis: avosDecimoDisponiveis,
		AvosDecimoDevidos:     avosDecimoDevidos,
		AvosFeriasPagos:       avosFeriasPagos,
		AvosFeriasDisponiveis: avosFeriasDisponiveis,
		AvosFeriasDevidos:     avosFeriasDevidos,
	}, nil
}

func contarDiasMesDesligamento(data time.Time) int {
	dias := data.Day()
	ultimoDia := time.Date(data.Year(), data.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if dias == 31 || dias == ultimoDia {
		return 30
	}
	return dias
}

// ParseMoney parses a non-negative Brazilian currency value into cents.
func ParseMoney(value string) (Money, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimSpace(strings.TrimPrefix(value, "R$"))
	value = strings.ReplaceAll(value, " ", "")
	if value == "" || strings.HasPrefix(value, "-") || strings.HasPrefix(value, "+") {
		return 0, fmt.Errorf("informe um valor monetário não negativo")
	}

	parts := strings.Split(value, ",")
	if len(parts) > 2 {
		return 0, fmt.Errorf("use o formato 1.234,56")
	}

	whole := parts[0]
	if !validWholeAmount(whole) {
		return 0, fmt.Errorf("use o formato 1.234,56")
	}

	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1]
		if !isDigits(fraction) || len(fraction) == 0 || len(fraction) > 2 {
			return 0, fmt.Errorf("informe até duas casas decimais")
		}
		if len(fraction) == 1 {
			fraction += "0"
		}
	}

	whole = strings.ReplaceAll(whole, ".", "")
	reais, err := strconv.ParseUint(whole, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("valor monetário fora do limite permitido")
	}
	centavos, _ := strconv.ParseUint(fraction, 10, 8)
	maxMoney := uint64(math.MaxInt64)
	if reais > (maxMoney-centavos)/100 {
		return 0, fmt.Errorf("valor monetário fora do limite permitido")
	}

	return Money(reais*100 + centavos), nil
}

func parseDate(value, field string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("data de %s inválida; use AAAA-MM-DD", field)
	}
	return date, nil
}

func validWholeAmount(value string) bool {
	groups := strings.Split(value, ".")
	if len(groups) == 1 {
		return isDigits(groups[0])
	}
	if len(groups[0]) < 1 || len(groups[0]) > 3 || !isDigits(groups[0]) {
		return false
	}
	for _, group := range groups[1:] {
		if len(group) != 3 || !isDigits(group) {
			return false
		}
	}
	return true
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}
