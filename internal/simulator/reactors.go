package simulator

import (
	"fmt"
	"math/rand"

	"projeto-lanxess/internal/model"
)

func BuildReactors(count int, radiusM float64, seed int64) []model.Reactor {
	rng := rand.New(rand.NewSource(seed))
	reactors := make([]model.Reactor, 0, count)
	step := 360.0 / float64(count)

	for i := 0; i < count; i++ {
		reactors = append(reactors, model.Reactor{
			ID:         i + 1,
			Name:       fmt.Sprintf("Reator %02d", i+1),
			AngleDeg:   step * float64(i),
			RadiusM:    radiusM,
			CapacityKg: 4500 + rng.Float64()*2500,
			Status:     "normal",
			LidStatus:  "unknown",
		})
	}

	return reactors
}

func ApplyLoads(reactors []model.Reactor, events []model.InferredEvent) []model.Reactor {
	out := make([]model.Reactor, len(reactors))
	copy(out, reactors)

	index := map[int]int{}
	for i, reactor := range out {
		index[reactor.ID] = i
	}

	for _, event := range events {
		i, ok := index[event.ReactorID]
		if !ok {
			continue
		}
		out[i].LoadKg += event.AmountKg
		out[i].CycleCount++
		out[i].LastLoadedAt = &event.Timestamp
		out[i].LastLoadMassKg = event.AmountKg
		if out[i].LoadKg >= out[i].CapacityKg*0.9 {
			out[i].Status = "attention"
		}
	}

	return out
}
