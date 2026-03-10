package evaluator

import (
	"sort"

	"simulator/internal/simulation"
)

// ROIEvaluator implements Evaluator using Monte Carlo simulation to estimate ROI.
type ROIEvaluator struct {
	AvgROI    float64
	Runs      int     // default 1000
	NoiseMean float64 // default 1.0
	NoiseStd  float64 // default 0.25
}

// Evaluate runs a Monte Carlo simulation for the given budget and returns
// the mean ROI along with P10/P90 confidence bounds.
func (e *ROIEvaluator) Evaluate(budget float64) Result {
	results := simulation.Simulate(budget, e.AvgROI, e.NoiseMean, e.NoiseStd, e.Runs)

	sort.Float64s(results)

	n := len(results)
	if n == 0 {
		return Result{}
	}

	sum := 0.0
	for _, v := range results {
		sum += v
	}
	mean := sum / float64(n)

	low := simulation.ComputePercentile(results, 10)
	high := simulation.ComputePercentile(results, 90)

	return Result{
		Mean: mean,
		Low:  low,
		High: high,
	}
}
