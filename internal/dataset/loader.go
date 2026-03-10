package dataset

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"simulator/internal/model"
)

// LoadMovies reads a CSV file and returns parsed Movie records.
// Skips the header row. Returns descriptive errors for malformed rows.
// Returns empty slice (not error) for empty/header-only files.
func LoadMovies(path string) ([]model.Movie, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Empty file or header-only: return empty slice
	if len(records) <= 1 {
		return []model.Movie{}, nil
	}

	// Skip header row (index 0), parse data rows
	movies := make([]model.Movie, 0, len(records)-1)
	for i, row := range records[1:] {
		rowNum := i + 2 // 1-indexed, accounting for header

		budget, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid budget: %w", rowNum, err)
		}

		revenue, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid revenue: %w", rowNum, err)
		}

		movies = append(movies, model.Movie{
			Title:   row[0],
			Genre:   row[1],
			Theme:   row[2],
			Budget:  budget,
			Revenue: revenue,
		})
	}

	return movies, nil
}
