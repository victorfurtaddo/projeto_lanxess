package simulator

import (
	"testing"

	"projeto-lanxess/internal/inference"
)

func TestGenerateCanBeInferred(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Events = 8
	sim := Generate(cfg)
	events := inference.DetectEvents(sim.Samples, sim.Reactors, inference.DefaultConfig())

	if len(events) != cfg.Events {
		t.Fatalf("inferred %d events, want %d", len(events), cfg.Events)
	}

	accuracy := inference.EvaluateAccuracy(sim.Truth, events)
	if accuracy < 0.85 {
		t.Fatalf("accuracy = %.2f, want at least 0.85", accuracy)
	}
}
