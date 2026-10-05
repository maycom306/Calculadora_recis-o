# Skill: PDF

## Objetivo

Gerar um Termo de Rescisão profissional, legível e adequado para
impressão e anexação a processos administrativos.

## Formato

O documento deve ser:

-   A4;
-   preferencialmente uma página quando o conteúdo permitir;
-   legível em impressão;
-   organizado;
-   com margens adequadas;
-   sem elementos decorativos desnecessários.

### Cabeçalho institucional

Usar este cabeçalho no topo do A4:

```text
ESTADO DO RIO GRANDE DO NORTE
PREFEITURA MUNICIPAL DE IPANGUAÇU
SECRETARIA MUNICIPAL DE PLANEJAMENTO, ADMINISTRAÇÃO E TRANSPORTES
DEPARTAMENTO DE RECURSOS HUMANOS
```

Seguir a identidade visual do material enviado: azul em tons institucionais e vermelho como acento discreto. Incorporar o logotipo oficial à esquerda quando o arquivo de imagem estiver disponível; não redesenhar ou inventar o brasão.

### Estrutura visual

Usar A4 com estrutura semelhante ao termo de referência: identificação do órgão e do servidor, dados funcionais, discriminação das verbas, totais e áreas de assinatura. Adaptar campos ao RescisaoDesk sem copiar campos legais ou dados não confirmados.


## Conteúdo

O PDF deve conter, conforme os requisitos do projeto:

### Identificação

-   órgão/município;
-   título do documento;
-   nome do servidor;
-   matrícula, quando informada;
-   cargo, quando informado;
-   tipo de desligamento, quando aplicável.

### Dados funcionais

-   data de admissão;
-   data de desligamento;
-   salário/base informado.

### Cálculo

Apresentar uma tabela clara com:

-   verba;
-   quantidade de dias/avos;
-   base de cálculo;
-   valor.

Itens do MVP:
-   13º proporcional;
-   férias proporcionais;
-   1/3 de férias.

Exibir em linhas separadas: dias considerados, dedução monetária dos dias, 13º proporcional, férias proporcionais e 1/3. A quantidade de dias não deve ser confundida com seu valor monetário.

### Total

Apresentar claramente o valor total calculado.

### Memória de cálculo

Sempre que possível, demonstrar as fórmulas utilizadas.

Exemplo conceitual:

``` text
Dias trabalhados (informativo):
XX dias no mês do desligamento.

Dedução dos dias:
Apresentar o valor calculado pelo domínio em linha própria.

13º proporcional:
Base × avos ÷ 12; apresentar a remuneração de referência e os avos.

Férias:
Férias proporcionais:
Remuneração do mês da exoneração × avos ÷ 12; não inferir quais parcelas compõem a remuneração.

1/3 de férias:
Valor das férias ÷ 3 = R$ ...
```

As fórmulas reais devem seguir
`docs/regras-calculo.md`.

## Formatação

Usar:

-   datas no padrão brasileiro;
-   valores em reais;
-   duas casas decimais;
-   números alinhados de forma consistente;
-   títulos e seções claramente diferenciados.

## Assinatura

Reservar espaço para assinatura quando exigido pelo modelo do documento.

## Nome do arquivo

Preferir um nome seguro e previsível, por exemplo:

``` text
MATRICULA_NOME_TERMO_RESCISAO.pdf
```

Se não houver matrícula, utilizar uma alternativa baseada no nome.

Remover caracteres inválidos para nomes de arquivos do Windows.

## Regras importantes

O gerador de PDF não deve recalcular valores.

Ele deve receber um resultado já calculado pelo domínio.

Fluxo:

``` text
Dados
 ↓
Cálculo
 ↓
Resultado
 ↓
PDF
```

Não:

``` text
Dados
 ↓
PDF calcula novamente
```

Isso evita divergência entre a tela e o documento.
