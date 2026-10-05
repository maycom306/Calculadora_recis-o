# Skill: Calculations

## Objetivo

Implementar os cálculos da rescisão de forma determinística, testável e
auditável.

## Regra principal

**Nunca invente uma regra de cálculo.**

Antes de implementar qualquer fórmula:

1.  Leia `docs/regras-calculo.md`.
2.  Identifique exatamente a regra aplicável.
3.  Verifique os exemplos documentados.
4.  Transforme a regra em código.
5.  Crie testes para os casos documentados.

Se houver dúvida sobre a interpretação, não faça uma suposição
silenciosa.

## Itens do MVP

A rescisão contempla somente os itens definidos na documentação municipal
do projeto:

-   13º proporcional;
-   férias proporcionais;
-   1/3 constitucional/proporcional de férias.

Dias trabalhados devem ser discriminados como informação e não entram no
total monetário.

Não adicionar automaticamente:

-   aviso prévio;
-   FGTS;
-   multa de 40%;
-   INSS;
-   IRRF;
-   outras verbas.

Somente implementar outros itens quando forem explicitamente adicionados
às regras do projeto.

## Avos

A contagem de avos deve seguir exatamente a regra documentada em
`docs/regras-calculo.md`.

Não utilizar automaticamente uma regra genérica encontrada na internet.

Especial atenção para:

-   data de admissão;
-   data de exoneração;
-   períodos completos;
-   contagem de avos do 13º proporcional, considerando 15 dias por avo conforme informado;
-   regra específica de 15 dias;
-   regra específica de 30 dias;
-   eventual período residual;
-   contagem informativa de dias trabalhados, sem inclusão no total monetário.

## Dinheiro

Não utilize `float64` para representar valores monetários quando isso
puder gerar perda de precisão.

Prefira uma representação determinística, como:

-   centavos em `int64`; ou
-   uma biblioteca decimal adequada, se realmente necessária.

Toda operação monetária deve possuir regra clara de arredondamento.

## Estrutura recomendada

Uma função de cálculo deve ser pequena e previsível.

Exemplo conceitual:

``` go
type ResultadoRescisao struct {
    DiasTrabalhadosInfo int
    DecimoTerceiroAvos  int
    ValorDecimoTerceiro int64
    FeriasAvos          int
    ValorFerias         int64
    ValorTercoFerias    int64
    Total               int64
}
```

Os nomes podem ser adaptados à arquitetura real.

## Memória de cálculo

O sistema deve conseguir explicar o resultado.

Para cada verba, mantenha informações suficientes para mostrar:

-   quantidade de dias ou avos;
-   base utilizada;
-   fórmula;
-   resultado;
-   arredondamento, quando aplicável.

O usuário deve conseguir conferir como o total foi obtido.

## Testes

Toda regra de cálculo nova deve possuir testes.

Inclua, quando aplicável:

-   caso normal;
-   zero dias;
-   salário inteiro;
-   salário com centavos;
-   datas no mesmo mês;
-   mudança de mês;
-   períodos completos;
-   limite da regra;
-   dia imediatamente anterior ao limite;
-   dia exatamente no limite;
-   dia imediatamente posterior ao limite;
-   férias na exoneração;
-   13º proporcional;
-   contagem informativa de dias trabalhados, sem valor monetário.

## Proibição

Não esconda regra de negócio em:

-   handlers HTTP;
-   métodos Wails;
-   componentes JavaScript;
-   templates HTML;
-   código de geração do PDF.

Os cálculos devem ficar em uma camada própria.
