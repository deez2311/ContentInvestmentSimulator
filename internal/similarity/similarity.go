package similarity

import (
	"log"
	"sort"
	"strings"

	"simulator/internal/model"
)

// FindSimilarMovies matches plot keywords against movie Genre and Theme fields.
// Returns up to 3 movies sorted by match relevance (descending).
// Matching is case-insensitive.
func FindSimilarMovies(plot string, movies []model.Movie) []model.Movie {
	words := strings.Fields(strings.ToLower(plot))
	if len(words) == 0 || len(movies) == 0 {
		return []model.Movie{}
	}

	type scored struct {
		movie model.Movie
		count int
	}

	var results []scored
	for _, m := range movies {
		genre := strings.ToLower(m.Genre)
		theme := strings.ToLower(m.Theme)
		count := 0
		for _, w := range words {
			if strings.Contains(genre, w) || strings.Contains(theme, w) {
				count++
			}
		}
		if count > 0 {
			results = append(results, scored{movie: m, count: count})
		}
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].count > results[j].count
	})

	limit := 3
	if len(results) < limit {
		limit = len(results)
	}

	out := make([]model.Movie, limit)
	for i := 0; i < limit; i++ {
		out[i] = results[i].movie
	}
	return out
}

// AverageROI computes the average ROI from a slice of Movies.
// Movies with zero budget are excluded and a warning is logged for each.
// Returns 0.0 if no valid movies remain after filtering.
func AverageROI(movies []model.Movie) float64 {
	var sum float64
	var count int
	for _, m := range movies {
		if m.Budget == 0 {
			log.Printf("warning: excluding movie %q from ROI calculation (zero budget)", m.Title)
			continue
		}
		sum += m.ROI()
		count++
	}
	if count == 0 {
		return 0.0
	}
	return sum / float64(count)
}
