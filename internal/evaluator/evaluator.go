package evaluator

// Result holds the output of an evaluation: mean ROI and confidence interval bounds.
type Result struct {
	Mean float64 // mean ROI from simulation
	Low  float64 // P10 ROI
	High float64 // P90 ROI
}

// Evaluator defines a pluggable evaluation strategy for budget scenarios.
type Evaluator interface {
	Evaluate(budget float64) Result
}
