# Content Investment Simulator

[Watch the demo](https://drive.google.com/file/d/1Fp0No5HGz3mO9Jtg4INEz1sEgZC3GRy7/view?usp=drive_link)

A command-line tool that answers a simple question: given a new movie plot, how much should you invest and what return can you expect?

You type in a movie plot. The system uses ChatGPT to analyze the plot into structured semantic features (genre, themes, keywords), finds historically similar titles from a dataset using weighted matching, computes their average ROI, runs a Monte Carlo simulation to model uncertainty, and recommends an optimal budget with a confidence interval. ChatGPT also generates a natural-language explanation of why those movies were selected.

## How It Works

1. You enter a movie plot as free text.
2. ChatGPT analyzes the plot and extracts structured features: genre, themes, and keywords.
3. The similarity engine scores movies from the dataset against those features using weighted matching (genre match = +3, theme match = +2 each, keyword match = +1 each). Up to 3 top-scoring movies are returned.
4. The average ROI (revenue / budget) of those similar movies is calculated.
5. A Monte Carlo simulation (1,000 runs) applies random noise around that average ROI to model real-world uncertainty.
6. Six candidate budgets ($20M–$120M) are evaluated. The one with the highest expected return is recommended.
7. ChatGPT generates a 2–3 sentence explanation focusing on narrative structure, themes, and genre similarities.

## Prerequisites

- Go 1.25+
- An OpenAI API key

### Installing Go 1.25

Download from the official site: https://go.dev/dl/

On macOS with Homebrew:

```bash
brew install go
```

Or download the installer directly:

```bash
# macOS (Apple Silicon)
curl -LO https://go.dev/dl/go1.25.0.darwin-arm64.pkg
sudo installer -pkg go1.25.0.darwin-arm64.pkg -target /

# macOS (Intel)
curl -LO https://go.dev/dl/go1.25.0.darwin-amd64.pkg
sudo installer -pkg go1.25.0.darwin-amd64.pkg -target /

# Linux (amd64)
curl -LO https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

Verify the installation:

```bash
go version
# go version go1.25.0 ...
```

### Setting the OpenAI API Key

```bash
export OPENAI_API_KEY=your-key-here
```

## Running

```bash
go run cmd/main.go
```

Or pipe input directly:

```bash
echo "A retired assassin hunts the gang that betrayed him." | go run cmd/main.go
```

## Example

Input:

```
Enter movie plot:
Jesse, who meets a French student, Céline, on a train heading through Europe. Feeling an instant connection, Jesse convinces Céline to get off the train with him in Vienna and spend the night exploring the city before his flight the next morning. As they wander through streets, cafés, parks, and landmarks, they talk deeply about life, love, relationships, and their fears about the future, gradually forming a romantic bond during their brief time together. Knowing their meeting is temporary, they part ways in the morning but agree to return to the same station in six months, leaving their reunion to fate.
```

Output:

```
Extracted Features: Genre="Romantic Drama", Themes=[chance encounter romance fate], Keywords=[train journey Vienna short romantic encounter longing]

Similar Titles
--------------
Broken Oath
Lost Fang
Lost Drift

Investment Recommendation
-------------------------
Optimal Budget: $120M
Expected ROI: 3.0x
Confidence Interval: 2.1x – 4.0x

Explanation
-----------
These movies are similar to the user plot in terms of their focus on relationships and deep emotional connections between the characters. They all explore themes of love, loss, and the passage of time, and have a strong emphasis on character development and introspection. Additionally, they all fall under the drama genre, which allows for a more nuanced exploration of the characters' emotions and experiences.
```

Note: Similar titles, ROI values, and the optimal budget will vary between runs due to the stochastic nature of the Monte Carlo simulation and LLM responses.

## Architecture

```
Plot Input
   │
   ▼
ChatGPT Plot Analyzer → PlotFeatures (genre, themes, keywords)
   │
   ▼
Weighted Similarity Engine
   │
   ▼
Similar Movies → Average ROI
   │
   ▼
Monte Carlo Simulation (1,000 runs)
   │
   ▼
Budget Optimizer ($20M–$120M)
   │
   ▼
ChatGPT Explanation Generator
```

The LLM handles creative interpretation only. All financial modeling (ROI calculation, Monte Carlo simulation, budget optimization) remains deterministic.

## Dataset

`data/movies.csv` contains ~1,500 synthetic movie records with title, genre, theme, budget (millions), and revenue (millions). This is mock data created for demonstration purposes — no machine learning models or real-world datasets are involved. The similarity matching and financial modeling are based entirely on this CSV file.

## Tests

```bash
go test ./...
```

## File List

```
.
├── go.mod
├── cmd/
│   └── main.go
├── data/
│   └── movies.csv
└── internal/
    ├── dataset/
    │   ├── loader.go
    │   └── loader_test.go
    ├── evaluator/
    │   ├── evaluator.go
    │   ├── roi.go
    │   └── roi_test.go
    ├── explain/
    │   ├── explanation.go
    │   └── explanation_test.go
    ├── llm/
    │   ├── client.go
    │   ├── plot_analyzer.go
    │   └── explanation.go
    ├── model/
    │   ├── movie.go
    │   └── movie_test.go
    ├── optimizer/
    │   ├── optimizer.go
    │   └── optimizer_test.go
    ├── similarity/
    │   ├── similarity.go
    │   └── similarity_test.go
    └── simulation/
        ├── montecarlo.go
        └── montecarlo_test.go
```
