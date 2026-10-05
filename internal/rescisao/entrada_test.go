package rescisao

import (
	"testing"
)

func TestParseMoney(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Money
		valid bool
	}{
		{name: "inteiro", input: "1500", want: 150000, valid: true},
		{name: "centavos brasileiros", input: "1.234,56", want: 123456, valid: true},
		{name: "uma casa decimal", input: "10,5", want: 1050, valid: true},
		{name: "prefixo real", input: "R$ 0,05", want: 5, valid: true},
		{name: "zero", input: "0,00", want: 0, valid: true},
		{name: "agrupamento inválido", input: "12.34,56", valid: false},
		{name: "decimais demais", input: "10,999", valid: false},
		{name: "negativo", input: "-10,00", valid: false},
		{name: "vazio", input: "", valid: false},
		{name: "valor acima do limite", input: "92233720368547758,08", valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseMoney(test.input)
			if test.valid {
				if err != nil {
					t.Fatalf("ParseMoney(%q) returned error: %v", test.input, err)
				}
				if got != test.want {
					t.Errorf("ParseMoney(%q) = %d, want %d", test.input, got, test.want)
				}
				return
			}
			if err == nil {
				t.Errorf("ParseMoney(%q) expected an error", test.input)
			}
		})
	}
}

func TestValidarEntrada(t *testing.T) {
	valid := EntradaRescisao{
		Nome:             "  Ana da Silva  ",
		Matricula:        "  123  ",
		Cargo:            "  Técnica  ",
		DataAdmissao:     "2024-02-29",
		DataDesligamento: "2024-03-01",
		SalarioBase:      "1.234,56",
		AvosFeriasPagos:  "0",
	}

	got, err := ValidarEntrada(valid)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.Servidor.Nome != "Ana da Silva" || got.Servidor.Matricula != "123" || got.Servidor.Cargo != "Técnica" {
		t.Errorf("ValidarEntrada() did not trim identification fields: %+v", got.Servidor)
	}
	if got.SalarioBase != 123456 {
		t.Errorf("SalarioBase = %d, want 123456 cents", got.SalarioBase)
	}
	if got.AvosFeriasPagos != 0 {
		t.Errorf("AvosFeriasPagos = %d, want 0", got.AvosFeriasPagos)
	}
	if got.DataAdmissao != "2024-02-29" {
		t.Errorf("DataAdmissao = %q, want 2024-02-29", got.DataAdmissao)
	}
}

func TestValidarEntradaContaSomenteDiasDoMesDeDesligamento(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2024-01-10",
		DataDesligamento: "2024-03-03",
		SalarioBase:      "3.000,00",
	}

	got, err := ValidarEntrada(input)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.DiasTrabalhados != 3 {
		t.Errorf("DiasTrabalhados = %d, want 3", got.DiasTrabalhados)
	}
}

func TestValidarEntradaRetornaAvosFeriasDisponiveis(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2025-01-31",
		DataDesligamento: "2025-02-28",
		SalarioBase:      "3.000,00",
		AvosFeriasPagos:  "0",
	}

	got, err := ValidarEntrada(input)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.AvosFeriasDisponiveis != 1 {
		t.Errorf("AvosFeriasDisponiveis = %d, want 1", got.AvosFeriasDisponiveis)
	}
	if got.AvosFeriasDevidos != 1 {
		t.Errorf("AvosFeriasDevidos = %d, want 1", got.AvosFeriasDevidos)
	}
}

func TestValidarEntradaCalculaAvosDecimoEDescontaPagos(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2025-01-01",
		DataDesligamento: "2025-02-15",
		SalarioBase:      "3.000,00",
		AvosDecimo:       "1",
	}

	got, err := ValidarEntrada(input)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.AvosDecimoDisponiveis != 2 || got.AvosDecimo != 1 || got.AvosDecimoDevidos != 1 {
		t.Errorf("13th avos = available %d, paid %d, due %d; want 2, 1, 1", got.AvosDecimoDisponiveis, got.AvosDecimo, got.AvosDecimoDevidos)
	}
}

func TestValidarEntradaPermiteAvosDecimoPagosAcimaDosAdquiridos(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2025-01-01",
		DataDesligamento: "2025-06-15",
		SalarioBase:      "3.000,00",
		AvosDecimo:       "8",
	}

	got, err := ValidarEntrada(input)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.AvosDecimoDisponiveis != 6 || got.AvosDecimo != 8 || got.AvosDecimoDevidos != 0 {
		t.Errorf("13th avos = available %d, paid %d, due %d; want 6, 8, 0", got.AvosDecimoDisponiveis, got.AvosDecimo, got.AvosDecimoDevidos)
	}
}

func TestValidarEntradaCalculaDezoitoAvosDeDecimoEDescontaSeisPagos(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2025-01-01",
		DataDesligamento: "2026-06-15",
		SalarioBase:      "2.000,00",
		AvosDecimo:       "6",
	}

	got, err := ValidarEntrada(input)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.AvosDecimoDisponiveis != 18 || got.AvosDecimo != 6 || got.AvosDecimoDevidos != 12 {
		t.Errorf("13th avos = available %d, paid %d, due %d; want 18, 6, 12", got.AvosDecimoDisponiveis, got.AvosDecimo, got.AvosDecimoDevidos)
	}
}

func TestValidarEntradaRejeitaFeriasPagasAcimaDoDisponivel(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2025-01-31",
		DataDesligamento: "2025-02-28",
		SalarioBase:      "3.000,00",
		AvosFeriasPagos:  "2",
	}
	if _, err := ValidarEntrada(input); err == nil {
		t.Fatal("ValidarEntrada() expected an error when paid vacation avos exceed available avos")
	}
}

func TestValidarEntradaSubtraiAvosDeFeriasJaPagos(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2024-01-01",
		DataDesligamento: "2025-01-01",
		SalarioBase:      "3.000,00",
		AvosFeriasPagos:  "12",
	}

	got, err := ValidarEntrada(input)
	if err != nil {
		t.Fatalf("ValidarEntrada() returned error: %v", err)
	}
	if got.AvosFeriasDisponiveis != 12 || got.AvosFeriasPagos != 12 || got.AvosFeriasDevidos != 0 {
		t.Errorf("vacation avos = available %d, paid %d, due %d; want 12, 12, 0", got.AvosFeriasDisponiveis, got.AvosFeriasPagos, got.AvosFeriasDevidos)
	}
}

func TestValidarEntradaErrors(t *testing.T) {
	valid := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2024-01-10",
		DataDesligamento: "2024-02-10",
		SalarioBase:      "100,00",
	}

	tests := []struct {
		name   string
		change func(*EntradaRescisao)
	}{
		{name: "nome obrigatório", change: func(input *EntradaRescisao) { input.Nome = " " }},
		{name: "data inválida", change: func(input *EntradaRescisao) { input.DataAdmissao = "2024-02-30" }},
		{name: "desligamento anterior", change: func(input *EntradaRescisao) { input.DataDesligamento = "2024-01-09" }},
		{name: "salário inválido", change: func(input *EntradaRescisao) { input.SalarioBase = "-1,00" }},
		{name: "avos de férias negativos", change: func(input *EntradaRescisao) { input.AvosFeriasPagos = "-1" }},
		{name: "avos de férias não inteiros", change: func(input *EntradaRescisao) { input.AvosFeriasPagos = "1,5" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := valid
			test.change(&input)
			if _, err := ValidarEntrada(input); err == nil {
				t.Fatal("ValidarEntrada() expected an error")
			}
		})
	}
}

func TestValidarEntradaPermiteDatasIguais(t *testing.T) {
	input := EntradaRescisao{
		Nome:             "Ana",
		DataAdmissao:     "2024-01-10",
		DataDesligamento: "2024-01-10",
		SalarioBase:      "100,00",
	}
	if _, err := ValidarEntrada(input); err != nil {
		t.Fatalf("same admission and termination date should be accepted: %v", err)
	}
}
