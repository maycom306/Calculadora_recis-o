package rescisao

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

// ResultadoRescisao contains the money items and day-work deduction.
type ResultadoRescisao struct {
	DecimoTerceiro    Money `json:"decimoTerceiroCentavos"`
	Ferias            Money `json:"feriasCentavos"`
	TercoFerias       Money `json:"tercoFeriasCentavos"`
	DiasTrabalhados   Money `json:"diasTrabalhadosCentavos"`
	AvosFeriasDevidos int   `json:"avosFeriasDevidos"`
	Total             Money `json:"totalCentavos"`
}

// CalcularValoresRescisao sums the confirmed monetary items and accepts vacation avos already paid.
func CalcularValoresRescisao(salarioBaseDecimo Money, avosDecimoDisponiveis int, avosDecimoPagos int, baseFerias Money, avosFeriasDisponiveis int, avosFeriasPagos int, baseDiasTrabalhados Money, quantidadeDiasTrabalhados int) (ResultadoRescisao, error) {
	decimo, err := CalcularDecimoProporcional(salarioBaseDecimo, avosDecimoDisponiveis, avosDecimoPagos)
	if err != nil {
		return ResultadoRescisao{}, fmt.Errorf("13º proporcional: %w", err)
	}

	avosFeriasDevidos, err := CalcularSaldoAvosFerias(avosFeriasDisponiveis, avosFeriasPagos)
	if err != nil {
		return ResultadoRescisao{}, err
	}

	ferias, err := calcularPorDozeAvos(baseFerias, avosFeriasDevidos)
	if err != nil {
		return ResultadoRescisao{}, fmt.Errorf("férias proporcionais: %w", err)
	}

	terco, err := CalcularTercoFerias(ferias)
	if err != nil {
		return ResultadoRescisao{}, err
	}

	diasTrabalhados, err := CalcularValorDiasTrabalhados(baseDiasTrabalhados, quantidadeDiasTrabalhados)
	if err != nil {
		return ResultadoRescisao{}, err
	}

	if decimo > Money(math.MaxInt64)-ferias || decimo+ferias > Money(math.MaxInt64)-terco || decimo+ferias+terco > Money(math.MaxInt64)-diasTrabalhados {
		return ResultadoRescisao{}, fmt.Errorf("total da rescisão fora do limite permitido")
	}

	return ResultadoRescisao{
		DecimoTerceiro:    decimo,
		Ferias:            ferias,
		TercoFerias:       terco,
		DiasTrabalhados:   diasTrabalhados,
		AvosFeriasDevidos: avosFeriasDevidos,
		Total:             decimo + ferias + terco,
	}, nil
}

// CalcularDecimoProporcional calcula o valor dos avos adquiridos ainda não pagos.
func CalcularDecimoProporcional(salarioBase Money, avosDisponiveis int, avosJaPagos int) (Money, error) {
	if avosDisponiveis < 0 {
		return 0, fmt.Errorf("avos do 13º adquiridos não podem ser negativos")
	}
	if avosJaPagos < 0 {
		return 0, fmt.Errorf("avos do 13º já recebidos não podem ser negativos")
	}
	avosDevidos := avosDisponiveis - avosJaPagos
	if avosDevidos < 0 {
		avosDevidos = 0
	}
	return calcularPorDozeAvos(salarioBase, avosDevidos)
}

// CalcularAvosDecimo grants one avo after 15 days in each monthly period
// anchored to the admission date, accumulating across admission anniversaries.
func CalcularAvosDecimo(dataAdmissao string, dataDesligamento string) (int, error) {
	admissao, err := parseDate(dataAdmissao, "admissão")
	if err != nil {
		return 0, err
	}
	desligamento, err := parseDate(dataDesligamento, "desligamento")
	if err != nil {
		return 0, err
	}
	if desligamento.Before(admissao) {
		return 0, fmt.Errorf("a data de desligamento não pode ser anterior à admissão")
	}

	avos := 0
	for mes := 0; ; mes++ {
		inicioPeriodo := aniversarioMensal(admissao, mes)
		if inicioPeriodo.After(desligamento) {
			break
		}
		if !inicioPeriodo.AddDate(0, 0, 14).After(desligamento) {
			avos++
		}
	}
	return avos, nil
}

// CalcularSaldoAvosFerias removes vacation avos already taken or paid.
func CalcularSaldoAvosFerias(avosDisponiveis int, avosJaPagos int) (int, error) {
	if avosDisponiveis < 0 || avosJaPagos < 0 {
		return 0, fmt.Errorf("avos de férias não podem ser negativos")
	}
	if avosJaPagos > avosDisponiveis {
		return 0, fmt.Errorf("avos de férias já pagos excedem os avos disponíveis")
	}
	return avosDisponiveis - avosJaPagos, nil
}

// CalcularAvosFerias uses the standard accrual rule: each 30 consecutive days gives 1 avo,
// and any residual of 15 days or more also adds 1 avo.
func CalcularAvosFerias(dataAdmissao string, dataDesligamento string) (int, error) {
	admissao, err := parseDate(dataAdmissao, "admissão")
	if err != nil {
		return 0, err
	}
	desligamento, err := parseDate(dataDesligamento, "desligamento")
	if err != nil {
		return 0, err
	}
	if desligamento.Before(admissao) {
		return 0, fmt.Errorf("a data de desligamento não pode ser anterior à admissão")
	}

	diasTrabalhados := int(desligamento.Sub(admissao).Hours()/24) + 1
	avos := diasTrabalhados / 30
	if diasTrabalhados%30 >= 15 {
		avos++
	}
	return avos, nil
}

// CalcularValorDiasTrabalhados uses the confirmed fixed divisor of 30.
func CalcularValorDiasTrabalhados(baseMensal Money, dias int) (Money, error) {
	if dias < 0 || dias > 31 {
		return 0, fmt.Errorf("quantidade de dias trabalhados deve estar entre 0 e 31")
	}
	return calcularProporcional(baseMensal, dias, 30)
}

// CalcularTercoFerias calculates one third of the already-rounded vacation amount.
func CalcularTercoFerias(ferias Money) (Money, error) {
	if ferias < 0 {
		return 0, fmt.Errorf("valor das férias não pode ser negativo")
	}

	terco := ferias / 3
	if ferias%3 >= 2 {
		terco++
	}
	return terco, nil
}

func calcularPorDozeAvos(base Money, avos int) (Money, error) {
	if base < 0 {
		return 0, fmt.Errorf("remuneração-base não pode ser negativa")
	}
	if avos < 0 {
		return 0, fmt.Errorf("quantidade de avos não pode ser negativa")
	}

	return calcularProporcional(base, avos, 12)
}

func calcularProporcional(base Money, quantidade int, divisor int64) (Money, error) {
	if base < 0 {
		return 0, fmt.Errorf("base monetária não pode ser negativa")
	}
	if quantidade < 0 {
		return 0, fmt.Errorf("quantidade não pode ser negativa")
	}

	numerador := new(big.Int).Mul(big.NewInt(int64(base)), big.NewInt(int64(quantidade)))
	quociente, resto := new(big.Int), new(big.Int)
	quociente.QuoRem(numerador, big.NewInt(divisor), resto)
	if new(big.Int).Mul(resto, big.NewInt(2)).Cmp(big.NewInt(divisor)) >= 0 {
		quociente.Add(quociente, big.NewInt(1))
	}
	if !quociente.IsInt64() {
		return 0, fmt.Errorf("resultado monetário fora do limite permitido")
	}
	return Money(quociente.Int64()), nil
}

func aniversarioMensal(admissao time.Time, meses int) time.Time {
	mesDestino := time.Date(admissao.Year(), admissao.Month()+time.Month(meses), 1, 0, 0, 0, 0, time.UTC)
	ultimoDia := time.Date(mesDestino.Year(), mesDestino.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	dia := admissao.Day()
	if dia > ultimoDia {
		dia = ultimoDia
	}
	return time.Date(mesDestino.Year(), mesDestino.Month(), dia, 0, 0, 0, 0, time.UTC)
}
