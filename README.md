# Content Investment Simulator

A command-line tool that answers a simple question: given a new movie plot, how much should you invest and what return can you expect?

You type in a movie plot. The system finds historically similar titles from a dataset, computes their average ROI, runs a Monte Carlo simulation to model uncertainty, and recommends an optimal budget with a confidence interval.

## How It Works

1. You enter a movie plot as free text.
2. The similarity engine matches keywords from your plot against the genre and theme of movies in the dataset, returning up to 3 similar titles.
3. The average ROI (revenue / budget) of those similar movies is calculated.
4. A Monte Carlo simulation (1,000 runs) applies random noise around that average ROI to model real-world uncertainty.
5. Six candidate budgets ($20M–$120M) are evaluated. The one with the highest expected return is recommended.
6. The output includes the recommendation, a P10–P90 confidence interval, and an explanation of which movies drove the prediction.

## Example

Input:

```
Enter movie plot:
A retired assassin hunts the gang that betrayed him.
```

Output:

```
Similar Titles
--------------
Shadow Strike
Iron Justice
Dragon Realm

Investment Recommendation
-------------------------
Optimal Budget: $120M
Expected ROI: 4.6x
Confidence Interval: 3.1x – 6.0x

Explanation
-----------
Based on similar titles:

Shadow Strike (ROI 5.3x)
Iron Justice (ROI 4.3x)
Dragon Realm (ROI 4.2x)

Average ROI across similar films: 4.6x
```

Note: ROI values and the optimal budget will vary slightly between runs due to the stochastic nature of the Monte Carlo simulation.

## Running

```
go run cmd/main.go
```

## Dataset

`data/movies.csv` contains 10 synthetic movie records with title, genre, theme, budget (millions), and revenue (millions).

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
