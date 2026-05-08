package inference

import (
	"testing"
	"time"

	"projeto-lanxess/internal/model"
)

func TestDetectEventsFindsStableWeightDrop(t *testing.T) {
	reactors := []model.Reactor{
		{ID: 1, AngleDeg: 0, RadiusM: 8},
		{ID: 2, AngleDeg: 90, RadiusM: 8},
	}
	samples := []model.CraneSample{
		{
			Timestamp: time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			AngleDeg:  89.8,
			RadiusM:   8.05,
			Speed:     0.04,
			WeightKg:  1000,
		},
		{
			Timestamp: time.Date(2026, 5, 8, 10, 0, 20, 0, time.UTC),
			AngleDeg:  90.1,
			RadiusM:   8.02,
			Speed:     0.03,
			WeightKg:  120,
		},
	}

	events := DetectEvents(samples, reactors, DefaultConfig())
	if len(events) != 1 {
		t.Fatalf("DetectEvents() returned %d events, want 1", len(events))
	}
	if events[0].ReactorID != 2 {
		t.Fatalf("event ReactorID = %d, want 2", events[0].ReactorID)
	}
	if events[0].AmountKg != 880 {
		t.Fatalf("event AmountKg = %v, want 880", events[0].AmountKg)
	}
}

func TestDetectEventsIgnoresMovingWeightDrop(t *testing.T) {
	reactors := []model.Reactor{{ID: 1, AngleDeg: 0, RadiusM: 8}}
	samples := []model.CraneSample{
		{Timestamp: time.Now(), AngleDeg: 0, RadiusM: 8, Speed: 1.2, WeightKg: 1000},
		{Timestamp: time.Now().Add(time.Second), AngleDeg: 2, RadiusM: 8.2, Speed: 1.1, WeightKg: 100},
	}

	events := DetectEvents(samples, reactors, DefaultConfig())
	if len(events) != 0 {
		t.Fatalf("DetectEvents() returned %d events, want 0", len(events))
	}
}
