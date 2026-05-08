# Arquitetura do MVP

## Fluxo

```text
Simulador
  -> amostras da grua
  -> pessoas simuladas
  -> calculo de massa por volume
  -> limites da grua e da rede
  -> deteccao de queda de peso
  -> ranking de reatores candidatos
  -> regras de seguranca
  -> evento inferido
  -> alertas
  -> estado acumulado do gemeo digital
```

## Modelagem dos Reatores

Os 28 reatores sao distribuidos em um circulo. Cada reator possui:

- identificador;
- angulo;
- raio;
- capacidade;
- carga acumulada.

Esta modelagem permite transformar uma posicao continua da grua em uma decisao discreta de destino.

## Motor de Inferencia

O motor atual e propositalmente simples e explicavel:

1. Compara duas amostras consecutivas da grua.
2. Verifica se houve queda relevante de peso.
3. Confirma se a grua estava praticamente parada.
4. Confirma se a posicao estava estavel.
5. Calcula candidatos por distancia angular e radial.
6. Normaliza os scores para obter uma confianca.

## Massa e Limites

A massa operacional e calculada com a densidade aparente das aparas de aco:

```text
massa_kg = volume_m3 * 500
```

Depois sao aplicados os limites:

```text
limite_operacional = min(capacidade_grua_por_raio, 2500 kg)
massa_final = min(massa_kg, limite_operacional)
```

No MVP, os reatores ficam em raio longo, entao a rede de icamento e a ponta da grua convergem para o mesmo limite de 2500 kg.

## Seguranca

O modulo `internal/safety` avalia:

- carga proxima ao limite;
- carga acima do limite;
- pessoa sob carga suspensa;
- pessoa em zona de exclusao de reator;
- movimentacao da grua com pessoa em area critica.

As deteccoes de pessoas ainda sao simuladas pelo modulo `internal/vision`, mas os tipos ja permitem substituir o mock por uma fonte real de visao computacional.

## Caminho de Evolucao

- incorporar filtros temporais com janelas moveis;
- adicionar calibracao por dados reais;
- persistir series temporais em PostgreSQL ou TimescaleDB;
- incluir sensores auxiliares quando existirem;
- trocar o simulador por stream real;
- adicionar alertas de capacidade e anomalias.
