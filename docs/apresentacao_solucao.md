# Apresentacao da Solucao

## Problema abordado

A area industrial possui uma grua central responsavel por abastecer 28 reatores com aparas de aco, mas opera com baixa instrumentacao.

Atualmente, nao ha medicao automatica da carga recebida por cada reator, nem monitoramento automatico da posicao da grua, da carga suspensa ou da presenca de pessoas em zonas de risco.

O principal desafio e reconstruir eventos fisicos reais a partir de dados indiretos, como angulo, raio, velocidade, carga estimada e timestamp, transformando esses sinais em informacoes operacionais uteis.

## Estrategia / abordagem utilizada

A abordagem adotada combina modelagem fisica, regras operacionais e simulacao.

O sistema foi estruturado como um mini gemeo digital industrial, representando virtualmente a grua, os 28 reatores, os ciclos de icamento, as cargas transportadas, as pessoas simuladas e os alertas de seguranca.

A inferencia de carregamento usa regras geometricas e logicas:

- deteccao de baixa velocidade da grua;
- identificacao de queda de carga;
- verificacao de estabilidade da posicao;
- comparacao entre posicao da grua e posicao dos reatores;
- calculo de confianca para o reator mais provavel.

Tambem foram incorporadas regras de seguranca:

- calculo de massa a partir do volume de aparas;
- aplicacao da densidade aparente de 500 kg/m3;
- limite da rede de icamento de 2,5 t;
- limite de capacidade da grua conforme o raio;
- deteccao de pessoas em zonas de risco.

## Solucao proposta

A solucao desenvolvida e um MVP funcional em Go com dashboard web.

Ela inclui:

- simulador da grua e dos ciclos de carga e descarga;
- criacao automatica dos 28 reatores em distribuicao radial;
- estimativa de massa transportada com base no volume e na densidade do material;
- aplicacao do limite operacional `min(capacidade_grua_por_raio, 2500 kg)`;
- inferencia automatica do reator carregado;
- acumulacao de carga e contagem de ciclos por reator;
- simulacao de pessoas na area operacional;
- geracao de alertas `INFO`, `WARNING` e `CRITICAL`;
- API HTTP para consulta do estado do gemeo digital;
- dashboard para visualizar grua, reatores, pessoas, alertas e historico.

Endpoints principais:

```text
GET /health
GET /api/state
GET /api/crane
GET /api/reactors
GET /api/events
GET /api/alerts
GET /api/demo
```

## Limitacoes identificadas

O MVP ainda utiliza dados simulados, pois nao ha integracao com sensores reais, cameras, CLPs ou sistemas industriais.

As principais limitacoes sao:

- ausencia de persistencia em banco de dados;
- ausencia de aquisicao em tempo real por sensores;
- pessoas simuladas por mock, sem visao computacional real;
- zonas de risco simplificadas por geometria 2D;
- status de tampas dos reatores ainda nao integrado a dados reais;
- dashboard atualizado por nova simulacao, sem WebSocket continuo;
- modelo de inferencia baseado em regras, ainda sem calibracao por historico operacional real.

Essas limitacoes foram mantidas de forma consciente para preservar clareza, rastreabilidade e foco no raciocinio do gemeo digital.

## Impacto potencial da solucao

A solucao tem potencial para aumentar a confiabilidade operacional e melhorar a seguranca da area.

Os principais impactos esperados sao:

- maior visibilidade sobre a distribuicao de carga nos reatores;
- reducao da dependencia de registros manuais;
- apoio a decisao do operador durante ciclos de icamento;
- deteccao antecipada de operacoes proximas ao limite de carga;
- alerta para presenca de pessoas sob carga suspensa;
- base tecnica para evolucao com cameras, sensores, CLPs e conectividade 5G;
- criacao de um historico operacional para analises futuras, previsao de carga e otimizacao de processo.

Em uma evolucao futura, o MVP pode ser conectado a dados reais para se tornar um gemeo digital operacional em tempo quase real, integrando monitoramento, seguranca e inteligencia de processo.
