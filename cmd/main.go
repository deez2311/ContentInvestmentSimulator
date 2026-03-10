package main

import (
	"bufio"
	"fmt"
	"os"

	"simulator/internal/dataset"
	"simulator/internal/evaluator"
	"simulator/internal/llm"
	"simulator/internal/optimizer"
	"simulator/internal/similarity"
)

func main() {
	movies, err := dataset.LoadMovies("data/movies.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load dataset: %v\n", err)
		os.Exit(1)
	}

	client, err := llm.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize LLM client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Enter movie plot:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	plot := scanner.Text()

	features, err := llm.AnalyzePlot(client, plot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to analyze plot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nExtracted Features: Genre=%q, Themes=%v, Keywords=%v\n", features.Genre, features.Themes, features.Keywords)

	similar := similarity.FindSimilarMovies(features, movies)

	if len(similar) == 0 {
		fmt.Println("No similar titles found for the given plot.")
		return
	}

	avgROI := similarity.AverageROI(similar)

	eval := &evaluator.ROIEvaluator{
		AvgROI:    avgROI,
		Runs:      1000,
		NoiseMean: 1.0,
		NoiseStd:  0.25,
	}

	optimal := optimizer.FindOptimalBudget(eval)

	explanation, err := llm.GenerateExplanation(client, plot, similar, avgROI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate explanation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Similar Titles")
	fmt.Println("--------------")
	for _, m := range similar {
		fmt.Println(m.Title)
	}

	fmt.Println()
	fmt.Println("Investment Recommendation")
	fmt.Println("-------------------------")
	fmt.Printf("Optimal Budget: $%.0fM\n", optimal.Budget)
	fmt.Printf("Expected ROI: %.1fx\n", optimal.Result.Mean)
	fmt.Printf("Confidence Interval: %.1fx – %.1fx\n", optimal.Result.Low, optimal.Result.High)

	fmt.Println()
	fmt.Println("Explanation")
	fmt.Println("-----------")
	fmt.Println(explanation)
}
