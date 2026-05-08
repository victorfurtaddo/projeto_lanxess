package inference

import (
	"math"
	"sort"

	"projeto-lanxess/internal/model"
)

type Config struct {
	SpeedThreshold        float64 `json:"limite_velocidade"`
	WeightDropThreshold   float64 `json:"limite_queda_peso_kg"`
	PositionWindowDeg     float64 `json:"janela_posicao_graus"`
	RadiusWindowM         float64 `json:"janela_raio_m"`
	MaxAngularDistanceDeg float64 `json:"distancia_angular_max"`
	MaxRadialDistanceM    float64 `json:"distancia_radial_max"`
}

func DefaultConfig() Config {
	return Config{
		SpeedThreshold:        0.18,
		WeightDropThreshold:   120,
		PositionWindowDeg:     1.5,
		RadiusWindowM:         0.35,
		MaxAngularDistanceDeg: 8,
		MaxRadialDistanceM:    1.2,
	}
}

func DetectEvents(samples []model.CraneSample, reactors []model.Reactor, cfg Config) []model.InferredEvent {
	if cfg == (Config{}) {
		cfg = DefaultConfig()
	}

	events := make([]model.InferredEvent, 0)
	for i := 1; i < len(samples); i++ {
		prev := samples[i-1]
		curr := samples[i]
		drop := prev.WeightKg - curr.WeightKg
		stableAngle := AngularDistance(prev.AngleDeg, curr.AngleDeg) <= cfg.PositionWindowDeg
		stableRadius := math.Abs(prev.RadiusM-curr.RadiusM) <= cfg.RadiusWindowM
		stopped := prev.Speed <= cfg.SpeedThreshold && curr.Speed <= cfg.SpeedThreshold

		if drop < cfg.WeightDropThreshold || !stopped || !stableAngle || !stableRadius {
			continue
		}

		candidates := RankCandidates(curr, reactors, cfg)
		if len(candidates) == 0 {
			continue
		}

		best := candidates[0]
		events = append(events, model.InferredEvent{
			Timestamp:  curr.Timestamp,
			ReactorID:  best.ReactorID,
			AmountKg:   drop,
			Confidence: best.Probability,
			AngleDeg:   curr.AngleDeg,
			RadiusM:    curr.RadiusM,
			Reason:     "baixa velocidade, posicao estavel e queda de peso",
			Candidates: candidates,
		})
	}

	return events
}

func RankCandidates(sample model.CraneSample, reactors []model.Reactor, cfg Config) []model.Candidate {
	if cfg == (Config{}) {
		cfg = DefaultConfig()
	}

	candidates := make([]model.Candidate, 0, len(reactors))
	for _, reactor := range reactors {
		angularDistance := AngularDistance(sample.AngleDeg, reactor.AngleDeg)
		radialDistance := math.Abs(sample.RadiusM - reactor.RadiusM)
		if angularDistance > cfg.MaxAngularDistanceDeg || radialDistance > cfg.MaxRadialDistanceM {
			continue
		}

		angularScore := 1 - (angularDistance / cfg.MaxAngularDistanceDeg)
		radialScore := 1 - (radialDistance / cfg.MaxRadialDistanceM)
		score := math.Max(0, angularScore)*0.75 + math.Max(0, radialScore)*0.25
		candidates = append(candidates, model.Candidate{
			ReactorID:          reactor.ID,
			AngularDistanceDeg: angularDistance,
			RadialDistanceM:    radialDistance,
			Score:              score,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].ReactorID < candidates[j].ReactorID
		}
		return candidates[i].Score > candidates[j].Score
	})

	total := 0.0
	for _, candidate := range candidates {
		total += candidate.Score
	}
	if total > 0 {
		for i := range candidates {
			candidates[i].Probability = candidates[i].Score / total
		}
	}

	if len(candidates) > 3 {
		return candidates[:3]
	}
	return candidates
}

func EvaluateAccuracy(truth []model.TruthEvent, inferred []model.InferredEvent) float64 {
	if len(truth) == 0 || len(inferred) == 0 {
		return 0
	}

	limit := len(truth)
	if len(inferred) < limit {
		limit = len(inferred)
	}

	matches := 0
	for i := 0; i < limit; i++ {
		if truth[i].ReactorID == inferred[i].ReactorID {
			matches++
		}
	}

	return float64(matches) / float64(len(truth))
}
