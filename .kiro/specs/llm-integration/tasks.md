# Implementation Plan: LLM Integration

## Overview

Integrate ChatGPT into the Content Investment Simulator by adding a new `internal/llm/` package (client, plot analyzer, explanation generator), updating the similarity engine to accept structured `PlotFeatures`, and wiring the LLM pipeline into `cmd/main.go`. All financial modeling remains unchanged. Property-based tests use `pgregory.net/rapid`.

## Tasks

- [ ] 1. Create LLM client package
  - [x] 1.1 Implement `llm.Client` in `internal/llm/client.go`
    - Define `Client` struct with `apiKey`, `httpClient`, `model`, and `endpoint` fields
    - Implement `NewClient()` that reads `OPENAI_API_KEY` from environment, returns descriptive error if not set
    - Implement `Chat(prompt string) (string, error)` that sends a request to the OpenAI chat completions endpoint, parses the response, and returns the text content from the first choice
    - Return error containing HTTP status code and response body on API failure
    - _Requirements: 1.1, 1.2, 1.3, 1.4_

  - [ ]* 1.2 Write property tests for `llm.Client` in `internal/llm/client_test.go`
    - **Property 1: Client returns API response text** — For any prompt string, when a mock HTTP server returns a successful response, `Chat` returns the text content from the first choice. (100 iterations)
    - **Validates: Requirements 1.2**
    - **Property 2: Client error includes status and body** — For any HTTP error status code (4xx/5xx) and response body from a mock server, the error from `Chat` contains both the status code and body text. (100 iterations)
    - **Validates: Requirements 1.4**

  - [ ]* 1.3 Write unit tests for `llm.Client` in `internal/llm/client_test.go`
    - Test `NewClient()` returns error when `OPENAI_API_KEY` is not set
    - Test `NewClient()` succeeds when `OPENAI_API_KEY` is set
    - _Requirements: 1.1, 1.3_

- [ ] 2. Implement plot analyzer
  - [x] 2.1 Implement `PlotFeatures` struct and `AnalyzePlot` function in `internal/llm/plot_analyzer.go`
    - Define `PlotFeatures` struct with `Genre string`, `Themes []string`, `Keywords []string` and JSON tags
    - Implement `AnalyzePlot(client *Client, plot string) (PlotFeatures, error)` that constructs a prompt with explicit JSON schema instructions, sends it via `Client.Chat`, and parses the JSON response
    - Return descriptive parse error if response is not valid JSON
    - Propagate client errors unchanged
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [ ]* 2.2 Write property tests for plot analyzer in `internal/llm/plot_analyzer_test.go`
    - **Property 3: PlotFeatures JSON round-trip** — For any valid `PlotFeatures`, serialize to JSON and parse back; assert equality. (100 iterations)
    - **Validates: Requirements 2.2**
    - **Property 4: Invalid JSON produces parse error** — For any non-JSON string, parsing should return a non-nil error. (100 iterations)
    - **Validates: Requirements 2.4**
    - **Property 5: Plot analyzer propagates client errors** — For any error from `Client.Chat`, `AnalyzePlot` returns that error or a wrapping error. (100 iterations)
    - **Validates: Requirements 2.3**

  - [ ]* 2.3 Write unit tests for plot analyzer in `internal/llm/plot_analyzer_test.go`
    - Test that the prompt sent to the client contains JSON schema instructions and the plot text
    - Test that `PlotFeatures` struct has no financial fields
    - _Requirements: 2.1, 2.5, 6.1_

- [x] 3. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 4. Update similarity engine for structured matching
  - [x] 4.1 Update `FindSimilarMovies` in `internal/similarity/similarity.go`
    - Change signature from `FindSimilarMovies(plot string, movies []model.Movie)` to `FindSimilarMovies(features llm.PlotFeatures, movies []model.Movie)`
    - Implement weighted scoring: genre match = +3, each theme match = +2, each keyword match = +1
    - All comparisons case-insensitive using `strings.EqualFold`
    - Return up to 3 movies sorted by descending score; exclude movies with score 0
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

  - [ ]* 4.2 Write property tests for similarity engine in `internal/similarity/similarity_test.go`
    - **Property 6: Similarity scoring formula** — For any `PlotFeatures` and `Movie`, the score equals `3*(genre match) + 2*(theme overlap count) + 1*(keyword overlap count)`, all case-insensitive. (100 iterations)
    - **Validates: Requirements 3.2**
    - **Property 7: Similarity returns at most 3 results in descending score order** — For any features and movie slice, result length ≤ 3, sorted descending by score, no zero-score movies. (100 iterations)
    - **Validates: Requirements 3.3, 3.4**
    - **Property 8: Case-insensitive matching** — For any features and movie, changing case of genre/themes/keywords does not change the score. (100 iterations)
    - **Validates: Requirements 3.5**

  - [ ]* 4.3 Write unit tests for similarity engine in `internal/similarity/similarity_test.go`
    - Test empty movie slice returns empty result
    - Test exact genre match scores 3
    - Test theme and keyword overlap scoring
    - _Requirements: 3.2, 3.3, 3.4_

- [ ] 5. Implement LLM explanation generator
  - [x] 5.1 Implement `GenerateExplanation` in `internal/llm/explanation.go`
    - Implement `GenerateExplanation(client *Client, plot string, similar []model.Movie, avgROI float64) (string, error)`
    - Construct prompt containing user plot, each similar movie's title/genre/theme, and average ROI
    - Instruct LLM to focus on narrative structure, themes, and genre; limit to 2-3 sentences
    - Do not instruct LLM to compute financial values — ROI is context only
    - Propagate client errors unchanged
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 6.2_

  - [ ]* 5.2 Write property tests for explanation generator in `internal/llm/explanation_test.go`
    - **Property 9: Explanation prompt contains all inputs** — For any plot, movie slice, and avgROI, the prompt contains the plot text, each movie title, and the ROI value. (100 iterations)
    - **Validates: Requirements 4.1**
    - **Property 10: Explanation generator propagates client errors** — For any error from `Client.Chat`, `GenerateExplanation` returns that error or a wrapping error. (100 iterations)
    - **Validates: Requirements 4.3**

  - [ ]* 5.3 Write unit tests for explanation generator in `internal/llm/explanation_test.go`
    - Test prompt contains narrative/theme focus instructions
    - Test prompt does not instruct financial computation
    - _Requirements: 4.4, 4.5, 6.2_

- [x] 6. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. Wire LLM pipeline into main
  - [x] 7.1 Update `cmd/main.go` to use LLM pipeline
    - Initialize `llm.Client` at startup; print error to stderr and `os.Exit(1)` on failure
    - Call `llm.AnalyzePlot` with user plot; exit on error
    - Pass `PlotFeatures` to updated `similarity.FindSimilarMovies` instead of raw plot string
    - Replace `explain.Generate` call with `llm.GenerateExplanation`; exit on error
    - Preserve existing output format: Similar Titles, Investment Recommendation, Explanation sections
    - Keep all financial calculations (ROIEvaluator, Monte Carlo, optimizer) unchanged
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.3_

  - [x] 7.2 Add `pgregory.net/rapid` dependency
    - Run `go get pgregory.net/rapid` to add the property-based testing library
    - _Requirements: Testing infrastructure_

- [x] 8. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property tests use `pgregory.net/rapid` with a minimum of 100 iterations each
- The existing `internal/explain/` package is preserved but no longer called from the main pipeline
- All financial modeling (ROI, Monte Carlo, optimizer) remains untouched
