# Skill: Wails

## Objetivo

Utilizar Wails como ponte entre a interface desktop e o backend Go.

## Arquitetura

``` text
HTML / CSS / JavaScript
          ↓
        Wails
          ↓
          Go
          ↓
   Regras de negócio
```

## Responsabilidade do frontend

O frontend deve cuidar principalmente de:

-   formulário;
-   validação básica de UX;
-   apresentação;
-   mensagens;
-   estado visual;
-   exibição dos resultados;
-   acionamento da geração do PDF.

## Responsabilidade do Go

Go deve cuidar de:

-   validação definitiva;
-   regras de rescisão;
-   cálculo;
-   valores monetários;
-   memória de cálculo;
-   geração do PDF;
-   tratamento de erros.

## Não duplicar regras

Não implemente uma fórmula financeira simultaneamente em Go e
JavaScript.

O frontend deve enviar os dados para Go e receber o resultado.

Evite:

``` text
JavaScript calcula férias
+
Go calcula férias
```

O correto é:

``` text
JavaScript → envia dados
Go → calcula
Go → retorna resultado
JavaScript → apresenta
```

## Interface com Go

Os métodos expostos pelo Wails devem ter nomes claros e
responsabilidades pequenas.

Exemplos conceituais:

``` text
CalcularRescisao()
GerarPDF()
```

A implementação interna pode utilizar serviços menores.

## Estado

Como o MVP não possui banco de dados, não criar persistência
desnecessária.

O estado da tela pode permanecer no frontend enquanto o usuário está
trabalhando no termo.

## UX

Após calcular:

1.  Mostrar o resultado.
2.  Mostrar a memória de cálculo.
3.  Permitir corrigir os dados.
4.  Permitir recalcular.
5.  Permitir gerar o PDF somente quando os dados estiverem válidos.

## Segurança

Não assumir que a validação do frontend é suficiente.

Toda entrada recebida pelo Go deve ser validada novamente.

## Build

O aplicativo deve funcionar como desktop Windows e não depender de
servidor externo para realizar o cálculo.
