package safety

import (
	"math"

	"projeto-lanxess/internal/model"
)

const (
	CriticalLoadRadiusM = 1.35
	ReactorExclusionM   = 0.85
)

func LoadPosition(sample model.CraneSample) (float64, float64) {
	radians := (sample.AngleDeg - 90) * math.Pi / 180
	return math.Cos(radians) * sample.RadiusM, math.Sin(radians) * sample.RadiusM
}

func DistanceToLoad(person model.PersonDetection, sample model.CraneSample) float64 {
	x, y := LoadPosition(sample)
	return math.Hypot(person.X-x, person.Y-y)
}

func DistanceToReactor(person model.PersonDetection, reactor model.Reactor) float64 {
	radians := (reactor.AngleDeg - 90) * math.Pi / 180
	x := math.Cos(radians) * reactor.RadiusM
	y := math.Sin(radians) * reactor.RadiusM
	return math.Hypot(person.X-x, person.Y-y)
}

func IsInReactorExclusionZone(person model.PersonDetection, reactors []model.Reactor) bool {
	for _, reactor := range reactors {
		if DistanceToReactor(person, reactor) <= ReactorExclusionM {
			return true
		}
	}
	return false
}
