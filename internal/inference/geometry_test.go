package inference

import "testing"

func TestAngularDistanceWrapAround(t *testing.T) {
	got := AngularDistance(358, 2)
	if got != 4 {
		t.Fatalf("AngularDistance() = %v, want 4", got)
	}
}

func TestNormalizeAngle(t *testing.T) {
	cases := map[float64]float64{
		360: 0,
		-10: 350,
		725: 5,
	}

	for input, want := range cases {
		if got := NormalizeAngle(input); got != want {
			t.Fatalf("NormalizeAngle(%v) = %v, want %v", input, got, want)
		}
	}
}
