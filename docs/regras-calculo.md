# Regras de Cálculo — RescisaoDesk

Este documento é a fonte oficial das regras de cálculo do RescisaoDesk. Não substituir regras municipais por fórmulas genéricas da CLT, de legislação federal ou de outros municípios.

## Estado das Regras

### Definidas

- O 13º usa exclusivamente salário-base/vencimento-base, pela fórmula salário-base × avos ÷ 12. Gratificações e demais parcelas ficam fora.
- No 13º, cada período mensal ancorado na data de admissão gera 1 avo ao completar 15 dias; a contagem acumula sem reiniciar no aniversário anual.
- Férias na exoneração usam meses completos de efetivo exercício desde a admissão e fração residual de 15 dias ou mais para um avo adicional.
- Se o dia da admissão não existir no mês de destino, o marco é o último dia daquele mês.
- Os dias trabalhados são contados do dia 1 até o dia do desligamento no último mês, incluindo a data de desligamento.
- O 1/3 corresponde a um terço do valor das férias proporcionais. Valores monetários são arredondados por verba antes do total.
- O valor dos dias do último mês usa divisor fixo 30 e é exibido separadamente, fora do total das verbas.
- Férias gozadas/pagas devem ser informadas e abatidas dos avos disponíveis; o saldo pode ultrapassar 12 avos antes do abatimento.

### Pendentes

- Como controlar períodos aquisitivos de férias, férias já gozadas/pagas e o saldo ainda devido ao longo de vários anos; não limitar automaticamente a 12 avos.
- Quais parcelas compõem a remuneração do mês da exoneração usada na base das férias.
- Confirmar com a área responsável se as bases atualmente assumidas (salário-base e divisor 30) correspondem às regras locais.

## Escopo

O cálculo contém:

- 13º proporcional;
- férias proporcionais;
- 1/3 de férias;
- dedução monetária dos dias trabalhados no mês do desligamento.

O total é:

```text
13º proporcional + férias proporcionais + 1/3 de férias
```

O valor dos dias trabalhados é exibido à parte e não compõe esse total.

Não calcular aviso prévio, FGTS, multa de 40%, INSS, IRRF ou outras verbas sem regra expressamente aprovada.

## 13º Proporcional

### Base

Usar exclusivamente o salário-base/vencimento-base do servidor. Não incluir gratificações, adicionais, funções gratificadas, vantagens pessoais, verbas eventuais ou quaisquer outras parcelas. O salário-base deve permanecer separado das demais verbas.

### Fórmula

```text
13º proporcional = salário-base × avos ÷ 12
```

Exemplo:

```text
Salário-base: R$ 3.000,00
Avos: 8
13º proporcional: R$ 3.000,00 × 8 ÷ 12 = R$ 2.000,00
```

### Avos

Seguir a regra confirmada pelo responsável: a contagem começa na admissão, usa períodos mensais ancorados nessa data e concede no máximo 1 avo por período ao completar 15 dias de efetivo exercício naquele período. A contagem acumula ao longo de todo o vínculo, sem reiniciar no aniversário anual nem teto automático de 12 avos. Os avos adquiridos já recebidos devem ser descontados; o saldo devido não pode ser negativo.

Exemplo para admissão em 01/01:

```text
01/01: início do período, 0 avo
15/01: 1 avo
30/01: continua em 1 avo
15/02: 2 avos
15/03: 3 avos
15/01 do ano seguinte: 13 avos
```

Não aplicar uma sequência contínua de blocos de 15 dias que conceda 2 avos entre 01/01 e 30/01; em cada novo período mensal ancorado na admissão, contar 15 dias inclusivamente a partir do marco ajustado daquele período. Uma data como 30/02 não existe no calendário; em fevereiro, o marco mensal usa o último dia do mês quando necessário.

## Férias Proporcionais na Exoneração

### Avos

- cada período de 30 dias corridos de efetivo exercício: 1 avo;
- resto de 14 dias: não gera avo adicional;
- resto de 15 dias ou mais: gera 1 avo adicional.

Os avos disponíveis podem abranger mais de um período aquisitivo e não têm teto automático de 12. A entrada deverá registrar os avos já gozados/pagos (zero quando nenhum); saldo devido = avos disponíveis − avos já gozados/pagos. Não contar novamente avos quitados.

Implementação atual: conta dias inclusivamente entre admissão e desligamento, divide por 30 e aplica a regra do resto acima. A regra ainda não separa saldos por período aquisitivo.

Os avos de férias são acumulados desde a admissão e não podem ser limitados simplesmente a 12 por ano. Para a exoneração, considerar somente os avos ainda devidos nos períodos aquisitivos, descontando férias já gozadas ou pagas. A regra exata para identificar períodos quitados e o saldo remanescente ainda precisa ser definida; esses dados devem ser informados separadamente e não podem ser presumidos.

### Base e fórmula

A base de referência é a remuneração do mês da exoneração:

```text
férias proporcionais = remuneração do mês da exoneração × avos ÷ 12
```

A composição dessa remuneração ainda depende de confirmação no Estatuto municipal. Não presumir quais gratificações, adicionais ou vantagens entram na base.

## 1/3 de Férias

```text
1/3 de férias = férias proporcionais ÷ 3
```

## Dias do Mês do Desligamento

Contar do dia 1 até o dia do desligamento, incluindo esse dia. Exemplo: desligamento no dia 3 corresponde a 3 dias. Meses anteriores não entram nessa contagem, pois já foram pagos.

O valor dos dias trabalhados é calculado como verba positiva, mas é exibido separadamente e não entra no total:

```text
valor dos dias trabalhados = base mensal × dias considerados ÷ 30
```

O sistema usa divisor fixo 30 e o salário-base informado como base. A área responsável ainda deve confirmar se essa base corresponde à regra local.

## Dinheiro e Arredondamento

Representar valores em centavos inteiros. Arredondar cada verba monetária para centavos depois do cálculo, usando arredondamento para o centavo mais próximo e meio centavo para cima. O total soma 13º proporcional, férias proporcionais e 1/3 arredondados; os dias trabalhados aparecem em separado.

Exibir separadamente a quantidade de dias e o valor da dedução. Não misturar essa parcela com os avos ou valores de 13º/férias.

## Casos de Teste Obrigatórios

- 13º: conferir os limites de 14, 15 e 16 dias, a permanência de apenas 1 avo no mesmo período mensal e o acúmulo após o aniversário anual da admissão.
- Férias na exoneração: 29, 30 e 31 dias completos; resíduos de 14, 15 e 16 dias; férias já gozadas/pagas e saldo de múltiplos períodos aquisitivos.
- Marco de admissão nos dias 1, 10 e 31; meses com 28, 29, 30 e 31 dias; virada de mês e de ano.
- Dias do mês do desligamento: desligamento nos dias 1 e 3; conferir contagem e dedução em linha separada. Não calcular o valor enquanto base e divisor estiverem pendentes.
- Valores: salário-base de R$ 3.000,00 com 8 avos deve resultar em R$ 2.000,00 de 13º antes do arredondamento.

## Pendências

Pontos de negócio ainda sujeitos a confirmação, sem bloquear o fluxo técnico atual:

- como representar férias gozadas/pagas e calcular o saldo ainda devido por período aquisitivo;
- quais parcelas, se houver, compõem a remuneração do mês da exoneração usada como base das férias;
- quais parcelas compõem a base monetária da dedução dos dias trabalhados;
- qual divisor aplicar à dedução dos dias trabalhados;
- demais particularidades aplicáveis do Estatuto ou regulamentação de Ipanguaçu.

Atualizar este documento e os testes sempre que uma regra for confirmada. Uma referência à Lei 8.112/1990 ou a qualquer legislação externa é apenas comparativa e não confirma, por si só, a regra municipal.