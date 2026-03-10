# Content Investment Simulator

### Multi-Stage Development Plan (POC)

## Project Goal

Build a prototype system that answers:

> **Given a new movie plot, how much should Netflix invest and what ROI should we expect?**

The system will:

1. Accept a **movie plot**
2. Find **similar historical titles**
3. Estimate **expected ROI**
4. Simulate **budget vs ROI**
5. Recommend **optimal investment**
6. Provide **explanations**
7. Output **confidence intervals**
8. Be **extensible for future ML models**

Language: **Go**

Dataset: **Synthetic CSV**

Estimated implementation time: **~2 hours**

---

# High-Level Architecture

```
Plot Input
   |
Similarity Engine
   |
Historical Dataset (CSV)
   |
Evaluator Interface
   |
ROI Evaluator
   |
Monte Carlo Simulation
   |
Budget Optimizer
   |
Explanation Generator
   |
Confidence Interval Estimator
```

This architecture demonstrates:

* modular design
* pluggable evaluation logic
* explainable predictions
* uncertainty modeling

---

# Project Structure

```
content-investment-simulator

go.mod

cmd/
    main.go

data/
    movies.csv

internal/

    model/
        movie.go

    dataset/
        loader.go

    similarity/
        similarity.go

    simulation/
        montecarlo.go

    evaluator/
        evaluator.go
        roi.go

    optimizer/
        optimizer.go

    explain/
        explanation.go
```

---

# Dataset

File:

```
data/movies.csv
```

Example:

```
title,genre,theme,budget,revenue
Shadow Strike,Action,Revenge,60,320
Iron Justice,Action,Police,70,300
Midnight Heist,Crime,Heist,50,210
Deep Harbor,Thriller,Crime,45,180
Dragon Realm,Fantasy,Adventure,120,500
Last Laugh,Comedy,Family,20,90
Urban Pursuit,Action,Police,55,230
Silent Witness,Thriller,Mystery,40,160
Cosmic Voyage,SciFi,Exploration,110,480
Broken Oath,Drama,Justice,30,100
```

Derived metric:

```
ROI = revenue / budget
```

---

# Stage 0 — Project Setup

Estimated Time: **10 minutes**

## Tasks

Create project folders:

```
cmd/
internal/
data/
```

Initialize Go module:

```
go mod init simulator
```

Add dataset file:

```
data/movies.csv
```

Verify the project compiles:

```
go run cmd/main.go
```

---

# Stage 1 — Minimal Working System (MVP)

Estimated Time: **35–45 minutes**

Goal:

> Accept a movie plot and produce a basic investment recommendation.

## Pipeline

```
Plot
 ↓
Find Similar Movies
 ↓
Compute Average ROI
 ↓
Simulate ROI
 ↓
Recommend Budget
```

## Tasks

### Dataset Loader

File:

```
internal/dataset/loader.go
```

Function:

```
LoadMovies(path string) ([]Movie, error)
```

Responsibilities:

* read CSV
* parse movie records
* return movie structs

---

### Similarity Engine

File:

```
internal/similarity/similarity.go
```

Function:

```
FindSimilarMovies(plot string, movies []Movie) []Movie
```

Simple keyword matching:

```
plot text
  ↓
match genre or theme
  ↓
return top 3 movies
```

---

### ROI Calculation

For each similar movie:

```
ROI = revenue / budget
```

Compute:

```
Average ROI across similar movies
```

---

### Budget Simulation

Test candidate budgets:

```
20M
40M
60M
80M
100M
120M
```

For each:

```
simulate ROI
```

Choose the highest expected return.

---

# Stage 2 — Decision Intelligence

Estimated Time: **30–40 minutes**

Goal:

> Provide reasoning for the recommendation.

Executives must understand **why** a model made a prediction.

## Tasks

### Explanation Generator

File:

```
internal/explain/explanation.go
```

Explain:

* which similar movies were used
* their ROI
* the average ROI

Example output:

```
Based on similar titles:

Shadow Strike (ROI 5.3x)
Iron Justice (ROI 4.2x)
Urban Pursuit (ROI 4.1x)

Average ROI across similar films: 4.5x
```

---

# Stage 3 — Uncertainty Modeling

Estimated Time: **20–25 minutes**

Goal:

> Show prediction uncertainty.

Instead of a single number:

```
ROI = 4.2x
```

Return:

```
Expected ROI: 4.2x
Confidence Range: 2.8x – 5.5x
```

## Tasks

### Monte Carlo Simulation

File:

```
internal/simulation/montecarlo.go
```

Simulation logic:

```
1000 runs
```

Noise model:

```
Normal distribution
mean = 1
stddev = 0.25
```

Revenue simulation:

```
revenue = budget * avgROI * noise
```

Compute:

```
mean ROI
P10 ROI
P90 ROI
```

---

# Stage 4 — Platform Architecture

Estimated Time: **15–25 minutes**

Goal:

> Make evaluation logic extensible.

Introduce:

```
Evaluator Interface
```

Architecture:

```
Optimizer
   |
Evaluator Interface
   |
ROI Evaluator
```

## Implementation

Create file:

```
internal/evaluator/evaluator.go
```

Interface:

```
type Evaluator interface {
    Evaluate(budget float64) Result
}
```

Result struct:

```
type Result struct {
    Mean float64
    Low  float64
    High float64
}
```

---

### ROI Evaluator

File:

```
internal/evaluator/roi.go
```

Responsibility:

* run simulation
* compute confidence interval
* return evaluation result

---

### Update Optimizer

File:

```
internal/optimizer/optimizer.go
```

Instead of passing ROI directly, pass an evaluator:

```
FindOptimalBudget(e Evaluator)
```

This allows future models:

```
ML predictor
engagement model
subscriber growth model
risk model
```

---

# Stage 5 — CLI Demo

Estimated Time: **5–10 minutes**

Goal:

Make the CLI output clear and impressive.

Example run:

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
Urban Pursuit


Investment Recommendation
-------------------------

Optimal Budget: $80M
Expected ROI: 4.3x
Confidence Interval: 2.9x – 5.6x


Explanation
-----------

Based on similar titles:

Shadow Strike (ROI 5.3x)
Iron Justice (ROI 4.2x)
Urban Pursuit (ROI 4.1x)

Average ROI across similar films: 4.5x
```

---

# Final Architecture

```
Plot Input
   |
Similarity Engine
   |
Historical Dataset
   |
Evaluator Interface
   |
ROI Evaluator
   |
Monte Carlo Simulation
   |
Optimizer
   |
Explanation Engine
```

---

# Estimated Total Development Time

| Stage                 | Time       |
| --------------------- | ---------- |
| Project Setup         | 10 minutes |
| MVP Implementation    | 40 minutes |
| Explanation Layer     | 30 minutes |
| Uncertainty Modeling  | 20 minutes |
| Platform Architecture | 20 minutes |
| Demo Output           | 10 minutes |

Total:

```
~2 hours
```

---

# Future Extensions (For Interview Discussion)

This architecture can support:

### ML Prediction

```
script embedding → performance prediction
```

### Knowledge Graph Integration

```
actors
directors
franchise relationships
```

### Engagement Metrics

```
watch hours
completion rate
```

### Netflix Decision Metrics

```
subscriber acquisition
retention impact
member value
```

---

# One Sentence Project Description

Use this during interviews:

> I built a prototype content investment simulator that estimates ROI for new titles by analyzing historically similar content, running Monte Carlo simulations to model uncertainty, and optimizing budget allocation through a modular evaluation framework.

This highlights:

* decision intelligence
* uncertainty modeling
* extensible system architecture
