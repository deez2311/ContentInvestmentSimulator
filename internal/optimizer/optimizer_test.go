package optimizer

import (
	"testing"

	"simulator/internal/evaluator"
)

// mockEvaluator returns deterministic results based on a lookup map.
type mockEvaluator struct {
	results map[float64]evaluator.Result
}

func (m *mockEvaluator) Evaluate(budget float64) evaluator.Result {
	if r, ok := m.results[budget]; ok {
		return r
	}
	return evaluator.Result{}
}

func TestFindOptimalBudget_HighestExpectedReturn(t *testing.T) {
	// Budget 60 has the highest Mean*Budget = 60*5.0 = 300
	mock := &mockEvaluator{
		results: map[float64]evaluator.Result{
			20:  {Mean: 4.0, Low: 3.0, High: 5.0}, // 80
			40:  {Mean: 4.5, Low: 3.2, High: 5.5}, // 180
			60:  {Mean: 5.0, Low: 3.5, High: 6.0}, // 300
			80:  {Mean: 3.5, Low: 2.5, High: 4.5}, // 280
			100: {Mean: 2.0, Low: 1.5, High: 2.5}, // 200
			120: {Mean: 1.5, Low: 1.0, High: 2.0}, // 180
		},
	}

	result := FindOptimalBudget(mock)

	if result.Budget != 60 {
		t.Errorf("expected budget 60, got %v", result.Budget)
	}
	if result.Result.Mean != 5.0 {
		t.Errorf("expected Mean 5.0, got %v", result.Result.Mean)
	}
	if result.Result.Low != 3.5 {
		t.Errorf("expected Low 3.5, got %v", result.Result.Low)
	}
	if result.Result.High != 6.0 {
		t.Errorf("expected High 6.0, got %v", result.Result.High)
	}
}

func TestFindOptimalBudget_HighBudgetWins(t *testing.T) {
	// All same Mean, so highest budget wins (Mean*Budget is largest).
	mock := &mockEvaluator{
		results: map[float64]evaluator.Result{
			20:  {Mean: 3.0, Low: 2.0, High: 4.0},
			40:  {Mean: 3.0, Low: 2.0, High: 4.0},
			60:  {Mean: 3.0, Low: 2.0, High: 4.0},
			80:  {Mean: 3.0, Low: 2.0, High: 4.0},
			100: {Mean: 3.0, Low: 2.0, High: 4.0},
			120: {Mean: 3.0, Low: 2.0, High: 4.0},
		},
	}

	result := FindOptimalBudget(mock)

	if result.Budget != 120 {
		t.Errorf("expected budget 120 when all Means equal, got %v", result.Budget)
	}
}

func TestFindOptimalBudget_SmallBudgetHighROI(t *testing.T) {
	// Budget 20 has Mean 20.0 → expected return 400, which beats all others.
	mock := &mockEvaluator{
		results: map[float64]evaluator.Result{
			20:  {Mean: 20.0, Low: 15.0, High: 25.0}, // 400
			40:  {Mean: 4.0, Low: 3.0, High: 5.0},    // 160
			60:  {Mean: 3.0, Low: 2.0, High: 4.0},    // 180
			80:  {Mean: 2.5, Low: 1.5, High: 3.5},    // 200
			100: {Mean: 2.0, Low: 1.0, High: 3.0},    // 200
			120: {Mean: 1.5, Low: 0.5, High: 2.5},    // 180
		},
	}

	result := FindOptimalBudget(mock)

	if result.Budget != 20 {
		t.Errorf("expected budget 20 with high ROI, got %v", result.Budget)
	}
}

func TestFindOptimalBudget_EvaluatesAllCandidates(t *testing.T) {
	called := make(map[float64]bool)
	mock := &trackingEvaluator{called: called}

	FindOptimalBudget(mock)

	for _, b := range CandidateBudgets {
		if !called[b] {
			t.Errorf("budget %v was not evaluated", b)
		}
	}
}

// trackingEvaluator records which budgets were evaluated.
type trackingEvaluator struct {
	called map[float64]bool
}

func (te *trackingEvaluator) Evaluate(budget float64) evaluator.Result {
	te.called[budget] = true
	return evaluator.Result{Mean: 1.0, Low: 0.5, High: 1.5}
}
