package safety

import (
	"testing"
	"time"

	"projeto-lanxess/internal/model"
	"projeto-lanxess/internal/vision"
)

func TestEvaluateSampleDetectsPersonUnderSuspendedLoad(t *testing.T) {
	sample := model.CraneSample{
		Timestamp:     time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
		AngleDeg:      90,
		RadiusM:       4,
		WeightKg:      1200,
		SuspendedLoad: true,
		People: []model.PersonDetection{
			vision.MockDetectionAt("person_01", 4, 0, time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)),
		},
	}

	alerts := EvaluateSample(sample, nil)
	if len(alerts) == 0 {
		t.Fatal("EvaluateSample() returned no alerts")
	}
	if alerts[0].Level != "CRITICAL" {
		t.Fatalf("alert level = %s, want CRITICAL", alerts[0].Level)
	}
}
