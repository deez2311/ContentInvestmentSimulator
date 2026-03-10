package optimizer

import "simulator/internal/evaluator"

// CandidateBudgets is the fixed set of budgets (in millions) to evaluate.
var CandidateBudgets = []float64{20, 40, 60, 80, 100, 120}

// OptimalResult holds the recommended budget and its evaluation.
type OptimalResult struct {
	Budget float64
	Result evaluator.Result
}

// FindOptimalBudget evaluates candidate budgets [20, 40, 60, 80, 100, 120]
// using the provided Evaluator and returns the one with the highest
// expected return (Mean * Budget).
func FindOptimalBudget(e evaluator.Evaluator) OptimalResult {
	var best OptimalResult
	bestReturn := -1.0

	for _, budget := range CandidateBudgets {
		result := e.Evaluate(budget)
		expectedReturn := result.Mean * budget
		if expectedReturn > bestReturn {
			bestReturn = expectedReturn
			best = OptimalResult{Budget: budget, Result: result}
		}
	}

	return best
}
