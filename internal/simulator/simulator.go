package simulator

import (
	"math"
	"math/rand"
	"time"

	"projeto-lanxess/internal/inference"
	"projeto-lanxess/internal/model"
	"projeto-lanxess/internal/safety"
	"projeto-lanxess/internal/vision"
)

type Config struct {
	Seed       int64
	Events     int
	ReactorN   int
	RadiusM    float64
	Start      time.Time
	SampleStep time.Duration
}

func DefaultConfig() Config {
	return Config{
		Seed:       7,
		Events:     18,
		ReactorN:   28,
		RadiusM:    8,
		Start:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
		SampleStep: 20 * time.Second,
	}
}

func Generate(cfg Config) model.SimulationResult {
	if cfg.Events <= 0 {
		cfg.Events = 1
	}
	if cfg.ReactorN <= 0 {
		cfg.ReactorN = 28
	}
	if cfg.RadiusM <= 0 {
		cfg.RadiusM = 8
	}
	if cfg.SampleStep <= 0 {
		cfg.SampleStep = 20 * time.Second
	}
	if cfg.Start.IsZero() {
		cfg.Start = DefaultConfig().Start
	}

	rng := rand.New(rand.NewSource(cfg.Seed))
	reactors := BuildReactors(cfg.ReactorN, cfg.RadiusM, cfg.Seed+100)
	samples := make([]model.CraneSample, 0, cfg.Events*8)
	truth := make([]model.TruthEvent, 0, cfg.Events)
	detector := vision.NewMockDetector(cfg.Seed + 300)

	t := cfg.Start
	angle := 270.0
	radius := 2.5
	weight := 0.0
	stepIndex := 0

	for n := 0; n < cfg.Events; n++ {
		reactor := reactors[rng.Intn(len(reactors))]
		volume := 1.1 + rng.Float64()*5.2
		mass := safety.EstimateMass(volume, reactor.RadiusM)
		amount := mass.FinalMassKg
		weight = amount

		loadSample := sample(t, angle, radius, 3.2, 0.05, weight, volume, "lifting", true, reactors, rng)
		loadSample.People = detector.DetectPeople(loadSample, stepIndex)
		markRiskZones(&loadSample)
		samples = append(samples, loadSample)
		t = t.Add(cfg.SampleStep)
		stepIndex++

		steps := 4 + rng.Intn(3)
		startAngle := angle
		startRadius := radius
		for step := 1; step <= steps; step++ {
			progress := float64(step) / float64(steps)
			angle = inference.NormalizeAngle(startAngle + shortestDelta(startAngle, reactor.AngleDeg)*progress)
			radius = startRadius + (reactor.RadiusM-startRadius)*progress
			moveSample := sample(t, angle, radius, 3.7, 1.2+rng.Float64()*0.8, weight, volume, "moving", true, reactors, rng)
			moveSample.People = detector.DetectPeople(moveSample, stepIndex)
			markRiskZones(&moveSample)
			samples = append(samples, moveSample)
			t = t.Add(cfg.SampleStep)
			stepIndex++
		}

		stopAngle := inference.NormalizeAngle(reactor.AngleDeg + rng.NormFloat64()*0.7)
		stopRadius := reactor.RadiusM + rng.NormFloat64()*0.12
		stopSample := sample(t, stopAngle, stopRadius, 2.1, 0.04, weight, volume, "positioning", true, reactors, rng)
		stopSample.People = detector.DetectPeople(stopSample, stepIndex)
		markRiskZones(&stopSample)
		samples = append(samples, stopSample)
		t = t.Add(cfg.SampleStep)
		stepIndex++

		residual := rng.Float64() * 15
		dischargeSample := sample(t, stopAngle+rng.NormFloat64()*0.15, stopRadius+rng.NormFloat64()*0.03, 1.6, 0.03, residual, 0, "discharged", false, reactors, rng)
		dischargeSample.People = detector.DetectPeople(dischargeSample, stepIndex)
		markRiskZones(&dischargeSample)
		samples = append(samples, dischargeSample)
		truth = append(truth, model.TruthEvent{
			Timestamp:       t,
			ReactorID:       reactor.ID,
			VolumeM3:        volume,
			EstimatedMassKg: mass.EstimatedMassKg,
			CraneLimitKg:    mass.CraneLimitKg,
			NetworkLimitKg:  mass.NetworkLimitKg,
			FinalMassKg:     mass.FinalMassKg,
			AmountKg:        amount - residual,
		})
		t = t.Add(cfg.SampleStep)
		stepIndex++

		weight = residual
		angle = stopAngle
		radius = stopRadius
	}

	return model.SimulationResult{
		Reactors: reactors,
		Samples:  samples,
		Truth:    truth,
	}
}

func sample(t time.Time, angle, radius, height, speed, weight, volume float64, state string, suspended bool, reactors []model.Reactor, rng *rand.Rand) model.CraneSample {
	angle = inference.NormalizeAngle(angle + rng.NormFloat64()*0.18)
	radius = math.Max(0, radius+rng.NormFloat64()*0.04)
	operationalLimit := safety.OperationalLimit(radius)
	measuredWeight := math.Max(0, weight+rng.NormFloat64()*3)
	if weight > 0 {
		measuredWeight = math.Min(measuredWeight, operationalLimit)
	}
	return model.CraneSample{
		Timestamp:          t,
		AngleDeg:           angle,
		RadiusM:            radius,
		HeightM:            math.Max(0, height+rng.NormFloat64()*0.08),
		Speed:              math.Max(0, speed+rng.NormFloat64()*0.03),
		WeightKg:           measuredWeight,
		VolumeM3:           math.Max(0, volume),
		CraneCapacityKg:    safety.CraneCapacityByRadius(radius),
		OperationalLimitKg: operationalLimit,
		SuspendedLoad:      suspended,
		State:              state,
		NearestReactorID:   nearestReactorID(angle, radius, reactors),
	}
}

func shortestDelta(from, to float64) float64 {
	delta := inference.NormalizeAngle(to) - inference.NormalizeAngle(from)
	if delta > 180 {
		delta -= 360
	}
	if delta < -180 {
		delta += 360
	}
	return delta
}

func nearestReactorID(angle, radius float64, reactors []model.Reactor) int {
	bestID := 0
	bestScore := math.MaxFloat64
	for _, reactor := range reactors {
		score := inference.AngularDistance(angle, reactor.AngleDeg) + math.Abs(radius-reactor.RadiusM)*4
		if score < bestScore {
			bestScore = score
			bestID = reactor.ID
		}
	}
	return bestID
}

func markRiskZones(sample *model.CraneSample) {
	for i := range sample.People {
		sample.People[i].InRiskZone = safety.DistanceToLoad(sample.People[i], *sample) <= safety.CriticalLoadRadiusM
	}
}
