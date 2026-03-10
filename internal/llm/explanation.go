package llm

import (
	"fmt"
	"strings"

	"simulator/internal/model"
)

// GenerateExplanation asks ChatGPT to explain why the similar movies
// were selected, focusing on narrative structure, themes, and genre.
// The avgROI is provided for context only — the prompt instructs the LLM
// not to compute financial values.
func GenerateExplanation(client *Client, plot string, similar []model.Movie, avgROI float64) (string, error) {
	var movieLines []string
	for _, m := range similar {
		movieLines = append(movieLines, fmt.Sprintf("- %s (%s, %s)", m.Title, m.Genre, m.Theme))
	}

	prompt := fmt.Sprintf(`User plot:
%s

Similar historical movies:
%s

Average ROI: %.1fx

Explain why these movies are similar to the user plot.
Focus on narrative structure, themes, and genre.
Do not compute or predict any financial values — the ROI is provided for context only.

Limit your explanation to 2-3 sentences.`, plot, strings.Join(movieLines, "\n"), avgROI)

	response, err := client.Chat(prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
