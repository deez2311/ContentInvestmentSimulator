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
// All comparisons are case-insensitive using strings.Contains for partial matching
// (e.g. LLM genre "Action Thriller" matches dataset genre "Action").
// Empty feature strings are skipped to avoid false-positive matches.
func FindSimilarMovies(features llm.PlotFeatures, movies []model.Movie) []model.Movie {
	if len(movies) == 0 {
		return []model.Movie{}
	}

	type scored struct {
		movie model.Movie
		score int
	}

	featGenre := strings.ToLower(features.Genre)

	var results []scored
	for _, m := range movies {
		score := 0

		movieGenre := strings.ToLower(m.Genre)
		movieTheme := strings.ToLower(m.Theme)

		// Genre match: +3 (partial match — e.g. "Romance" matches "Romantic Drama")
		if featGenre != "" && (strings.Contains(featGenre, movieGenre) || strings.Contains(movieGenre, featGenre)) {
			score += 3
		}

		// Theme overlap: +2 for each matching theme (partial match)
		for _, theme := range features.Themes {
			t := strings.ToLower(theme)
			if t != "" && (strings.Contains(t, movieTheme) || strings.Contains(movieTheme, t)) {
				score += 2
			}
		}

		// Keyword overlap: +1 for each keyword matching genre or theme (partial match)
		for _, kw := range features.Keywords {
			k := strings.ToLower(kw)
			if k != "" && (strings.Contains(k, movieGenre) || strings.Contains(movieGenre, k) ||
				strings.Contains(k, movieTheme) || strings.Contains(movieTheme, k)) {
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
	for i := range limit {
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
