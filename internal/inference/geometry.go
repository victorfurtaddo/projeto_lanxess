package inference

import "math"

func NormalizeAngle(angle float64) float64 {
	normalized := math.Mod(angle, 360)
	if normalized < 0 {
		normalized += 360
	}
	return normalized
}

func AngularDistance(a, b float64) float64 {
	diff := math.Abs(NormalizeAngle(a) - NormalizeAngle(b))
	if diff > 180 {
		return 360 - diff
	}
	return diff
}
