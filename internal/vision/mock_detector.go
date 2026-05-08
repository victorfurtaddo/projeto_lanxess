package vision

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"projeto-lanxess/internal/model"
)

type MockDetector struct {
	rng *rand.Rand
}

func NewMockDetector(seed int64) *MockDetector {
	return &MockDetector{rng: rand.New(rand.NewSource(seed))}
}

func (d *MockDetector) DetectPeople(sample model.CraneSample, step int) []model.PersonDetection {
	people := make([]model.PersonDetection, 0, 2)

	angle := float64((step*31)%360) * math.Pi / 180
	radius := 3.2 + math.Sin(float64(step)/3)*1.8
	people = append(people, model.PersonDetection{
		ID:         "person_01",
		X:          math.Cos(angle) * radius,
		Y:          math.Sin(angle) * radius,
		SpeedMps:   0.6 + d.rng.Float64()*0.4,
		State:      "walking",
		Confidence: 0.82 + d.rng.Float64()*0.12,
		Timestamp:  sample.Timestamp,
	})

	if step%9 == 4 && sample.SuspendedLoad {
		x, y := loadPosition(sample)
		people = append(people, model.PersonDetection{
			ID:         fmt.Sprintf("person_%02d", 2+step%5),
			X:          x + 0.35,
			Y:          y - 0.25,
			SpeedMps:   0.35,
			State:      "walking",
			Confidence: 0.91,
			Timestamp:  sample.Timestamp,
		})
	}

	return people
}

func loadPosition(sample model.CraneSample) (float64, float64) {
	radians := (sample.AngleDeg - 90) * math.Pi / 180
	return math.Cos(radians) * sample.RadiusM, math.Sin(radians) * sample.RadiusM
}

func MockDetectionAt(id string, x, y float64, timestamp time.Time) model.PersonDetection {
	return model.PersonDetection{
		ID:         id,
		X:          x,
		Y:          y,
		SpeedMps:   0.4,
		State:      "walking",
		Confidence: 0.9,
		Timestamp:  timestamp,
	}
}
