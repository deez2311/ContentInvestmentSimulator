package similarity

import (
	"testing"

	"simulator/internal/llm"
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

func TestFindSimilarMovies_KnownFeatures(t *testing.T) {
	features := llm.PlotFeatures{
		Genre:    "Action",
		Themes:   []string{"Revenge"},
		Keywords: []string{"police"},
	}
	result := FindSimilarMovies(features, testMovies)
	if len(result) == 0 {
		t.Fatal("expected at least one result")
	}
	if len(result) > 3 {
		t.Fatalf("expected at most 3 results, got %d", len(result))
	}
	// Shadow Strike: genre Action=+3, theme Revenge=+2 → score 5
	if result[0].Title != "Shadow Strike" {
		t.Errorf("expected first result to be Shadow Strike, got %s", result[0].Title)
	}
}

func TestFindSimilarMovies_NoMatch(t *testing.T) {
	features := llm.PlotFeatures{
		Genre:    "Romance",
		Themes:   []string{"Musical"},
		Keywords: []string{"dance"},
	}
	result := FindSimilarMovies(features, testMovies)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d movies", len(result))
	}
}

func TestFindSimilarMovies_EmptyMovies(t *testing.T) {
	features := llm.PlotFeatures{
		Genre:    "Action",
		Themes:   []string{"Revenge"},
		Keywords: []string{"assassin"},
	}
	result := FindSimilarMovies(features, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil movies, got %d movies", len(result))
	}
}

func TestFindSimilarMovies_CaseInsensitive(t *testing.T) {
	lower := llm.PlotFeatures{Genre: "action", Themes: []string{"revenge"}, Keywords: []string{"police"}}
	upper := llm.PlotFeatures{Genre: "ACTION", Themes: []string{"REVENGE"}, Keywords: []string{"POLICE"}}
	mixed := llm.PlotFeatures{Genre: "AcTiOn", Themes: []string{"ReVeNgE"}, Keywords: []string{"PoLiCe"}}

	resultLower := FindSimilarMovies(lower, testMovies)
	resultUpper := FindSimilarMovies(upper, testMovies)
	resultMixed := FindSimilarMovies(mixed, testMovies)

	if len(resultLower) != len(resultUpper) || len(resultLower) != len(resultMixed) {
		t.Fatalf("case-insensitive mismatch: lower=%d upper=%d mixed=%d",
			len(resultLower), len(resultUpper), len(resultMixed))
	}
	for i := range resultLower {
		if resultLower[i].Title != resultUpper[i].Title || resultLower[i].Title != resultMixed[i].Title {
			t.Errorf("result mismatch at index %d: lower=%s upper=%s mixed=%s",
				i, resultLower[i].Title, resultUpper[i].Title, resultMixed[i].Title)
		}
	}
}

func TestFindSimilarMovies_ReturnsAtMost3(t *testing.T) {
	// Genre "Action" matches Shadow Strike, Iron Justice, Urban Pursuit
	// Keyword "thriller" matches Deep Harbor, Silent Witness
	// That's 5 movies with score > 0, but we should get at most 3
	features := llm.PlotFeatures{
		Genre:    "Action",
		Themes:   []string{},
		Keywords: []string{"Thriller"},
	}
	result := FindSimilarMovies(features, testMovies)
	if len(result) > 3 {
		t.Errorf("expected at most 3 results, got %d", len(result))
	}
}

func TestFindSimilarMovies_SortedByScore(t *testing.T) {
	// Shadow Strike: genre Action=+3, theme Revenge=+2 → score 5
	// Iron Justice: genre Action=+3 → score 3
	// Urban Pursuit: genre Action=+3 → score 3
	features := llm.PlotFeatures{
		Genre:  "Action",
		Themes: []string{"Revenge"},
	}
	result := FindSimilarMovies(features, testMovies)
	if len(result) < 2 {
		t.Fatal("expected at least 2 results")
	}
	if result[0].Title != "Shadow Strike" {
		t.Errorf("expected Shadow Strike first (score 5), got %s", result[0].Title)
	}
}

func TestFindSimilarMovies_GenreMatchScores3(t *testing.T) {
	features := llm.PlotFeatures{Genre: "Comedy"}
	result := FindSimilarMovies(features, testMovies)
	if len(result) != 1 {
		t.Fatalf("expected 1 result for Comedy genre, got %d", len(result))
	}
	if result[0].Title != "Last Laugh" {
		t.Errorf("expected Last Laugh, got %s", result[0].Title)
	}
}

func TestFindSimilarMovies_ThemeAndKeywordScoring(t *testing.T) {
	// Test that themes and keywords contribute to scoring
	// Deep Harbor: Genre=Thriller, Theme=Crime
	// features: Genre=Thriller (+3), Themes=Crime (+2), Keywords=Thriller (+1 via genre match)
	features := llm.PlotFeatures{
		Genre:    "Thriller",
		Themes:   []string{"Crime"},
		Keywords: []string{"Thriller"},
	}
	result := FindSimilarMovies(features, testMovies)
	if len(result) == 0 {
		t.Fatal("expected at least one result")
	}
	// Deep Harbor: genre Thriller=+3, theme Crime=+2, keyword Thriller matches genre=+1 → score 6
	if result[0].Title != "Deep Harbor" {
		t.Errorf("expected Deep Harbor first, got %s", result[0].Title)
	}
}

func TestAverageROI_ValidMovies(t *testing.T) {
	movies := []model.Movie{
		{Title: "A", Genre: "Action", Theme: "Revenge", Budget: 60, Revenue: 320},
		{Title: "B", Genre: "Action", Theme: "Police", Budget: 70, Revenue: 300},
		{Title: "C", Genre: "Action", Theme: "Police", Budget: 55, Revenue: 230},
	}
	got := AverageROI(movies)
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
	expected := 200.0 / 50.0
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
