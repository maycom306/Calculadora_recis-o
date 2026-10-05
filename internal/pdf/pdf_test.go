package pdf

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestGerarArquivo(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "arquivo.pdf")

	doc := Documento{
		Nome:                 strings.Repeat("Ana da Silva ", 11),
		Matricula:            strings.Repeat("1234", 10),
		Cargo:                strings.Repeat("Técnica administrativa ", 5),
		Banco:                strings.Repeat("Banco ", 12),
		Agencia:              strings.Repeat("12345", 4),
		Conta:                strings.Repeat("12345-", 5),
		Pix:                  strings.Repeat("chave-pix-", 12),
		DataAdmissao:         "2024-01-10",
		DataDesligamento:     "2025-01-10",
		SalarioBase:          "R$ 3.000,00",
		DiasTrabalhados:      10,
		ValorDiasTrabalhados: 150000,
		DecimoTerceiro:       200000,
		Ferias:               150000,
		TercoFerias:          50000,
		Total:                400000,
		AvosFeriasDevidos:    1,
	}

	path, err := GerarArquivo(doc, dest)
	if err != nil {
		t.Fatalf("GerarArquivo() returned error: %v", err)
	}
	if path != dest {
		t.Fatalf("GerarArquivo() path = %q, want %q", path, dest)
	}

	bytes, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}
	if !strings.HasPrefix(string(bytes), "%PDF") {
		t.Fatalf("generated pdf does not contain PDF header")
	}
	if len(bytes) < 512 {
		t.Fatalf("generated pdf is too small to be a valid document")
	}
	if !strings.Contains(string(bytes), "%%EOF") {
		t.Fatalf("generated pdf does not contain the final PDF terminator")
	}
	pageCount := regexp.MustCompile(`/Type\s*/Page\b`).FindAll(bytes, -1)
	if len(pageCount) != 1 {
		t.Fatalf("generated pdf has %d pages, want exactly 1", len(pageCount))
	}
}
