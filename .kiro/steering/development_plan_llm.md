# LLM Integration Development Plan

## Content Investment Simulator – ChatGPT Enhancement

## Objective

Improve the quality of movie similarity matching and explanations by integrating **ChatGPT (LLM)** into the pipeline.

The LLM will be used for two purposes:

1. **Plot Understanding**

   * Convert a raw movie plot into structured semantic features.
   * This allows the similarity engine to find more relevant historical titles.

2. **Explanation Generation**

   * Generate human-readable reasoning explaining why historical movies were selected.
   * Improve trust and interpretability of the investment recommendation.

Important design principle:

```
LLM → creative interpretation
Simulation engine → financial modeling
```

The LLM **must not control financial calculations** such as ROI or budget optimization.

---

# Current System Pipeline

```
Plot Input
   ↓
Keyword Similarity
   ↓
Similar Movies
   ↓
Average ROI
   ↓
Monte Carlo Simulation
   ↓
Budget Recommendation
   ↓
Basic Explanation
```

---

# Target System Pipeline

```
Plot Input
   ↓
LLM Plot Analyzer
   ↓
Structured Plot Features
   ↓
Similarity Engine
   ↓
Similar Movies
   ↓
ROI Evaluator
   ↓
Monte Carlo Simulation
   ↓
Budget Optimizer
   ↓
LLM Explanation Generator
```

---

# New Components

Two new modules will be introduced.

```
internal/llm/
    client.go
    plot_analyzer.go
    explanation.go
```

---

# Stage 1 — ChatGPT Client

Estimated Time: **15 minutes**

Create a reusable OpenAI/ChatGPT client.

File:

```
internal/llm/client.go
```

Responsibilities:

* Call ChatGPT API
* Handle authentication
* Return text responses

Environment variable:

```
OPENAI_API_KEY
```

Example implementation outline:

```go
type Client struct {
    APIKey string
}

func NewClient() *Client
func (c *Client) Chat(prompt string) (string, error)
```

The client will send prompts to the ChatGPT API and return the response.

---

# Stage 2 — Plot Analyzer

Estimated Time: **25 minutes**

Purpose:

Convert a movie plot into **structured semantic features**.

File:

```
internal/llm/plot_analyzer.go
```

Define a feature struct:

```go
type PlotFeatures struct {
    Genre    string
    Themes   []string
    Keywords []string
}
```

Function:

```go
func AnalyzePlot(client *Client, plot string) (PlotFeatures, error)
```

---

## Prompt Design

Prompt example:

```
Analyze the following movie plot.

Extract structured metadata and return JSON.

Fields:
- genre
- themes (array)
- keywords (array)

Plot:
<plot>
```

Expected output:

```json
{
  "genre": "Action Thriller",
  "themes": ["revenge", "betrayal"],
  "keywords": ["assassin", "crime syndicate"]
}
```

Parse the JSON response into `PlotFeatures`.

---

# Stage 3 — Update Similarity Engine

Estimated Time: **25 minutes**

Modify the similarity engine to accept structured plot features.

Current:

```go
FindSimilarMovies(plot string, movies []Movie)
```

New:

```go
FindSimilarMovies(features PlotFeatures, movies []Movie)
```

Similarity scoring example:

```
score =
  genre match * 3 +
  theme overlap * 2 +
  keyword overlap
```

Example scoring logic:

```
genre match → +3
theme match → +2 each
keyword match → +1 each
```

Return the **top 3 highest scoring movies**.

---

# Stage 4 — LLM Explanation Generator

Estimated Time: **20 minutes**

Use ChatGPT to generate a natural-language explanation for similarity.

File:

```
internal/llm/explanation.go
```

Function:

```go
func GenerateExplanation(
    client *Client,
    plot string,
    similar []model.Movie,
    avgROI float64,
) (string, error)
```

---

## Prompt Design

Prompt example:

```
User plot:
<plot>

Similar historical movies:
- <title> (<genre>, <theme>)
- <title> (<genre>, <theme>)

Average ROI: <value>

Explain why these movies are similar to the user plot.
Focus on narrative structure, themes, and genre.

Limit explanation to 2–3 sentences.
```

Example output:

```
The plot resembles action revenge narratives like Shadow Strike and Iron Justice. 
These films feature lone protagonists confronting criminal organizations, a pattern that has historically performed well with audiences. 
Mid-budget action thrillers often deliver strong ROI due to global appeal and repeat viewership.
```

---

# Stage 5 — Update Main Pipeline

Estimated Time: **10 minutes**

Modify `cmd/main.go`.

Current flow:

```
plot
↓
FindSimilarMovies(plot)
```

New flow:

```
plot
↓
LLM AnalyzePlot
↓
FindSimilarMovies(features)
```

Implementation outline:

```go
client := llm.NewClient()

features, err := llm.AnalyzePlot(client, plot)

similar := similarity.FindSimilarMovies(features, movies)
```

Then replace explanation generation:

```go
explanation, _ := llm.GenerateExplanation(client, plot, similar, avgROI)
```

---

# Updated Main Pipeline

```
plot input
   │
   ▼
ChatGPT Plot Analyzer
   │
   ▼
structured plot features
   │
   ▼
Similarity Engine
   │
   ▼
Similar Movies
   │
   ▼
ROI Evaluator
   │
   ▼
Monte Carlo Simulation
   │
   ▼
Budget Optimizer
   │
   ▼
ChatGPT Explanation Generator
```

---

# Testing Strategy

Test scenarios:

### Test 1 — Revenge Plot

Input:

```
A retired assassin hunts the gang that betrayed him.
```

Expected result:

* Similar movies: action / revenge themes

---

### Test 2 — SciFi Exploration

Input:

```
A crew travels across the galaxy searching for a lost civilization.
```

Expected result:

* Similar movies: sci-fi / exploration

---

### Test 3 — Family Comedy

Input:

```
A struggling father starts a bakery with his kids.
```

Expected result:

* Similar movies: comedy / family

---

# Expected Output Example

```
Similar Titles
--------------

Shadow Strike
Iron Justice
Urban Pursuit

Historical Average ROI: 4.5x


Investment Recommendation
-------------------------

Optimal Budget: $80M

Expected ROI: 4.3x

Confidence Interval:
2.9x – 5.6x


Explanation
-----------

The plot resembles action revenge stories such as Shadow Strike and Iron Justice. 
These films feature lone protagonists confronting criminal organizations, a structure that historically performs well with audiences. 
Mid-budget action thrillers often achieve strong ROI due to global appeal.
```

---

# Estimated Development Time

| Stage                    | Time       |
| ------------------------ | ---------- |
| ChatGPT client           | 15 minutes |
| Plot analyzer            | 25 minutes |
| Similarity engine update | 25 minutes |
| Explanation generator    | 20 minutes |
| Pipeline integration     | 10 minutes |

Total:

```
~1.5 hours
```

---

# Future Improvements

Possible future upgrades:

### Embedding Similarity

```
plot → embedding
movie dataset → embeddings
cosine similarity search
```

### Knowledge Graph

Add relationships:

```
actor
director
franchise
```

### Feature Storage

Store structured metadata in dataset:

```
keywords
tone
setting
```

### ML-Based ROI Predictor

Replace heuristic ROI with trained model.

---

# Interview Framing

Explain the design like this:

> LLMs interpret creative input such as plots and narrative themes, while deterministic simulation models evaluate financial outcomes like ROI and investment risk.

This highlights:

* **AI-assisted creative analysis**
* **deterministic financial modeling**
* **explainable decision systems**
