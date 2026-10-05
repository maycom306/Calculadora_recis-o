# Calculadora_recis-o

Aplicativo desktop Windows para conferência de verbas rescisórias e geração de termo PDF. O nome interno do aplicativo e do módulo Go é `RescisaoDesk`.

## Estado atual

- Formulário em quatro etapas: dados do servidor, rescisão, dados bancários e conferência.
- Go valida e calcula os valores; a interface apresenta o detalhamento e oferece geração do PDF.
- O PDF é aberto no leitor padrão e o teste garante uma página A4, inclusive com campos longos.
- Os dias trabalhados são exibidos separadamente e não integram o total das verbas.
- Regras e ressalvas estão em [docs/regras-calculo.md](docs/regras-calculo.md).
- Instruções para continuar o desenvolvimento estão em [docs/GUIA-PARA-RETOMADA.md](docs/GUIA-PARA-RETOMADA.md).

## Requisitos

- Go 1.25 ou compatível com `go.mod`.
- Node.js/npm.
- Wails v2.16.0.
- WebView2 no Windows.

Instale o CLI Wails com:

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0
```

## Desenvolvimento e validação

Execute na raiz do projeto:

```powershell
go test ./...
npm.cmd --prefix frontend run build
wails dev
```

## Build Windows

```powershell
wails build
```

A configuração Wails fica em `wails.json` na raiz e o executável vai para `build/bin/`. Se o app já estiver aberto, use outro nome para não interromper a sessão:

```powershell
wails build -o RescisaoDesk-atualizado.exe
```

O aplicativo funciona localmente. Os PDFs são gravados na pasta temporária do usuário, dentro de `RescisaoDesk`.
