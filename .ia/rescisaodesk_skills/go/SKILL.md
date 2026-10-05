# Skill: Go

## Objetivo

Escrever o backend do RescisaoDesk utilizando Go idiomático, simples e
testável.

## Princípios

-   Prefira simplicidade.
-   Evite abstrações prematuras.
-   Use nomes claros.
-   Funções devem ter responsabilidade bem definida.
-   Erros devem ser explícitos.
-   Evite estado global.
-   Evite dependências sem necessidade.
-   Não misture apresentação com regra de negócio.

## Organização

A estrutura pode seguir:

``` text
internal/
├── domain/
├── calculation/
└── pdf/
```

Adapte a estrutura somente quando houver uma razão clara.

## Domínio

Os modelos devem representar conceitos do sistema.

Exemplos:

-   Servidor;
-   DadosRescisao;
-   ResultadoRescisao;
-   VerbaRescisoria;
-   MemoriaCalculo.

Não criar modelos desnecessários.

## Validação

Valide os dados antes do cálculo.

Exemplos:

-   nome obrigatório;
-   salário válido;
-   data de admissão válida;
-   data de desligamento válida;
-   desligamento não pode ocorrer antes da admissão;
-   campos obrigatórios preenchidos.

Mensagens de erro devem ser compreensíveis para o usuário.

## Valores monetários

Evite `float64` para dinheiro.

Prefira valores inteiros em centavos ou uma solução decimal apropriada.

Exemplo conceitual:

``` go
type Money int64
```

Se for utilizada outra abordagem, mantenha precisão e arredondamento
determinísticos.

## Wails

Os métodos expostos ao frontend devem funcionar como uma camada de
entrada.

Eles não devem conter a lógica completa de cálculo.

Fluxo esperado:

``` text
Wails method
    ↓
validação
    ↓
serviço/domínio
    ↓
resultado
```

## Erros

Não ignore erros.

Evite:

``` go
result, _ := algumaFuncao()
```

Prefira:

``` go
result, err := algumaFuncao()
if err != nil {
    return err
}
```

## Testes

Utilize os testes nativos do Go.

Prefira testes de tabela para várias combinações de entrada.

Exemplo:

``` go
tests := []struct {
    name string
    input Input
    want Output
}{
    // casos
}
```

## Antes de alterar código

1.  Leia o código existente.
2.  Identifique o fluxo atual.
3.  Faça a menor alteração necessária.
4.  Rode os testes.
5.  Verifique se não quebrou outra parte.

Não reestruture o projeto inteiro para resolver uma alteração pequena.
