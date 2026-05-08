# Premissas e Simplificacoes

- O MVP usa simulacao deterministica por requisicao, sem estado persistente em banco.
- Os endpoints `POST /api/simulation/*` retornam aceite, mas ainda nao controlam um loop persistente.
- Os 28 reatores sao distribuidos uniformemente em um circulo de raio 8 m.
- A capacidade por raio usa tres faixas configuradas em codigo: 3 m, 6 m e raio longo.
- A densidade aparente das aparas de aco e fixa em 500 kg/m3.
- A rede de icamento limita todos os ciclos a 2500 kg.
- Pessoas sao geradas por um detector mock em `internal/vision`.
- As zonas de risco sao geometricas e simplificadas: raio critico ao redor da carga e exclusao circular ao redor de reatores.
- O dashboard usa polling manual pelo botao Simular, sem WebSocket nesta versao.
