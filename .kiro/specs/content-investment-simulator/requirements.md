# Requirements Document

## Introduction

The Content Investment Simulator is a Go-based CLI prototype that estimates ROI for new movie titles. Given a movie plot description, the system finds historically similar titles from a synthetic CSV dataset, runs Monte Carlo simulations to model uncertainty, optimizes budget allocation, and provides explainable investment recommendations with confidence intervals. The architecture is modular with a pluggable Evaluator interface for future extensibility.

## Glossary

- **Simulator**: The Content Investment Simulator system as a whole
- **Dataset_Loader**: The component responsible for reading and parsing the CSV movie dataset
- **Similarity_Engine**: The component that finds historically similar movies based on plot keyword matching against genre and theme fields
- **ROI**: Return on Investment, calculated as revenue divided by budget for a given movie
- **Monte_Carlo_Simulator**: The component that runs stochastic simulations to model ROI uncertainty using normally distributed noise
- **Evaluator**: An interface that accepts a budget and returns a simulation result containing mean, low (P10), and high (P90) ROI estimates
- **ROI_Evaluator**: A concrete implementation of the Evaluator interface that uses Monte Carlo simulation to evaluate budget scenarios
- **Budget_Optimizer**: The component that tests candidate budgets and selects the one with the highest expected return
- **Explanation_Generator**: The component that produces human-readable reasoning for investment recommendations
- **CLI**: The command-line interface that accepts user input and displays formatted output
- **Movie**: A data record containing title, genre, theme, budget (in millions), and revenue (in millions)
- **Confidence_Interval**: A range defined by P10 and P90 percentiles from Monte Carlo simulation results

## Requirements

### Requirement 1: CSV Dataset Loading

**User Story:** As a developer, I want to load movie data from a CSV file, so that the system has historical titles to compare against.

#### Acceptance Criteria

1. WHEN a valid CSV file path is provided, THE Dataset_Loader SHALL parse the file and return a list of Movie records containing title, genre, theme, budget, and revenue fields
2. WHEN a CSV file contains malformed rows with non-numeric budget or revenue values, THE Dataset_Loader SHALL return a descriptive error indicating the row and field that failed parsing
3. IF the CSV file does not exist at the specified path, THEN THE Dataset_Loader SHALL return an error indicating the file was not found
4. IF the CSV file is empty or contains only a header row, THEN THE Dataset_Loader SHALL return an empty list of Movie records without error

### Requirement 2: Plot Similarity Matching

**User Story:** As a user, I want the system to find movies similar to my plot description, so that historical performance data can inform the investment recommendation.

#### Acceptance Criteria

1. WHEN a plot description is provided, THE Similarity_Engine SHALL match keywords in the plot against the genre and theme fields of all movies in the dataset
2. WHEN similar movies are found, THE Similarity_Engine SHALL return the top 3 movies ranked by keyword match relevance
3. IF no movies match the plot keywords, THEN THE Similarity_Engine SHALL return an empty list
4. THE Similarity_Engine SHALL perform case-insensitive matching when comparing plot keywords to genre and theme fields

### Requirement 3: ROI Calculation

**User Story:** As a user, I want to see the ROI of similar historical titles, so that I can understand the financial performance of comparable content.

#### Acceptance Criteria

1. THE Simulator SHALL compute ROI for each Movie as revenue divided by budget
2. WHEN a set of similar movies is identified, THE Simulator SHALL compute the average ROI across all similar movies in the set
3. IF a Movie has a budget of zero, THEN THE Simulator SHALL exclude that Movie from ROI calculations and log a warning

### Requirement 4: Monte Carlo Simulation

**User Story:** As a user, I want the system to model uncertainty in ROI predictions, so that I understand the range of possible outcomes.

#### Acceptance Criteria

1. THE Monte_Carlo_Simulator SHALL execute 1000 simulation runs for each budget scenario
2. THE Monte_Carlo_Simulator SHALL apply normally distributed noise with mean 1.0 and standard deviation 0.25 to each simulation run
3. THE Monte_Carlo_Simulator SHALL compute simulated revenue for each run as budget multiplied by average ROI multiplied by the noise factor
4. WHEN simulation runs are complete, THE Monte_Carlo_Simulator SHALL compute the mean ROI, P10 ROI (10th percentile), and P90 ROI (90th percentile) from the results
5. THE Monte_Carlo_Simulator SHALL return a result containing the mean, P10, and P90 values as the Confidence_Interval

### Requirement 5: Budget Optimization

**User Story:** As a user, I want the system to recommend an optimal budget, so that I can maximize expected return on investment.

#### Acceptance Criteria

1. THE Budget_Optimizer SHALL evaluate candidate budgets of 20M, 40M, 60M, 80M, 100M, and 120M
2. WHEN evaluating candidate budgets, THE Budget_Optimizer SHALL use the Evaluator interface to obtain simulation results for each candidate
3. WHEN all candidates are evaluated, THE Budget_Optimizer SHALL select the candidate budget with the highest mean expected return
4. THE Budget_Optimizer SHALL return the optimal budget along with the corresponding mean ROI and Confidence_Interval

### Requirement 6: Evaluator Interface

**User Story:** As a developer, I want a pluggable evaluation interface, so that the system can support future prediction models beyond ROI-based evaluation.

#### Acceptance Criteria

1. THE Simulator SHALL define an Evaluator interface with an Evaluate method that accepts a budget (float64) and returns a Result containing Mean, Low, and High float64 fields
2. THE ROI_Evaluator SHALL implement the Evaluator interface by running Monte Carlo simulation and returning the computed Confidence_Interval
3. THE Budget_Optimizer SHALL accept any implementation of the Evaluator interface, enabling substitution of evaluation strategies without modifying optimizer logic

### Requirement 7: Explanation Generation

**User Story:** As a user, I want a human-readable explanation of the recommendation, so that I can understand the reasoning behind the investment suggestion.

#### Acceptance Criteria

1. WHEN a recommendation is produced, THE Explanation_Generator SHALL list each similar movie used in the analysis along with the individual ROI of that movie
2. WHEN a recommendation is produced, THE Explanation_Generator SHALL display the average ROI computed across all similar movies
3. THE Explanation_Generator SHALL format the explanation as plain text suitable for CLI display
4. THE Explanation_Generator SHALL format ROI values with one decimal place followed by an "x" suffix (e.g., "4.3x")

### Requirement 8: CLI Input and Output

**User Story:** As a user, I want to enter a movie plot via the command line and receive a formatted investment recommendation, so that I can quickly evaluate content investment scenarios.

#### Acceptance Criteria

1. WHEN the CLI starts, THE CLI SHALL prompt the user with "Enter movie plot:" and accept a free-text plot description from standard input
2. WHEN a recommendation is computed, THE CLI SHALL display a "Similar Titles" section listing the titles of matched movies
3. WHEN a recommendation is computed, THE CLI SHALL display an "Investment Recommendation" section showing the optimal budget formatted as "$NNM", the expected ROI, and the Confidence_Interval
4. WHEN a recommendation is computed, THE CLI SHALL display an "Explanation" section containing the output from the Explanation_Generator
5. IF the Similarity_Engine returns no similar movies, THEN THE CLI SHALL display a message indicating no similar titles were found and skip the recommendation

### Requirement 9: CSV Dataset Format

**User Story:** As a developer, I want a well-defined CSV dataset format, so that the dataset can be extended or replaced without code changes.

#### Acceptance Criteria

1. THE Simulator SHALL read a CSV file with columns: title, genre, theme, budget, revenue in that order
2. THE Simulator SHALL treat budget and revenue values as floating-point numbers representing millions of dollars
3. THE Dataset_Loader SHALL skip the first row of the CSV file, treating the first row as a header row
