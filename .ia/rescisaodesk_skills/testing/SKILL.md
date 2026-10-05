# Skill: Testing

## Objetivo

Garantir que o RescisaoDesk produza resultados corretos e não sofra
regressões nas regras de cálculo.

## Regra principal

Toda regra de cálculo nova deve possuir testes.

Testes devem validar o resultado esperado e não apenas se a função
executa sem erro.

## Prioridade

A maior cobertura deve estar em:

1.  cálculos;
2.  contagem de dias;
3.  contagem de avos;
4.  valores monetários;
5.  arredondamentos;
6.  validação de datas.

## Casos de borda

Para regras baseadas em dias, testar pelo menos:

``` text
limite - 1
limite
limite + 1
```

Exemplo:

``` text
14 dias
15 dias
16 dias
```

Quando a regra utilizar 30 dias:

``` text
29 dias
30 dias
31 dias
```

## Datas

Testar:

-   mesmo mês;
-   meses diferentes;
-   virada de ano;
-   meses com 28 dias;
-   meses com 29 dias;
-   meses com 30 dias;
-   meses com 31 dias;
-   datas iguais, quando permitido;
-   desligamento anterior à admissão;
-   períodos completos;
-   períodos incompletos.

## Valores

Testar:

-   salário inteiro;
-   salário com centavos;
-   salário baixo;
-   salário alto;
-   resultado exato;
-   resultado que exige arredondamento.

## Regras de negócio

Cada exemplo definido em `docs/regras-calculo.md` deve virar um teste automatizado.

Exemplo conceitual:

``` go
func TestFeriasAvos(t *testing.T) {
    tests := []struct {
        name string
        admissao time.Time
        desligamento time.Time
        want int
    }{
        // exemplos oficiais do projeto
    }

    // execução
}
```

## Não inventar expected values

Se o valor esperado não estiver definido ou não puder ser calculado de
forma inequívoca pela documentação, não invente.

Peça confirmação da regra.

## Teste de integração

Depois de estabilizar os cálculos, testar o fluxo:

``` text
formulário
→ Go
→ cálculo
→ resultado
→ geração do PDF
```

## Antes de finalizar uma tarefa

Executar os testes relevantes.

Se possível, executar:

``` bash
go test ./...
```

Corrigir regressões antes de considerar a tarefa concluída.

## Qualidade

Um teste deve explicar qual regra está protegendo.

Prefira:

``` text
TestFeriasAvos_Com45Dias_DeveTerDoisAvos
```

a nomes genéricos como:

``` text
TestCalculo1
```
