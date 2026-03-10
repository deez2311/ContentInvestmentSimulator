package model

import (
	"math"
	"testing"
)

func TestROI_PositiveBudget(t *testing.T) {
	m := Movie{Title: "Test", Budget: 60, Revenue: 320}
	got := m.ROI()
	want := 320.0 / 60.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("ROI() = %f, want %f", got, want)
	}
}

func TestROI_ZeroBudget(t *testing.T) {
	m := Movie{Title: "Zero", Budget: 0, Revenue: 100}
	got := m.ROI()
	if got != 0 {
		t.Errorf("ROI() = %f, want 0 for zero budget", got)
	}
}

func TestROI_ZeroRevenue(t *testing.T) {
	m := Movie{Title: "Flop", Budget: 50, Revenue: 0}
	got := m.ROI()
	if got != 0 {
		t.Errorf("ROI() = %f, want 0 for zero revenue", got)
	}
}
