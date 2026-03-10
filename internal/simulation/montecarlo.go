package simulation

import (
	"math"
	"math/rand"
)

// Simulate runs n Monte Carlo iterations and returns the slice of simulated ROI values.
// Each run: revenue = budget * avgROI * noise, where noise ~ N(noiseMean, noiseStddev).
// The returned values are simulatedRevenue / budget (i.e., avgROI * noise).
func Simulate(budget, avgROI, noiseMean, noiseStddev float64, runs int) []float64 {
	if runs <= 0 {
		return []float64{}
	}
	results := make([]float64, runs)
	for i := 0; i < runs; i++ {
		noise := rand.NormFloat64()*noiseStddev + noiseMean
		revenue := budget * avgROI * noise
		if budget == 0 {
			results[i] = 0
		} else {
			results[i] = revenue / budget
		}
	}
	return results
}

// ComputePercentile returns the value at the given percentile (0-100) from a sorted slice.
// Uses linear interpolation between adjacent ranks.
// Returns 0 for an empty slice.
func ComputePercentile(sorted []float64, percentile float64) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return sorted[0]
	}

	// Clamp percentile to [0, 100]
	percentile = math.Max(0, math.Min(100, percentile))

	// Convert percentile to a 0-based fractional index
	rank := (percentile / 100) * float64(n-1)
	lower := int(math.Floor(rank))
	upper := int(math.Ceil(rank))

	if lower == upper {
		return sorted[lower]
	}

	// Linear interpolation
	frac := rank - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}
