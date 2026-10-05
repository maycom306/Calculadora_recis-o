package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// Documento representa os dados que devem constar no PDF final do termo.
type Documento struct {
	Nome                 string
	Matricula            string
	Cargo                string
	Banco                string
	Agencia              string
	Conta                string
	Pix                  string
	DataAdmissao         string
	DataDesligamento     string
	SalarioBase          string
	DiasTrabalhados      int
	ValorDiasTrabalhados int64
	DecimoTerceiro       int64
	Ferias               int64
	TercoFerias          int64
	Total                int64
	AvosFeriasDevidos    int
}

// GerarArquivo cria um PDF profissional em A4 com layout institucional.
func GerarArquivo(doc Documento, destino string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Termo de Rescisão", false)
	pdf.SetAuthor("RescisaoDesk", false)
	pdf.SetCreator("RescisaoDesk", false)
	pdf.SetMargins(12, 10, 12)
	pdf.SetAutoPageBreak(true, 10)
	pdf.AddPage()

	configureFonts(pdf)
	layoutTermo(pdf, doc)
	if pdf.PageNo() != 1 {
		return "", fmt.Errorf("o termo excedeu o limite de uma página")
	}

	if err := pdf.OutputFileAndClose(destino); err != nil {
		return "", err
	}
	return destino, nil
}

func configureFonts(pdf *gofpdf.Fpdf) {
	candidates := []string{
		"C:/Windows/Fonts/arial.ttf",
		"C:/Windows/Fonts/calibri.ttf",
		"C:/Windows/Fonts/arialbd.ttf",
		"C:/Windows/Fonts/calibrib.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation2/LiberationSans-Regular.ttf",
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			pdf.AddUTF8Font("Arial", "", candidate)
			pdf.AddUTF8Font("Arial", "B", candidate)
			pdf.SetFont("Arial", "", 11)
			return
		}
	}
	pdf.SetFont("Helvetica", "", 11)
}

func layoutTermo(pdf *gofpdf.Fpdf, doc Documento) {
	pdf.SetFillColor(12, 52, 98)
	pdf.Rect(12, 10, 186, 28, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 14)
	pdf.Text(17, 21, "ESTADO DO RIO GRANDE DO NORTE")
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(255, 255, 255)
	pdf.Text(17, 28, "PREFEITURA MUNICIPAL DE IPANGUAÇU")
	pdf.Text(17, 33, "SECRETARIA MUNICIPAL DE PLANEJAMENTO, ADMINISTRAÇÃO E TRANSPORTES")

	pdf.SetTextColor(30, 30, 30)
	pdf.SetY(43)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, "TERMO DE RESCISÃO", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	printSectionHeading(pdf, "DADOS DO SERVIDOR")
	printCamposEmColunas(pdf, [][2]string{
		{"Servidor", doc.Nome},
		{"Matrícula", strings.TrimSpace(doc.Matricula)},
		{"Cargo", strings.TrimSpace(doc.Cargo)},
		{"Admissão", formatData(doc.DataAdmissao)},
		{"Desligamento", formatData(doc.DataDesligamento)},
		{"Salário-base", formatMoeda(doc.SalarioBase)},
		{"Dias do mês do desligamento", fmt.Sprintf("%d", doc.DiasTrabalhados)},
	})

	printSectionHeading(pdf, "DADOS PARA PAGAMENTO")
	printCamposEmColunas(pdf, [][2]string{
		{"Banco", strings.TrimSpace(doc.Banco)},
		{"Agência", strings.TrimSpace(doc.Agencia)},
		{"Conta", strings.TrimSpace(doc.Conta)},
		{"Pix", strings.TrimSpace(doc.Pix)},
	})

	pdf.Ln(2)
	printSectionHeading(pdf, "VERBAS RESCISÓRIAS")
	printValor(pdf, "13º proporcional devido (após desconto dos avos recebidos)", formatMoedaInteiro(doc.DecimoTerceiro))
	printValor(pdf, "Férias proporcionais devidas (após desconto dos avos pagos)", formatMoedaInteiro(doc.Ferias))
	printValor(pdf, "Terço constitucional de férias", formatMoedaInteiro(doc.TercoFerias))
	pdf.Ln(1)
	printValorTotal(pdf, "TOTAL DAS VERBAS RESCISÓRIAS", formatMoedaInteiro(doc.Total))

	pdf.Ln(2)
	printSectionHeading(pdf, "VALOR SEPARADO (NÃO INCLUÍDO NO TOTAL)")
	printValor(pdf, fmt.Sprintf("Dias trabalhados (%d dias)", doc.DiasTrabalhados), formatMoedaInteiro(doc.ValorDiasTrabalhados))

	pdf.Ln(2)
	pdf.SetFont("Arial", "", 8)
	pdf.MultiCell(0, 4, "Cálculo realizado com os dados informados no sistema.", "", "L", false)

	pdf.Ln(3)
	ensureSpace(pdf, 18)
	pdf.SetDrawColor(180, 180, 180)
	y := pdf.GetY() + 7
	pdf.Line(30, y, 95, y)
	pdf.Line(120, y, 185, y)
	pdf.SetY(y + 1)
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(65, 4, "Assinatura do servidor", "", 0, "C", false, 0, "")
	pdf.CellFormat(65, 4, "Responsável / setor", "", 0, "C", false, 0, "")
}

func printSectionHeading(pdf *gofpdf.Fpdf, title string) {
	ensureSpace(pdf, 15)
	y := pdf.GetY()
	pdf.SetFillColor(235, 240, 245)
	pdf.Rect(12, y, 186, 7, "F")
	pdf.SetTextColor(30, 30, 30)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(17, y+0.5)
	pdf.CellFormat(0, 6, title, "", 0, "L", false, 0, "")
	pdf.SetY(y + 9)
}

func printCamposEmColunas(pdf *gofpdf.Fpdf, campos [][2]string) {
	const larguraColuna = 90.0
	const espacoColunas = 6.0
	const alturaLinha = 4.5

	for i := 0; i < len(campos); i += 2 {
		left := campos[i]
		leftValue := strings.TrimSpace(left[1])
		if leftValue == "" {
			leftValue = "-"
		}

		right := [2]string{}
		rightLines := []string{""}
		if i+1 < len(campos) {
			right = campos[i+1]
			right[1] = strings.TrimSpace(right[1])
			if right[1] == "" {
				right[1] = "-"
			}
			pdf.SetFont("Arial", "", 10)
			rightLines = pdf.SplitText(right[1], larguraColuna)
		}

		pdf.SetFont("Arial", "", 10)
		leftLines := pdf.SplitText(leftValue, larguraColuna)
		lineCount := len(leftLines)
		if len(rightLines) > lineCount {
			lineCount = len(rightLines)
		}
		rowHeight := 5 + float64(lineCount)*alturaLinha
		ensureSpace(pdf, rowHeight+2)
		y := pdf.GetY()
		printCampoColuna(pdf, 12, y, larguraColuna, left[0], leftValue, alturaLinha)
		if i+1 < len(campos) {
			printCampoColuna(pdf, 12+larguraColuna+espacoColunas, y, larguraColuna, right[0], right[1], alturaLinha)
		}
		pdf.SetDrawColor(225, 229, 233)
		pdf.Line(12, y+rowHeight-1, 198, y+rowHeight-1)
		pdf.SetY(y + rowHeight)
	}
}

func printCampoColuna(pdf *gofpdf.Fpdf, x, y, width float64, label, value string, lineHeight float64) {
	pdf.SetXY(x, y)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(width, 4.5, label+":", "", 1, "L", false, 0, "")
	pdf.SetXY(x, y+4.5)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(width, lineHeight, value, "", "L", false)
}

func printValor(pdf *gofpdf.Fpdf, label, value string) {
	ensureSpace(pdf, 8)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(125, 6.5, label+":", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6.5, value, "", 1, "R", false, 0, "")
}

func printValorTotal(pdf *gofpdf.Fpdf, label, value string) {
	ensureSpace(pdf, 9)
	pdf.SetFillColor(231, 240, 233)
	pdf.SetTextColor(28, 78, 45)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetX(15)
	pdf.CellFormat(120, 8, label, "", 0, "L", true, 0, "")
	pdf.CellFormat(55, 8, value, "", 1, "R", true, 0, "")
	pdf.SetTextColor(30, 30, 30)
}

func ensureSpace(pdf *gofpdf.Fpdf, height float64) {
	_, pageHeight := pdf.GetPageSize()
	if pdf.GetY()+height > pageHeight-18 {
		pdf.AddPage()
	}
}

func formatData(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "-"
	}
	parts := strings.Split(value, "-")
	if len(parts) == 3 {
		return parts[2] + "/" + parts[1] + "/" + parts[0]
	}
	return value
}

func formatMoeda(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "R$ 0,00"
	}
	return value
}

func formatMoedaInteiro(value int64) string {
	reais := value / 100
	centavos := value % 100
	if centavos < 0 {
		reais = (value + 1) / 100
		centavos = value % 100
		if centavos < 0 {
			centavos = -centavos
		}
	}
	if centavos < 0 {
		centavos = -centavos
	}
	return fmt.Sprintf("R$ %d,%02d", reais, centavos)
}

func encodeWinAnsiText(value string) string {
	return value
}

// NomeArquivo gera um nome seguro para o arquivo do PDF.
func NomeArquivo(nome, matricula string) string {
	base := strings.TrimSpace(matricula)
	if base == "" {
		base = strings.TrimSpace(nome)
	}
	if base == "" {
		base = "termo_rescisao"
	}
	base = strings.ToUpper(base)
	base = regexp.MustCompile(`[^A-Za-z0-9_\-]+`).ReplaceAllString(base, "_")
	base = strings.Trim(base, "_")
	if base == "" {
		return "termo_rescisao"
	}
	return base + "_TERMO_RESCISAO"
}
