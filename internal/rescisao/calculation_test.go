package rescisao

import "testing"

func TestCalcularDecimoProporcionalSalarioBase3000OitoAvosJaRecebidos(t *testing.T) {
	got, err := CalcularDecimoProporcional(300000, 12, 8)
	if err != nil {
		t.Fatalf("CalcularDecimoProporcional() returned error: %v", err)
	}
	if got != 100000 {
		t.Errorf("CalcularDecimoProporcional() = %d, want 100000 cents", got)
	}
}

func TestCalcularDecimoProporcionalAvosPagosAcimaDosAdquiridos(t *testing.T) {
	got, err := CalcularDecimoProporcional(200000, 6, 8)
	if err != nil {
		t.Fatalf("CalcularDecimoProporcional() returned error: %v", err)
	}
	if got != 0 {
		t.Errorf("CalcularDecimoProporcional() = %d, want 0 cents", got)
	}
}

func TestCalcularDecimoProporcionalDezoitoAvosMenosSeisPagos(t *testing.T) {
	got, err := CalcularDecimoProporcional(200000, 18, 6)
	if err != nil {
		t.Fatalf("CalcularDecimoProporcional() returned error: %v", err)
	}
	if got != 200000 {
		t.Errorf("CalcularDecimoProporcional() = %d, want 200000 cents", got)
	}
}

func TestCalcularValoresRescisao(t *testing.T) {
	got, err := CalcularValoresRescisao(300000, 12, 6, 300000, 18, 12, 300000, 10)
	if err != nil {
		t.Fatalf("CalcularValoresRescisao() returned error: %v", err)
	}

	if got.DecimoTerceiro != 150000 {
		t.Errorf("DecimoTerceiro = %d, want 150000 cents", got.DecimoTerceiro)
	}
	if got.AvosFeriasDevidos != 6 {
		t.Errorf("AvosFeriasDevidos = %d, want 6", got.AvosFeriasDevidos)
	}
	if got.Ferias != 150000 {
		t.Errorf("Ferias = %d, want 150000 cents", got.Ferias)
	}
	if got.TercoFerias != 50000 {
		t.Errorf("TercoFerias = %d, want 50000 cents", got.TercoFerias)
	}
	if got.DiasTrabalhados != 100000 {
		t.Errorf("DiasTrabalhados = %d, want 100000 cents", got.DiasTrabalhados)
	}
	if got.Total != 350000 {
		t.Errorf("Total = %d, want 350000 cents", got.Total)
	}
}

func TestCalcularValoresRescisaoCasoInformadoOutubro2026(t *testing.T) {
	avosDisponiveis, err := CalcularAvosFerias("2025-05-06", "2026-10-31")
	if err != nil {
		t.Fatalf("CalcularAvosFerias() returned error: %v", err)
	}
	if avosDisponiveis != 18 {
		t.Fatalf("AvosFeriasDisponiveis = %d, want 18", avosDisponiveis)
	}

	got, err := CalcularValoresRescisao(200000, 12, 11, 200000, avosDisponiveis, 7, 200000, 30)
	if err != nil {
		t.Fatalf("CalcularValoresRescisao() returned error: %v", err)
	}
	if got.DecimoTerceiro != 16667 {
		t.Errorf("DecimoTerceiro = %d, want 16667 cents", got.DecimoTerceiro)
	}
	if got.AvosFeriasDevidos != 11 {
		t.Errorf("AvosFeriasDevidos = %d, want 11", got.AvosFeriasDevidos)
	}
	if got.Ferias != 183333 {
		t.Errorf("Ferias = %d cents, want 183333 cents (R$ 1.833,33)", got.Ferias)
	}
	if got.TercoFerias != 61111 {
		t.Errorf("TercoFerias = %d cents, want 61111 cents (R$ 611,11)", got.TercoFerias)
	}
	if got.Total != 261111 {
		t.Errorf("Total = %d cents, want 261111 cents (R$ 2.611,11)", got.Total)
	}
}

func TestCalcularValoresRescisaoDescontaAvosDeDecimoJaPagos(t *testing.T) {
	got, err := CalcularValoresRescisao(200000, 12, 11, 200000, 0, 0, 200000, 30)
	if err != nil {
		t.Fatalf("CalcularValoresRescisao() returned error: %v", err)
	}
	if got.DecimoTerceiro != 16667 {
		t.Errorf("DecimoTerceiro = %d, want 16667 cents", got.DecimoTerceiro)
	}
}

func TestCalcularValoresRescisaoPermiteFeriasAcimaDeDozeAvos(t *testing.T) {
	got, err := CalcularValoresRescisao(300000, 8, 0, 300000, 15, 0, 0, 0)
	if err != nil {
		t.Fatalf("CalcularValoresRescisao() returned error: %v", err)
	}
	if got.Ferias != 375000 {
		t.Errorf("Ferias = %d, want 375000 cents", got.Ferias)
	}
}

func TestCalcularAvosDecimoPeriodoMensalAncoradoNaAdmissao(t *testing.T) {
	tests := []struct {
		dataDesligamento string
		want             int
	}{
		{dataDesligamento: "2025-01-14", want: 0},
		{dataDesligamento: "2025-01-15", want: 1},
		{dataDesligamento: "2025-01-30", want: 1},
		{dataDesligamento: "2025-02-14", want: 1},
		{dataDesligamento: "2025-02-15", want: 2},
		{dataDesligamento: "2025-03-15", want: 3},
		{dataDesligamento: "2026-06-15", want: 18},
	}
	for _, test := range tests {
		t.Run(test.dataDesligamento, func(t *testing.T) {
			got, err := CalcularAvosDecimo("2025-01-01", test.dataDesligamento)
			if err != nil {
				t.Fatalf("CalcularAvosDecimo() returned error: %v", err)
			}
			if got != test.want {
				t.Errorf("CalcularAvosDecimo() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestCalcularAvosFeriasVencimentoPadrao30Dias(t *testing.T) {
	tests := []struct {
		name             string
		dataAdmissao     string
		dataDesligamento string
		want             int
	}{
		{name: "30 dias completos", dataAdmissao: "2025-04-05", dataDesligamento: "2025-05-05", want: 1},
		{name: "45 dias completos", dataAdmissao: "2025-04-05", dataDesligamento: "2025-05-20", want: 2},
		{name: "14 dias não geram avo", dataAdmissao: "2025-02-01", dataDesligamento: "2025-02-14", want: 0},
		{name: "15 dias geram um avo", dataAdmissao: "2025-02-01", dataDesligamento: "2025-02-15", want: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := CalcularAvosFerias(test.dataAdmissao, test.dataDesligamento)
			if err != nil {
				t.Fatalf("CalcularAvosFerias() returned error: %v", err)
			}
			if got != test.want {
				t.Errorf("CalcularAvosFerias() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestCalcularPorDozeAvosArredondaMetadeParaCima(t *testing.T) {
	got, err := calcularPorDozeAvos(10002, 1)
	if err != nil {
		t.Fatalf("calcularPorDozeAvos() returned error: %v", err)
	}
	if got != 834 {
		t.Errorf("calcularPorDozeAvos() = %d, want 834 cents", got)
	}
}

func TestCalcularTercoFeriasArredondaParaCentavos(t *testing.T) {
	got, err := CalcularTercoFerias(101)
	if err != nil {
		t.Fatalf("CalcularTercoFerias() returned error: %v", err)
	}
	if got != 34 {
		t.Errorf("CalcularTercoFerias() = %d, want 34 cents", got)
	}
}

func TestCalcularValoresRescisaoRejeitaAvosNegativos(t *testing.T) {
	if _, err := CalcularValoresRescisao(10000, -1, 0, 10000, 1, 0, 0, 0); err == nil {
		t.Fatal("CalcularValoresRescisao() expected an error for negative avos")
	}
}

func TestCalcularValoresRescisaoRejeitaAvosFeriasPagosAcimaDoDisponivel(t *testing.T) {
	if _, err := CalcularValoresRescisao(10000, 1, 0, 10000, 5, 6, 0, 0); err == nil {
		t.Fatal("CalcularValoresRescisao() expected an error when paid vacation avos exceed available avos")
	}
}

func TestCalcularValorDiasTrabalhadosBase3000DezDias(t *testing.T) {
	got, err := CalcularValorDiasTrabalhados(300000, 10)
	if err != nil {
		t.Fatalf("CalcularValorDiasTrabalhados() returned error: %v", err)
	}
	if got != 100000 {
		t.Errorf("CalcularValorDiasTrabalhados() = %d, want 100000 cents", got)
	}
}
