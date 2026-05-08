package safety

import (
	"fmt"

	"projeto-lanxess/internal/model"
)

func EvaluateSample(sample model.CraneSample, reactors []model.Reactor) []model.Alert {
	alerts := make([]model.Alert, 0)
	limit := OperationalLimit(sample.RadiusM)

	if sample.WeightKg > limit {
		alerts = append(alerts, newAlert("CRITICAL", "Carga acima da capacidade permitida para o raio atual.", sample, map[string]any{
			"load_kg":              sample.WeightKg,
			"operational_limit_kg": limit,
			"crane_capacity_kg":    CraneCapacityByRadius(sample.RadiusM),
			"lifting_net_limit_kg": LiftingNetLimitKg,
			"radius_m":             sample.RadiusM,
		}))
	} else if sample.WeightKg >= limit*0.9 && sample.WeightKg > 0 {
		alerts = append(alerts, newAlert("WARNING", "Carga proxima ao limite operacional.", sample, map[string]any{
			"load_kg":              sample.WeightKg,
			"operational_limit_kg": limit,
			"radius_m":             sample.RadiusM,
		}))
	}

	for _, person := range sample.People {
		distance := DistanceToLoad(person, sample)
		if sample.SuspendedLoad && sample.WeightKg > 50 && distance <= CriticalLoadRadiusM {
			alerts = append(alerts, newAlert("CRITICAL", "Pessoa detectada sob carga suspensa.", sample, map[string]any{
				"person_id":          person.ID,
				"distance_to_load_m": distance,
			}))
		}

		if IsInReactorExclusionZone(person, reactors) {
			alerts = append(alerts, newAlert("WARNING", "Pessoa detectada em zona de exclusao de reator.", sample, map[string]any{
				"person_id": person.ID,
			}))
		}

		if sample.Speed > 0.5 && sample.SuspendedLoad && distance <= CriticalLoadRadiusM*1.8 {
			alerts = append(alerts, newAlert("CRITICAL", "Movimentacao da grua com pessoa em area critica.", sample, map[string]any{
				"person_id":          person.ID,
				"distance_to_load_m": distance,
				"speed":              sample.Speed,
			}))
		}
	}

	return alerts
}

func newAlert(level, message string, sample model.CraneSample, metadata map[string]any) model.Alert {
	return model.Alert{
		ID:        fmt.Sprintf("alert_%s_%03d", sample.Timestamp.Format("150405"), len(metadata)+len(message)),
		Level:     level,
		Message:   message,
		Timestamp: sample.Timestamp,
		Metadata:  metadata,
	}
}
