# Guia de Retomada — RescisaoDesk

Este documento permite que outra IA ou pessoa continue o projeto sem reabrir decisões já tomadas nem inventar regras municipais. A fonte canônica de cálculo é [`regras-calculo.md`](regras-calculo.md).

## Objetivo e limites

Aplicativo desktop local para preencher dados de rescisão, conferir verbas e gerar um termo PDF. Stack: Go 1.25, Wails v2.16, JavaScript/CSS sem framework e Vite. A lógica financeira pertence a `internal/rescisao`; o PDF recebe valores já calculados e nunca deve recalcular.

O app não persiste dados dos servidores nem depende de serviço remoto. O PDF é salvo em `os.TempDir()/RescisaoDesk` e aberto no leitor padrão.

## Estado do projeto

- Wizard de quatro etapas: servidor, dados da rescisão, dados bancários e conferência.
- A interface chama `ValidarDados` e `CalcularResultado`, mostra valores reais e oferece `Gerar PDF`.
- O backend gera o PDF e tenta abri-lo automaticamente; o teste exige exatamente uma página A4, inclusive com campos longos.
- Valores monetários são armazenados em centavos (`Money`, `int64`) e arredondados por verba.
- O total inclui 13º, férias e 1/3. O valor dos dias trabalhados fica separado e não compõe o total.
- Executáveis, `frontend/dist` e `node_modules` são ignorados pelo Git; publique o código-fonte, não os binários gerados.

## Arquivos principais

- [`app.go`](../app.go): métodos Wails para validação, cálculo e geração/abertura do PDF.
- [`internal/rescisao/entrada.go`](../internal/rescisao/entrada.go): entrada, validação, moeda e datas.
- [`internal/rescisao/calculation.go`](../internal/rescisao/calculation.go): avos, fórmulas monetárias, arredondamento e total.
- [`internal/rescisao/entrada_test.go`](../internal/rescisao/entrada_test.go) e [`internal/rescisao/calculation_test.go`](../internal/rescisao/calculation_test.go): testes de validação e cálculo.
- [`internal/pdf/pdf.go`](../internal/pdf/pdf.go): layout e geração de PDF A4 em uma página.
- [`internal/pdf/pdf_test.go`](../internal/pdf/pdf_test.go): validação do arquivo e contagem de páginas.
- [`frontend/src/main.js`](../frontend/src/main.js): wizard, apresentação dos valores e chamadas Wails.
- [`wails.json`](../wails.json): configuração Wails única, na raiz do projeto.
- [`docs/regras-calculo.md`](regras-calculo.md): regras canônicas e ressalvas.
- [`.ia/rescisaodesk_skills/pdf/SKILL.md`](../.ia/rescisaodesk_skills/pdf/SKILL.md): orientação visual para o documento.

## Regras atualmente codificadas

- **13º:** um avo ao completar 15 dias em cada período mensal ancorado na admissão. Acumula durante todo o vínculo, sem reinício anual ou teto fixo. Subtrai os avos recebidos; saldo devido mínimo zero. Exemplo testado: 18 adquiridos − 6 recebidos = 12 devidos. Fórmula: `salário-base × avos devidos ÷ 12`.
- **Férias:** implementação atual usa dias inclusivos, um avo a cada 30 dias e um adicional quando o resto é de pelo menos 15. Subtrai avos pagos; rejeita pagos acima do total disponível. Ainda não registra saldos por período aquisitivo.
- **Dias trabalhados:** conta do dia 1 até o desligamento, inclusive. Dia 31 ou último dia do mês é tratado como 30. O valor usa salário-base × dias ÷ 30 e fica fora do total.
- **Base provisória:** salário-base é usado nas verbas; a composição legal das bases ainda precisa de confirmação da área responsável.
- **1/3:** um terço das férias proporcionais já arredondadas.

## Pendências conhecidas

- O abatimento de férias é agregado e não separa períodos aquisitivos, férias gozadas e pagamentos por período.
- Confirmar se gratificações/adicionais devem compor as bases de férias e dias; hoje o sistema assume salário-base.
- Rever casos de admissão nos dias 29–31 e anos bissextos para avos de férias, com exemplos confirmados.
- Validar o modelo do termo com RH e incluir logotipo oficial apenas quando o arquivo for fornecido.

## Próximos passos sugeridos

1. Confirmar com RH a composição das bases de férias e dos dias trabalhados; registrar a regra adotada em `regras-calculo.md` e nos testes.
2. Definir como férias gozadas/pagas se relacionam aos períodos aquisitivos e então modelar os saldos individualmente, se necessário.
3. Adicionar casos aprovados para datas de admissão 29–31, fevereiro e anos bissextos antes de alterar a contagem de férias.
4. Validar o modelo final do PDF com RH; adicionar logotipo somente se fornecido.
5. Antes de entregar alterações, executar os testes Go, build do frontend e build Wails; confirmar que o PDF segue em uma página.

## Comandos úteis

```powershell
go test ./...
npm.cmd --prefix frontend run build
wails dev
wails build
```

Execute Wails a partir da raiz, onde está `wails.json`. Para não sobrescrever um executável aberto, use `wails build -o RescisaoDesk-atualizado.exe`; a saída fica em `build/bin/`.
