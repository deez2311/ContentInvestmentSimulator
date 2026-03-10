# Implementation Plan: Content Investment Simulator

## Overview

Implement a Go CLI application that recommends optimal content investment budgets. The implementation follows a staged approach: project setup, core data layer, similarity matching, Monte Carlo simulation, budget optimization with a pluggable Evaluator interface, explanation generation, and CLI wiring. Each stage builds incrementally on the previous one, with property-based tests validating correctness properties from the design.

## Tasks

- [ ] 1. Project setup and data layer
  - [x] 1.1 Initialize Go module and create project structure
    - Run `go mod init simulator`
    - Create directory structure: `cmd/`, `internal/model/`, `internal/dataset/`, `internal/similarity/`, `internal/simulation/`, `internal/evaluator/`, `internal/optimizer/`, `internal/explain/`, `data/`
    - Create `data/movies.csv` with the synthetic dataset (10 movie records)
    - _Requirements: 9.1, 9.2_

  - [x] 1.2 Implement Movie model with ROI method
    - Create `internal/model/movie.go` with `Movie` struct (Title, Genre, Theme, Budget, Revenue)
    - Implement `ROI()` method returning `revenue / budget`, returning 0 if budget is zero
    - _Requirements: 3.1, 3.3_

  - [ ]* 1.3 Write property test for ROI computation
    - **Property 6: ROI Computation**
    - For any Movie with positive budget, `ROI()` equals `revenue / budget`
    - **Validates: Requirements 3.1**

  - [x] 1.4 Implement CSV dataset loader
    - Create `internal/dataset/loader.go` with `LoadMovies(path string) ([]Movie, error)`
    - Use `encoding/csv` to parse, skip header row
    - Parse budget/revenue with `strconv.ParseFloat`
    - Return descriptive error with row number and field name on parse failure
    - Return empty slice for empty/header-only files
    - Return error for file not found
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 9.1, 9.2, 9.3_

  - [ ]* 1.5 Write property tests for CSV loading
    - **Property 1: CSV Parsing Round-Trip**
    - **Validates: Requirements 1.1, 9.1, 9.2**

  - [ ]* 1.6 Write property test for malformed CSV error reporting
    - **Property 2: Malformed CSV Error Reporting**
    - **Validates: Requirements 1.2**

  - [ ]* 1.7 Write property test for header row skipping
    - **Property 3: Header Row Skipping**
    - **Validates: Requirements 9.3**

- [x] 2. Checkpoint - Verify data layer
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 3. Similarity engine and ROI averaging
  - [x] 3.1 Implement similarity engine
    - Create `internal/similarity/similarity.go` with `FindSimilarMovies(plot string, movies []Movie) []Movie`
    - Tokenize plot into lowercase words
    - Count keyword matches against `strings.ToLower(genre)` and `strings.ToLower(theme)` for each movie
    - Sort by match count descending, return top 3
    - Return empty slice if no matches
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

  - [ ]* 3.2 Write property test for similarity results validity
    - **Property 4: Similarity Results Are Valid, Bounded, and Sorted**
    - **Validates: Requirements 2.1, 2.2**

  - [ ]* 3.3 Write property test for case-insensitive matching
    - **Property 5: Case-Insensitive Similarity Matching**
    - **Validates: Requirements 2.4**

  - [x] 3.4 Implement average ROI computation
    - Add a helper function to compute average ROI from a slice of Movies, excluding zero-budget movies with a log warning
    - _Requirements: 3.2, 3.3_

  - [ ]* 3.5 Write property test for average ROI correctness
    - **Property 7: Average ROI Correctness**
    - **Validates: Requirements 3.2**

- [x] 4. Checkpoint - Verify similarity and ROI
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 5. Monte Carlo simulation and evaluator
  - [x] 5.1 Implement Monte Carlo simulation
    - Create `internal/simulation/montecarlo.go` with `Simulate(budget, avgROI, noiseMean, noiseStddev float64, runs int) []float64`
    - Implement `ComputePercentile(sorted []float64, percentile float64) float64`
    - Each run: `revenue = budget * avgROI * noise` where noise ~ N(noiseMean, noiseStddev)
    - Return slice of simulated ROI values (simulatedRevenue / budget)
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

  - [ ]* 5.2 Write property test for simulation output count
    - **Property 8: Simulation Output Count and Formula**
    - **Validates: Requirements 4.1, 4.3**

  - [ ]* 5.3 Write property test for noise distribution
    - **Property 9: Simulation Noise Distribution**
    - **Validates: Requirements 4.2**

  - [x] 5.4 Implement Evaluator interface and Result struct
    - Create `internal/evaluator/evaluator.go` with `Evaluator` interface and `Result` struct
    - _Requirements: 6.1_

  - [x] 5.5 Implement ROI Evaluator
    - Create `internal/evaluator/roi.go` with `ROIEvaluator` struct implementing `Evaluator`
    - Use Monte Carlo simulation internally, compute mean/P10/P90
    - _Requirements: 6.2_

  - [ ]* 5.6 Write property test for confidence interval ordering
    - **Property 10: Confidence Interval Ordering**
    - **Validates: Requirements 4.4, 4.5**

- [x] 6. Checkpoint - Verify simulation and evaluator
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. Budget optimizer and explanation generator
  - [x] 7.1 Implement budget optimizer
    - Create `internal/optimizer/optimizer.go` with `FindOptimalBudget(e evaluator.Evaluator) OptimalResult`
    - Evaluate candidate budgets [20, 40, 60, 80, 100, 120]
    - Select candidate with highest `Mean * Budget` (expected return)
    - Return `OptimalResult` with budget and corresponding Result
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 6.3_

  - [ ]* 7.2 Write property test for optimizer selection
    - **Property 11: Optimizer Selects Maximum Expected Return**
    - **Validates: Requirements 5.3, 5.4**

  - [x] 7.3 Implement explanation generator
    - Create `internal/explain/explanation.go` with `Generate(similarMovies []model.Movie, avgROI float64) string`
    - List each similar movie with individual ROI formatted as "N.Nx"
    - Display average ROI formatted as "N.Nx"
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [ ]* 7.4 Write property test for explanation content
    - **Property 12: Explanation Contains All Movies, ROIs, and Average**
    - **Validates: Requirements 7.1, 7.2, 7.4**

- [x] 8. Checkpoint - Verify optimizer and explanation
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. CLI wiring and demo
  - [x] 9.1 Implement CLI entry point
    - Create `cmd/main.go` orchestrating the full pipeline
    - Prompt user with "Enter movie plot:" and read from stdin via `bufio.Scanner`
    - Load dataset, find similar movies, compute average ROI
    - Construct `ROIEvaluator`, run optimizer, generate explanation
    - Print "Similar Titles", "Investment Recommendation" (budget as "$NNM", ROI, confidence interval), and "Explanation" sections
    - Handle no-similar-movies case with graceful message
    - Handle dataset load failure with stderr message and non-zero exit
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

  - [ ]* 9.2 Write property test for budget formatting
    - **Property 13: Budget Formatting**
    - **Validates: Requirements 8.3**

- [x] 10. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- All property-based tests use the `rapid` library (`github.com/flyingmutant/rapid`)
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation between stages
- The `go test ./...` command runs all tests; use `-run TestProperty` to run property tests only
