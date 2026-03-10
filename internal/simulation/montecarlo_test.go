package simulation

import (
	"math"
	"sort"
	"testing"
)

func TestSimulate_ReturnsCorrectCount(t *testing.T) {
	results := Simulate(80, 4.0, 1.0, 0.25, 1000)
	if len(results) != 1000 {
		t.Errorf("expected 1000 results, got %d", len(results))
	}
}

func TestSimulate_ZeroRuns(t *testing.T) {
	results := Simulate(80, 4.0, 1.0, 0.25, 0)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSimulate_NegativeRuns(t *testing.T) {
	results := Simulate(80, 4.0, 1.0, 0.25, -5)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSimulate_ZeroBudget(t *testing.T) {
	results := Simulate(0, 4.0, 1.0, 0.25, 10)
	for i, v := range results {
		if v != 0 {
			t.Errorf("result[%d] = %f, expected 0 for zero budget", i, v)
		}
	}
}

func TestSimulate_MeanApproximatesAvgROI(t *testing.T) {
	avgROI := 4.0
	results := Simulate(80, avgROI, 1.0, 0.25, 10000)

	var sum float64
	for _, v := range results {
		sum += v
	}
	mean := sum / float64(len(results))

	// With noise mean=1.0, the expected mean ROI ≈ avgROI
	if math.Abs(mean-avgROI) > 0.15 {
		t.Errorf("mean ROI = %.4f, expected approximately %.1f", mean, avgROI)
	}
}

func TestSimulate_ZeroAvgROI(t *testing.T) {
	results := Simulate(80, 0, 1.0, 0.25, 100)
	for i, v := range results {
		if v != 0 {
			t.Errorf("result[%d] = %f, expected 0 for zero avgROI", i, v)
		}
	}
}

func TestComputePercentile_EmptySlice(t *testing.T) {
	result := ComputePercentile([]float64{}, 50)
	if result != 0 {
		t.Errorf("expected 0 for empty slice, got %f", result)
	}
}

func TestComputePercentile_SingleElement(t *testing.T) {
	result := ComputePercentile([]float64{5.0}, 50)
	if result != 5.0 {
		t.Errorf("expected 5.0, got %f", result)
	}
}

func TestComputePercentile_KnownValues(t *testing.T) {
	// Sorted slice: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// P0 = first element
	if v := ComputePercentile(data, 0); v != 1.0 {
		t.Errorf("P0 = %f, expected 1.0", v)
	}

	// P100 = last element
	if v := ComputePercentile(data, 100); v != 10.0 {
		t.Errorf("P100 = %f, expected 10.0", v)
	}

	// P50 = median, index 4.5 → interpolation between 5 and 6 = 5.5
	if v := ComputePercentile(data, 50); math.Abs(v-5.5) > 0.001 {
		t.Errorf("P50 = %f, expected 5.5", v)
	}
}

func TestComputePercentile_P10P90(t *testing.T) {
	// Generate 1000 sorted values
	results := Simulate(80, 4.0, 1.0, 0.25, 1000)
	sort.Float64s(results)

	p10 := ComputePercentile(results, 10)
	p90 := ComputePercentile(results, 90)
	mean := 0.0
	for _, v := range results {
		mean += v
	}
	mean /= float64(len(results))

	// P10 <= mean <= P90
	if p10 > mean {
		t.Errorf("P10 (%f) > mean (%f)", p10, mean)
	}
	if mean > p90 {
		t.Errorf("mean (%f) > P90 (%f)", mean, p90)
	}
}
