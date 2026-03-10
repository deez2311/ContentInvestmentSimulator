package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"simulator/internal/model"
	"testing"
)

// newFakeServer returns an httptest.Server that responds with the given content string
// wrapped in a valid ChatGPT API response format.
func newFakeServer(content string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{
			Choices: []chatChoice{
				{Message: chatMessage{Role: "assistant", Content: content}},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestChat_Success(t *testing.T) {
	server := newFakeServer("hello world")
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	result, err := client.Chat("test prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected 'hello world', got %q", result)
	}
}

func TestChat_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid key"}`))
	}))
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	_, err := client.Chat("test")
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestChat_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatResponse{Choices: []chatChoice{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	_, err := client.Chat("test")
	if err == nil {
		t.Fatal("expected error for empty choices")
	}
}

func TestNewClient_MissingAPIKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error when OPENAI_API_KEY is not set")
	}
}

func TestAnalyzePlot_ValidResponse(t *testing.T) {
	response := `{"genre":"Action Thriller","themes":["revenge","betrayal"],"keywords":["assassin","crime"]}`
	server := newFakeServer(response)
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	features, err := AnalyzePlot(client, "A retired assassin hunts the gang that betrayed him.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if features.Genre != "Action Thriller" {
		t.Errorf("expected genre 'Action Thriller', got %q", features.Genre)
	}
	if len(features.Themes) != 2 {
		t.Errorf("expected 2 themes, got %d", len(features.Themes))
	}
	if len(features.Keywords) != 2 {
		t.Errorf("expected 2 keywords, got %d", len(features.Keywords))
	}
}

func TestAnalyzePlot_InvalidJSON(t *testing.T) {
	server := newFakeServer("this is not json")
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	_, err := AnalyzePlot(client, "some plot")
	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}
}

func TestGenerateExplanation_Success(t *testing.T) {
	server := newFakeServer("These movies share revenge themes and action genre.")
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	movies := []model.Movie{
		{Title: "Shadow Strike", Genre: "Action", Theme: "Revenge", Budget: 60, Revenue: 320},
	}
	explanation, err := GenerateExplanation(client, "test plot", movies, 5.3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if explanation == "" {
		t.Error("expected non-empty explanation")
	}
}

func TestAnalyzePlot_MarkdownFencedJSON(t *testing.T) {
	// LLMs sometimes wrap JSON in markdown code fences
	response := "```json\n{\"genre\":\"SciFi\",\"themes\":[\"exploration\"],\"keywords\":[\"space\"]}\n```"
	server := newFakeServer(response)
	defer server.Close()

	client := NewTestClient(server.URL, server.Client())
	features, err := AnalyzePlot(client, "A crew travels across the galaxy.")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if features.Genre != "SciFi" {
		t.Errorf("expected genre 'SciFi', got %q", features.Genre)
	}
}
