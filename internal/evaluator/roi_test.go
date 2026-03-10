package evaluator

import (
	"math"
	"testing"
)

func TestROIEvaluator_Evaluate_DefaultParams(t *testing.T) {
	e := &ROIEvaluator{
		AvgROI:    4.0,
		Runs:      1000,
		NoiseMean: 1.0,
		NoiseStd:  0.25,
	}

	result := e.Evaluate(80)

	// Mean should be approximately AvgROI (within tolerance for 1000 runs)
	if math.Abs(result.Mean-4.0) > 0.5 {
		t.Errorf("expected mean ~4.0, got %f", result.Mean)
	}

	// Confidence interval ordering: Low <= Mean <= High
	if result.Low > result.Mean {
		t.Errorf("expected Low <= Mean, got Low=%f Mean=%f", result.Low, result.Mean)
	}
	if result.Mean > result.High {
		t.Errorf("expected Mean <= High, got Mean=%f High=%f", result.Mean, result.High)
	}

	// Low should be less than High (spread exists with noise)
	if result.Low >= result.High {
		t.Errorf("expected Low < High, got Low=%f High=%f", result.Low, result.High)
	}
}

func TestROIEvaluator_Evaluate_ZeroRuns(t *testing.T) {
	e := &ROIEvaluator{
		AvgROI:    4.0,
		Runs:      0,
		NoiseMean: 1.0,
		NoiseStd:  0.25,
	}

	result := e.Evaluate(80)

	if result.Mean != 0 || result.Low != 0 || result.High != 0 {
		t.Errorf("expected zero result for zero runs, got %+v", result)
	}
}

func TestROIEvaluator_Evaluate_ZeroBudget(t *testing.T) {
	e := &ROIEvaluator{
		AvgROI:    4.0,
		Runs:      100,
		NoiseMean: 1.0,
		NoiseStd:  0.25,
	}

	result := e.Evaluate(0)

	if result.Mean != 0 {
		t.Errorf("expected mean 0 for zero budget, got %f", result.Mean)
	}
}

func TestROIEvaluator_ImplementsEvaluator(t *testing.T) {
	var _ Evaluator = &ROIEvaluator{}
}
