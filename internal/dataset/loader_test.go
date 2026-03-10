package dataset

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMovies_ValidFile(t *testing.T) {
	movies, err := LoadMovies("../../data/movies.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(movies) != 1559 {
		t.Fatalf("expected 1559 movies, got %d", len(movies))
	}
	// Spot-check first record
	if movies[0].Title != "Shadow Strike" {
		t.Errorf("expected first title 'Shadow Strike', got %q", movies[0].Title)
	}
	if movies[0].Budget != 60 {
		t.Errorf("expected budget 60, got %f", movies[0].Budget)
	}
	if movies[0].Revenue != 320 {
		t.Errorf("expected revenue 320, got %f", movies[0].Revenue)
	}
}

func TestLoadMovies_FileNotFound(t *testing.T) {
	_, err := LoadMovies("nonexistent.csv")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestLoadMovies_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")
	os.WriteFile(path, []byte(""), 0644)

	movies, err := LoadMovies(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(movies) != 0 {
		t.Fatalf("expected 0 movies, got %d", len(movies))
	}
}

func TestLoadMovies_HeaderOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "header.csv")
	os.WriteFile(path, []byte("title,genre,theme,budget,revenue\n"), 0644)

	movies, err := LoadMovies(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(movies) != 0 {
		t.Fatalf("expected 0 movies, got %d", len(movies))
	}
}

func TestLoadMovies_MalformedBudget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.csv")
	content := "title,genre,theme,budget,revenue\nTest,Action,Revenge,abc,100\n"
	os.WriteFile(path, []byte(content), 0644)

	_, err := LoadMovies(path)
	if err == nil {
		t.Fatal("expected error for malformed budget")
	}
	errMsg := err.Error()
	if !contains(errMsg, "row 2") || !contains(errMsg, "budget") {
		t.Errorf("error should mention row number and field name, got: %s", errMsg)
	}
}

func TestLoadMovies_MalformedRevenue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.csv")
	content := "title,genre,theme,budget,revenue\nTest,Action,Revenge,50,xyz\n"
	os.WriteFile(path, []byte(content), 0644)

	_, err := LoadMovies(path)
	if err == nil {
		t.Fatal("expected error for malformed revenue")
	}
	errMsg := err.Error()
	if !contains(errMsg, "row 2") || !contains(errMsg, "revenue") {
		t.Errorf("error should mention row number and field name, got: %s", errMsg)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
