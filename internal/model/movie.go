package model

// Movie represents a historical movie record from the dataset.
type Movie struct {
	Title   string
	Genre   string
	Theme   string
	Budget  float64 // millions of dollars
	Revenue float64 // millions of dollars
}

// ROI computes return on investment as revenue / budget.
// Returns 0 if budget is zero to avoid division by zero.
func (m Movie) ROI() float64 {
	if m.Budget == 0 {
		return 0
	}
	return m.Revenue / m.Budget
}
