package llm

import (
	"encoding/json"
	"fmt"
)

// PlotFeatures holds structured semantic metadata extracted from a movie plot.
type PlotFeatures struct {
	Genre    string   `json:"genre"`
	Themes   []string `json:"themes"`
	Keywords []string `json:"keywords"`
}

// AnalyzePlot sends the plot to ChatGPT with a structured prompt and parses
// the JSON response into PlotFeatures.
// Returns a parse error if the response is not valid JSON.
// Propagates client errors unchanged.
func AnalyzePlot(client *Client, plot string) (PlotFeatures, error) {
	prompt := `Analyze the following movie plot.

Extract structured metadata and return JSON only, with no additional text.

The JSON must have exactly these fields:
- "genre": a string with the primary genre (e.g. "Action Thriller")
- "themes": an array of strings with narrative themes (e.g. ["revenge", "betrayal"])
- "keywords": an array of strings with descriptive keywords (e.g. ["assassin", "crime syndicate"])

Plot:
` + plot

	response, err := client.Chat(prompt)
	if err != nil {
		return PlotFeatures{}, err
	}

	var features PlotFeatures
	if err := json.Unmarshal([]byte(response), &features); err != nil {
		return PlotFeatures{}, fmt.Errorf("failed to parse LLM response as JSON: %w", err)
	}

	return features, nil
}
