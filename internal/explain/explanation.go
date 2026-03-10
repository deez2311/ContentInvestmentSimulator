package explain

import (
	"fmt"
	"strings"

	"simulator/internal/model"
)

// Generate produces a human-readable explanation string listing similar movies,
// their individual ROIs (formatted as "N.Nx"), and the average ROI.
func Generate(similarMovies []model.Movie, avgROI float64) string {
	var b strings.Builder

	b.WriteString("Based on similar titles:\n\n")

	for _, m := range similarMovies {
		fmt.Fprintf(&b, "%s (ROI %.1fx)\n", m.Title, m.ROI())
	}

	fmt.Fprintf(&b, "\nAverage ROI across similar films: %.1fx\n", avgROI)

	return b.String()
}
