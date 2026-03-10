package explain

import (
	"strings"
	"testing"

	"simulator/internal/model"
)

func TestGenerate_KnownMovies(t *testing.T) {
	movies := []model.Movie{
		{Title: "Shadow Strike", Genre: "Action", Theme: "Revenge", Budget: 60, Revenue: 320},
		{Title: "Iron Justice", Genre: "Action", Theme: "Police", Budget: 70, Revenue: 300},
		{Title: "Urban Pursuit", Genre: "Action", Theme: "Police", Budget: 55, Revenue: 230},
	}
	avgROI := 4.6

	result := Generate(movies, avgROI)

	// Check header
	if !strings.Contains(result, "Based on similar titles:") {
		t.Error("expected header 'Based on similar titles:'")
	}

	// Check each movie with ROI
	if !strings.Contains(result, "Shadow Strike (ROI 5.3x)") {
		t.Errorf("expected 'Shadow Strike (ROI 5.3x)', got:\n%s", result)
	}
	if !strings.Contains(result, "Iron Justice (ROI 4.3x)") {
		t.Errorf("expected 'Iron Justice (ROI 4.3x)', got:\n%s", result)
	}
	if !strings.Contains(result, "Urban Pursuit (ROI 4.2x)") {
		t.Errorf("expected 'Urban Pursuit (ROI 4.2x)', got:\n%s", result)
	}

	// Check average ROI
	if !strings.Contains(result, "Average ROI across similar films: 4.6x") {
		t.Errorf("expected 'Average ROI across similar films: 4.6x', got:\n%s", result)
	}
}

func TestGenerate_EmptyMovies(t *testing.T) {
	result := Generate(nil, 0)

	if !strings.Contains(result, "Based on similar titles:") {
		t.Error("expected header even with empty movies")
	}
	if !strings.Contains(result, "Average ROI across similar films: 0.0x") {
		t.Errorf("expected average ROI 0.0x for empty input, got:\n%s", result)
	}
}

func TestGenerate_SingleMovie(t *testing.T) {
	movies := []model.Movie{
		{Title: "Cosmic Voyage", Genre: "SciFi", Theme: "Exploration", Budget: 110, Revenue: 480},
	}

	result := Generate(movies, 4.4)

	if !strings.Contains(result, "Cosmic Voyage (ROI 4.4x)") {
		t.Errorf("expected 'Cosmic Voyage (ROI 4.4x)', got:\n%s", result)
	}
	if !strings.Contains(result, "Average ROI across similar films: 4.4x") {
		t.Errorf("expected average ROI line, got:\n%s", result)
	}
}
