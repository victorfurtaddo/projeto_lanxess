package safety

import "testing"

func TestEstimateMassRespectsDensityAndNetworkLimit(t *testing.T) {
	estimate := EstimateMass(6, 8)

	if estimate.EstimatedMassKg != 3000 {
		t.Fatalf("EstimatedMassKg = %v, want 3000", estimate.EstimatedMassKg)
	}
	if estimate.FinalMassKg != 2500 {
		t.Fatalf("FinalMassKg = %v, want 2500", estimate.FinalMassKg)
	}
	if estimate.OperationalLimitKg != 2500 {
		t.Fatalf("OperationalLimitKg = %v, want 2500", estimate.OperationalLimitKg)
	}
}

func TestCraneCapacityByRadius(t *testing.T) {
	cases := []struct {
		radius float64
		want   float64
	}{
		{radius: 2.5, want: 10000},
		{radius: 5, want: 5000},
		{radius: 8, want: 2500},
	}

	for _, tc := range cases {
		if got := CraneCapacityByRadius(tc.radius); got != tc.want {
			t.Fatalf("CraneCapacityByRadius(%v) = %v, want %v", tc.radius, got, tc.want)
		}
	}
}
