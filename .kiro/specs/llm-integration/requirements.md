# Requirements Document

## Introduction

This feature integrates a Large Language Model (ChatGPT) into the Content Investment Simulator pipeline. The LLM serves two purposes: (1) converting raw movie plot text into structured semantic features for improved similarity matching, and (2) generating human-readable explanations for investment recommendations. The LLM handles creative interpretation only — all financial modeling (ROI, Monte Carlo simulation, budget optimization) remains deterministic and unchanged.

## Glossary

- **Simulator**: The Content Investment Simulator application (`cmd/main.go` entry point)
- **LLM_Client**: A reusable HTTP client that communicates with the OpenAI ChatGPT API
- **Plot_Analyzer**: A module that uses the LLM_Client to extract structured semantic features from raw movie plot text
- **PlotFeatures**: A structured data object containing genre, themes, and keywords extracted from a movie plot
- **Similarity_Engine**: The module (`internal/similarity/similarity.go`) that finds historically similar movies
- **Explanation_Generator**: A module that uses the LLM_Client to produce natural-language reasoning for why similar movies were selected
- **Movie**: A historical movie record from the dataset containing title, genre, theme, budget, and revenue
- **ROI**: Return on Investment, calculated as revenue divided by budget
- **Pipeline**: The end-to-end flow from plot input to investment recommendation output

## Requirements

### Requirement 1: LLM Client Initialization

**User Story:** As a developer, I want a reusable ChatGPT API client, so that multiple modules can interact with the LLM without duplicating connection logic.

#### Acceptance Criteria

1. THE LLM_Client SHALL authenticate with the OpenAI API using an API key read from the `OPENAI_API_KEY` environment variable
2. WHEN a prompt is sent to the LLM_Client, THE LLM_Client SHALL return the text response from the ChatGPT API
3. IF the `OPENAI_API_KEY` environment variable is not set, THEN THE LLM_Client SHALL return a descriptive error during initialization
4. IF the ChatGPT API returns an HTTP error, THEN THE LLM_Client SHALL return an error containing the HTTP status code and response body

### Requirement 2: Plot Analysis

**User Story:** As a developer, I want to convert a raw movie plot into structured semantic features, so that the similarity engine can find more relevant historical matches.

#### Acceptance Criteria

1. WHEN a raw plot string is provided, THE Plot_Analyzer SHALL send a structured prompt to the LLM_Client requesting genre, themes, and keywords extraction
2. WHEN the LLM_Client returns a valid JSON response, THE Plot_Analyzer SHALL parse the response into a PlotFeatures object containing a genre string, a themes string array, and a keywords string array
3. IF the LLM_Client returns an error, THEN THE Plot_Analyzer SHALL propagate the error to the caller
4. IF the LLM response contains invalid JSON, THEN THE Plot_Analyzer SHALL return a descriptive parse error
5. THE Plot_Analyzer SHALL include explicit instructions in the prompt specifying the expected JSON schema with fields: genre, themes, and keywords

### Requirement 3: Structured Similarity Matching

**User Story:** As a developer, I want the similarity engine to accept structured plot features instead of raw text, so that matching is based on semantic understanding rather than keyword overlap.

#### Acceptance Criteria

1. THE Similarity_Engine SHALL accept a PlotFeatures object and a slice of Movie records as input
2. WHEN scoring a Movie against PlotFeatures, THE Similarity_Engine SHALL assign a weight of 3 for a genre match, a weight of 2 for each matching theme, and a weight of 1 for each matching keyword
3. THE Similarity_Engine SHALL return up to 3 movies sorted by descending match score
4. WHEN no movies match any feature in the PlotFeatures, THE Similarity_Engine SHALL return an empty result set
5. THE Similarity_Engine SHALL perform case-insensitive comparison for all feature matching

### Requirement 4: LLM Explanation Generation

**User Story:** As a developer, I want the LLM to generate a human-readable explanation of why similar movies were selected, so that investment recommendations are more trustworthy and interpretable.

#### Acceptance Criteria

1. WHEN similar movies and an average ROI value are provided, THE Explanation_Generator SHALL send a prompt to the LLM_Client containing the user plot, similar movie details, and average ROI
2. WHEN the LLM_Client returns a response, THE Explanation_Generator SHALL return the explanation text limited to 2-3 sentences
3. IF the LLM_Client returns an error, THEN THE Explanation_Generator SHALL return an error to the caller
4. THE Explanation_Generator SHALL instruct the LLM to focus on narrative structure, themes, and genre in the explanation
5. THE Explanation_Generator SHALL not include any financial calculations or ROI predictions in the prompt instructions — financial data is provided for context only

### Requirement 5: Pipeline Integration

**User Story:** As a developer, I want the main pipeline to use LLM-powered plot analysis and explanation generation, so that the end-to-end user experience benefits from improved similarity matching and richer explanations.

#### Acceptance Criteria

1. WHEN the Simulator starts, THE Simulator SHALL initialize the LLM_Client before processing any plot input
2. WHEN a plot is entered by the user, THE Simulator SHALL pass the plot through the Plot_Analyzer to obtain PlotFeatures before invoking the Similarity_Engine
3. THE Simulator SHALL pass the PlotFeatures object to the updated Similarity_Engine instead of the raw plot string
4. WHEN similar movies are found, THE Simulator SHALL use the Explanation_Generator to produce the explanation instead of the existing static explanation module
5. IF the LLM_Client fails to initialize, THEN THE Simulator SHALL print an error message to stderr and exit with a non-zero status code
6. THE Simulator SHALL preserve the existing output format: Similar Titles, Investment Recommendation (budget, expected ROI, confidence interval), and Explanation sections

### Requirement 6: Separation of Concerns

**User Story:** As a developer, I want a clear boundary between LLM-driven creative analysis and deterministic financial modeling, so that financial calculations remain predictable and auditable.

#### Acceptance Criteria

1. THE Plot_Analyzer SHALL produce only semantic metadata (genre, themes, keywords) and SHALL not produce budget, revenue, or ROI values
2. THE Explanation_Generator SHALL not modify or compute ROI, budget, or confidence interval values
3. THE Simulator SHALL continue to use the existing ROIEvaluator, Monte Carlo simulation, and budget optimizer for all financial calculations
