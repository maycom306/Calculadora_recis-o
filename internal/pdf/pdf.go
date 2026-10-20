package pdf

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

//go:embed assets/logo-ipanguacu.png
var logoIpanguacu []byte

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
	pdf.SetFillColor(25, 63, 101)
	pdf.Rect(12, 10, 186, 34, "F")
	pdf.SetFillColor(190, 44, 53)
	pdf.Rect(12, 10, 186, 1.2, "F")

	imageOptions := gofpdf.ImageOptions{ImageType: "PNG"}
	pdf.RegisterImageOptionsReader("logo-ipanguacu", imageOptions, bytes.NewReader(logoIpanguacu))
	pdf.ImageOptions("logo-ipanguacu", 16, 15, 63, 0, false, imageOptions, 0, "")

	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8)
	pdf.Text(84, 18.5, "ESTADO DO RIO GRANDE DO NORTE")
	pdf.SetFont("Arial", "B", 9)
	pdf.Text(84, 24.5, "PREFEITURA MUNICIPAL DE IPANGUAÇU")
	pdf.SetFont("Arial", "", 6.8)
	pdf.Text(84, 31, "SECRETARIA MUNICIPAL DE PLANEJAMENTO, ADMINISTRAÇÃO E TRANSPORTES")
	pdf.SetFont("Arial", "B", 7.5)
	pdf.Text(84, 37, "DEPARTAMENTO DE RECURSOS HUMANOS")

	pdf.SetTextColor(30, 30, 30)
	pdf.SetY(48)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, "TERMO DE RESCISÃO", "", 1, "C", false, 0, "")
	pdf.SetDrawColor(190, 44, 53)
	pdf.SetLineWidth(0.6)
	pdf.Line(84, 57, 126, 57)
	pdf.Ln(6)

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
	pdf.Ln(2.5)
	printValorTotal(pdf, "TOTAL DAS VERBAS RESCISÓRIAS", formatMoedaInteiro(doc.Total))

	pdf.Ln(2.5)
	printSectionHeading(pdf, "VALOR SEPARADO (NÃO INCLUÍDO NO TOTAL)")
	printValor(pdf, fmt.Sprintf("Dias trabalhados (%d dias)", doc.DiasTrabalhados), formatMoedaInteiro(doc.ValorDiasTrabalhados))

	pdf.Ln(2.5)
	pdf.SetTextColor(105, 116, 128)
	pdf.SetFont("Arial", "", 7.5)
	pdf.MultiCell(0, 4, "Cálculo realizado com os dados informados no sistema.", "", "L", false)
	pdf.SetTextColor(30, 30, 30)

	printAssinaturas(pdf)
}

func printSectionHeading(pdf *gofpdf.Fpdf, title string) {
	ensureSpace(pdf, 12)
	y := pdf.GetY()
	pdf.SetFillColor(232, 239, 246)
	pdf.Rect(12, y, 186, 7.5, "F")
	pdf.SetFillColor(38, 91, 137)
	pdf.Rect(12, y, 1.2, 7.5, "F")
	pdf.SetTextColor(25, 63, 101)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.SetXY(17, y+0.5)
	pdf.CellFormat(0, 6.5, title, "", 0, "L", false, 0, "")
	pdf.SetY(y + 9)
	pdf.SetTextColor(30, 30, 30)
}

func printCamposEmColunas(pdf *gofpdf.Fpdf, campos [][2]string) {
	const larguraColuna = 90.0
	const espacoColunas = 6.0
	const alturaLinha = 4.4

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
		rowHeight := 5.5 + float64(lineCount)*alturaLinha
		ensureSpace(pdf, rowHeight+1)
		y := pdf.GetY()
		pdf.SetFillColor(248, 250, 252)
		pdf.Rect(12, y, larguraColuna, rowHeight-1, "F")
		if i+1 < len(campos) {
			pdf.Rect(12+larguraColuna+espacoColunas, y, larguraColuna, rowHeight-1, "F")
		}
		printCampoColuna(pdf, 12, y, larguraColuna, left[0], leftValue, alturaLinha)
		if i+1 < len(campos) {
			printCampoColuna(pdf, 12+larguraColuna+espacoColunas, y, larguraColuna, right[0], right[1], alturaLinha)
		}
		pdf.SetDrawColor(216, 224, 232)
		pdf.Line(12, y+rowHeight-1, 198, y+rowHeight-1)
		pdf.SetY(y + rowHeight + 0.5)
	}
}

func printCampoColuna(pdf *gofpdf.Fpdf, x, y, width float64, label, value string, lineHeight float64) {
	pdf.SetXY(x, y)
	pdf.SetTextColor(38, 91, 137)
	pdf.SetFont("Arial", "B", 7.5)
	pdf.CellFormat(width, 4.2, label+":", "", 1, "L", false, 0, "")
	pdf.SetXY(x, y+4.2)
	pdf.SetTextColor(35, 43, 52)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(width, lineHeight, value, "", "L", false)
	pdf.SetTextColor(30, 30, 30)
}

func printValor(pdf *gofpdf.Fpdf, label, value string) {
	const labelWidth = 147.0
	lineHeight := 4.7
	pdf.SetFont("Arial", "B", 8.5)
	lineCount := len(pdf.SplitText(label, labelWidth))
	rowHeight := float64(lineCount) * lineHeight
	if rowHeight < 8 {
		rowHeight = 8
	}
	ensureSpace(pdf, rowHeight+1)
	y := pdf.GetY()
	pdf.SetFillColor(248, 250, 252)
	pdf.Rect(12, y, 186, rowHeight, "F")
	pdf.SetTextColor(42, 51, 61)
	pdf.SetXY(15, y+(rowHeight-float64(lineCount)*lineHeight)/2)
	pdf.MultiCell(labelWidth, lineHeight, label, "", "L", false)
	pdf.SetTextColor(25, 63, 101)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(162, y)
	pdf.CellFormat(33, rowHeight, value, "", 0, "R", false, 0, "")
	pdf.SetY(y + rowHeight + 1)
	pdf.SetTextColor(30, 30, 30)
}

func printValorTotal(pdf *gofpdf.Fpdf, label, value string) {
	ensureSpace(pdf, 9)
	pdf.SetFillColor(25, 63, 101)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9.5)
	pdf.SetX(15)
	pdf.CellFormat(120, 8, label, "", 0, "L", true, 0, "")
	pdf.CellFormat(55, 8, value, "", 1, "R", true, 0, "")
	pdf.SetTextColor(30, 30, 30)
	pdf.SetLineWidth(0.2)
}

func printAssinaturas(pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(190, 44, 53)
	pdf.SetLineWidth(0.5)
	pdf.Line(12, 263, 198, 263)
	pdf.SetTextColor(25, 63, 101)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetXY(12, 266)
	pdf.CellFormat(186, 5, "ASSINATURAS", "", 0, "L", false, 0, "")

	pdf.SetDrawColor(145, 155, 165)
	pdf.SetLineWidth(0.25)
	pdf.Line(24, 280, 96, 280)
	pdf.Line(114, 280, 186, 280)
	pdf.SetTextColor(82, 91, 101)
	pdf.SetFont("Arial", "", 7.5)
	pdf.SetXY(24, 282)
	pdf.CellFormat(72, 4, "Assinatura do servidor", "", 0, "C", false, 0, "")
	pdf.SetXY(114, 282)
	pdf.CellFormat(72, 4, "Responsável / setor", "", 0, "C", false, 0, "")
	pdf.SetTextColor(30, 30, 30)
	pdf.SetLineWidth(0.2)
}

func ensureSpace(pdf *gofpdf.Fpdf, height float64) {
	_, pageHeight := pdf.GetPageSize()
	if pdf.GetY()+height > pageHeight-36 {
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
