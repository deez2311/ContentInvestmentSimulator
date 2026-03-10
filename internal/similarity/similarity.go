package similarity

import (
	"log"
	"sort"
	"strings"

	"simulator/internal/llm"
	"simulator/internal/model"
)

// FindSimilarMovies scores movies against PlotFeatures using weighted matching:
//
//	genre match = +3, each theme match = +2, each keyword match = +1
//
// Returns up to 3 movies sorted by descending score. Movies with score 0 are excluded.
// All comparisons are case-insensitive using strings.EqualFold.
func FindSimilarMovies(features llm.PlotFeatures, movies []model.Movie) []model.Movie {
	if len(movies) == 0 {
		return []model.Movie{}
	}

	type scored struct {
		movie model.Movie
		score int
	}

	var results []scored
	for _, m := range movies {
		score := 0

		// Genre match: +3
		if strings.EqualFold(m.Genre, features.Genre) {
			score += 3
		}

		// Theme overlap: +2 for each matching theme
		for _, theme := range features.Themes {
			if strings.EqualFold(m.Theme, theme) {
				score += 2
			}
		}

		// Keyword overlap: +1 for each keyword matching genre or theme
		for _, kw := range features.Keywords {
			if strings.EqualFold(m.Genre, kw) || strings.EqualFold(m.Theme, kw) {
				score += 1
			}
		}

		if score > 0 {
			results = append(results, scored{movie: m, score: score})
		}
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	limit := min(3, len(results))

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
