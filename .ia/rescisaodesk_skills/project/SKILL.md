# Skill: Project

## Objetivo

Você está trabalhando no projeto **RescisaoDesk**, um aplicativo desktop
Windows para cálculo de rescisão de servidores municipais e geração de
um Termo de Rescisão em PDF.

O objetivo é criar uma aplicação simples, confiável e fácil de usar por
servidores do setor de Recursos Humanos.

## Stack

-   Go para backend e regras de negócio.
-   Wails para integração entre Go e a interface desktop.
-   HTML, CSS e JavaScript para frontend.
-   PDF gerado localmente pelo aplicativo.
-   Aplicação offline.
-   Executável Windows (`.exe`).

## Escopo do MVP

O MVP deve:

1.  Permitir informar manualmente os dados do servidor.
2.  Calcular os valores da rescisão conforme
    `docs/regras-calculo.md`.
3.  Exibir a memória de cálculo.
4.  Permitir revisar os resultados.
5.  Gerar o Termo de Rescisão em PDF.

O MVP NÃO deve possuir:

-   banco de dados;
-   SQLite;
-   importação de CSV;
-   servidor web;
-   API externa;
-   login;
-   autenticação;
-   internet;
-   cadastro permanente de servidores.

Os dados são informados para cada cálculo e não precisam ser
armazenados.

## Arquitetura

Mantenha uma separação clara:

``` text
Frontend
    ↓
Wails
    ↓
Go
    ↓
Domínio / Cálculos
    ↓
PDF
```

O frontend apresenta dados e resultados. As regras de negócio pertencem
ao Go.

## Regras gerais

-   Não invente regras trabalhistas ou municipais.
-   Consulte `docs/regras-calculo.md` antes de implementar qualquer cálculo.
-   Se uma regra estiver ambígua, pare e peça confirmação em vez de
    assumir.
-   Não coloque regras financeiras importantes no JavaScript.
-   Não altere arquivos não relacionados à tarefa.
-   Evite dependências externas desnecessárias.
-   Prefira código simples e explícito.
-   Toda alteração relevante deve ser acompanhada de testes quando
    aplicável.
-   Não implemente funcionalidades fora do escopo sem solicitação.

## Fonte das regras

`docs/regras-calculo.md` é a fonte de verdade para as regras de cálculo.

As skills explicam COMO implementar. A documentação de regras explica
QUAL cálculo deve ser feito.

## Prioridade

Em caso de conflito:

1.  Regras explicitamente confirmadas pelo usuário.
2.  `docs/regras-calculo.md`.
3.  Requisitos documentados.
4.  Arquitetura do projeto.
5.  Convenções gerais de programação.

Nunca substitua uma regra municipal documentada por uma regra genérica
da CLT sem autorização.
