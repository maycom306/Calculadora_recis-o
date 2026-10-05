# Skill: UI

## Objetivo

Criar uma interface simples para uso cotidiano por servidores do setor
de Recursos Humanos.

## Princípios

A interface deve priorizar:

-   clareza;
-   poucos campos;
-   leitura fácil;
-   fluxo previsível;
-   baixo número de cliques;
-   mensagens claras;
-   aparência profissional.

Não transformar o MVP em um sistema administrativo complexo.

## Fluxo

O fluxo principal deve ser:

``` text
Novo Termo
    ↓
Dados do servidor
    ↓
Dados da rescisão
    ↓
Conferência do cálculo
    ↓
Calcular
    ↓
Conferir resultado
    ↓
Gerar PDF
```

## Formulário

Campos podem incluir:

-   nome;
-   matrícula;
-   cargo;
-   data de admissão;
-   data de desligamento;
-   salário.

Somente incluir campos realmente necessários.

## Resultado

Depois do cálculo, mostrar:

``` text
Dias trabalhados (informativo)  XX dias
13º proporcional                R$ ...
Férias proporcionais            R$ ...
1/3 de férias                   R$ ...
-----------------------------------------
TOTAL                           R$ ...
```

Também disponibilizar a memória de cálculo.

## Erros

Evitar mensagens técnicas como:

``` text
panic: runtime error
```

Mostrar mensagens úteis:

``` text
A data de desligamento deve ser posterior à data de admissão.
```

## Responsividade

Mesmo sendo desktop, a interface deve se adaptar a diferentes tamanhos
de janela.

## JavaScript

Não colocar as regras financeiras no frontend.

JavaScript pode realizar:

-   máscara de moeda;
-   formatação;
-   validação visual;
-   controle de campos;
-   comunicação com Wails.

O cálculo oficial deve ser feito no Go.

## Visual

Priorizar:

-   tipografia legível;
-   bom espaçamento;
-   contraste adequado;
-   botões claros;
-   estados de erro/sucesso;
-   aparência institucional.

Evitar excesso de animações, gradientes, cards e elementos decorativos.
