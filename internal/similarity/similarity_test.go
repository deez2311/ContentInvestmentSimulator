package similarity

import (
	"testing"

	"simulator/internal/model"
)

var testMovies = []model.Movie{
	{Title: "Shadow Strike", Genre: "Action", Theme: "Revenge", Budget: 60, Revenue: 320},
	{Title: "Iron Justice", Genre: "Action", Theme: "Police", Budget: 70, Revenue: 300},
	{Title: "Midnight Heist", Genre: "Crime", Theme: "Heist", Budget: 50, Revenue: 210},
	{Title: "Deep Harbor", Genre: "Thriller", Theme: "Crime", Budget: 45, Revenue: 180},
	{Title: "Dragon Realm", Genre: "Fantasy", Theme: "Adventure", Budget: 120, Revenue: 500},
	{Title: "Last Laugh", Genre: "Comedy", Theme: "Family", Budget: 20, Revenue: 90},
	{Title: "Urban Pursuit", Genre: "Action", Theme: "Police", Budget: 55, Revenue: 230},
	{Title: "Silent Witness", Genre: "Thriller", Theme: "Mystery", Budget: 40, Revenue: 160},
	{Title: "Cosmic Voyage", Genre: "SciFi", Theme: "Exploration", Budget: 110, Revenue: 480},
	{Title: "Broken Oath", Genre: "Drama", Theme: "Justice", Budget: 30, Revenue: 100},
}

func TestFindSimilarMovies_KnownPlot(t *testing.T) {
	result := FindSimilarMovies("action revenge", testMovies)
	if len(result) == 0 {
		t.Fatal("expected at least one result")
	}
	if len(result) > 3 {
		t.Fatalf("expected at most 3 results, got %d", len(result))
	}
	// "Shadow Strike" matches both "action" (genre) and "revenge" (theme) → 2 matches
	if result[0].Title != "Shadow Strike" {
		t.Errorf("expected first result to be Shadow Strike, got %s", result[0].Title)
	}
}

func TestFindSimilarMovies_NoMatch(t *testing.T) {
	result := FindSimilarMovies("romance musical", testMovies)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d movies", len(result))
	}
}

func TestFindSimilarMovies_EmptyPlot(t *testing.T) {
	result := FindSimilarMovies("", testMovies)
	if len(result) != 0 {
		t.Errorf("expected empty result for empty plot, got %d movies", len(result))
	}
}

func TestFindSimilarMovies_EmptyMovies(t *testing.T) {
	result := FindSimilarMovies("action revenge", nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil movies, got %d movies", len(result))
	}
}

func TestFindSimilarMovies_CaseInsensitive(t *testing.T) {
	lower := FindSimilarMovies("action revenge", testMovies)
	upper := FindSimilarMovies("ACTION REVENGE", testMovies)
	mixed := FindSimilarMovies("AcTiOn ReVeNgE", testMovies)

	if len(lower) != len(upper) || len(lower) != len(mixed) {
		t.Fatalf("case-insensitive mismatch: lower=%d upper=%d mixed=%d", len(lower), len(upper), len(mixed))
	}
	for i := range lower {
		if lower[i].Title != upper[i].Title || lower[i].Title != mixed[i].Title {
			t.Errorf("result mismatch at index %d: lower=%s upper=%s mixed=%s",
				i, lower[i].Title, upper[i].Title, mixed[i].Title)
		}
	}
}

func TestFindSimilarMovies_ReturnsAtMost3(t *testing.T) {
	// "action" matches Shadow Strike, Iron Justice, Urban Pursuit (all Action genre)
	result := FindSimilarMovies("action", testMovies)
	if len(result) > 3 {
		t.Errorf("expected at most 3 results, got %d", len(result))
	}
}

func TestFindSimilarMovies_SortedByRelevance(t *testing.T) {
	// "action revenge" → Shadow Strike has 2 matches (action+revenge), others have 1
	result := FindSimilarMovies("action revenge", testMovies)
	if len(result) < 2 {
		t.Fatal("expected at least 2 results")
	}
	// First result should have the most matches
	if result[0].Title != "Shadow Strike" {
		t.Errorf("expected Shadow Strike first (2 matches), got %s", result[0].Title)
	}
}

func TestFindSimilarMovies_PartialWordMatch(t *testing.T) {
	// "crim" should match "Crime" genre/theme via strings.Contains
	result := FindSimilarMovies("crim", testMovies)
	if len(result) == 0 {
		t.Fatal("expected partial word match for 'crim' against 'Crime'")
	}
}

func TestAverageROI_ValidMovies(t *testing.T) {
	movies := []model.Movie{
		{Title: "A", Genre: "Action", Theme: "Revenge", Budget: 60, Revenue: 320},
		{Title: "B", Genre: "Action", Theme: "Police", Budget: 70, Revenue: 300},
		{Title: "C", Genre: "Action", Theme: "Police", Budget: 55, Revenue: 230},
	}
	got := AverageROI(movies)
	// Expected: (320/60 + 300/70 + 230/55) / 3
	expected := (320.0/60.0 + 300.0/70.0 + 230.0/55.0) / 3.0
	if diff := got - expected; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("AverageROI = %f, want %f", got, expected)
	}
}

func TestAverageROI_ExcludesZeroBudget(t *testing.T) {
	movies := []model.Movie{
		{Title: "Valid", Budget: 50, Revenue: 200},
		{Title: "ZeroBudget", Budget: 0, Revenue: 100},
	}
	got := AverageROI(movies)
	expected := 200.0 / 50.0 // only the valid movie
	if diff := got - expected; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("AverageROI = %f, want %f", got, expected)
	}
}

func TestAverageROI_AllZeroBudget(t *testing.T) {
	movies := []model.Movie{
		{Title: "A", Budget: 0, Revenue: 100},
		{Title: "B", Budget: 0, Revenue: 200},
	}
	got := AverageROI(movies)
	if got != 0.0 {
		t.Errorf("AverageROI = %f, want 0.0", got)
	}
}

func TestAverageROI_EmptySlice(t *testing.T) {
	got := AverageROI([]model.Movie{})
	if got != 0.0 {
		t.Errorf("AverageROI = %f, want 0.0", got)
	}
}

func TestAverageROI_NilSlice(t *testing.T) {
	got := AverageROI(nil)
	if got != 0.0 {
		t.Errorf("AverageROI = %f, want 0.0", got)
	}
}
