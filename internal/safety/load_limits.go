package safety

import (
	"math"

	"projeto-lanxess/internal/model"
)

const (
	SteelScrapDensityKgM3 = 500.0
	LiftingNetLimitKg     = 2500.0
	ShortRadiusLimitM     = 3.0
	MediumRadiusLimitM    = 6.0
)

func CraneCapacityByRadius(radiusM float64) float64 {
	switch {
	case radiusM <= ShortRadiusLimitM:
		return 10000
	case radiusM <= MediumRadiusLimitM:
		return 5000
	default:
		return 2500
	}
}

func OperationalLimit(radiusM float64) float64 {
	return math.Min(CraneCapacityByRadius(radiusM), LiftingNetLimitKg)
}

func EstimateMass(volumeM3, radiusM float64) model.MassEstimate {
	estimated := math.Max(0, volumeM3) * SteelScrapDensityKgM3
	craneLimit := CraneCapacityByRadius(radiusM)
	operationalLimit := math.Min(craneLimit, LiftingNetLimitKg)
	finalMass := math.Min(estimated, operationalLimit)

	limitedBy := "volume"
	if finalMass == LiftingNetLimitKg && finalMass == craneLimit && estimated > finalMass {
		limitedBy = "rede_icamento_e_grua"
	} else if finalMass == LiftingNetLimitKg && estimated > LiftingNetLimitKg {
		limitedBy = "rede_icamento"
	} else if finalMass == craneLimit && estimated > craneLimit {
		limitedBy = "capacidade_grua"
	}

	return model.MassEstimate{
		VolumeM3:           volumeM3,
		DensityKgM3:        SteelScrapDensityKgM3,
		EstimatedMassKg:    estimated,
		NetworkLimitKg:     LiftingNetLimitKg,
		CraneLimitKg:       craneLimit,
		OperationalLimitKg: operationalLimit,
		FinalMassKg:        finalMass,
		LimitedBy:          limitedBy,
	}
}
