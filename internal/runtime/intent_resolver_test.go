package runtime

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

var settings = domain.Settings{
	SimilarityThreshold: 0.3,
	AmbiguityThreshold:  0.1,
}

func TestIntentResolver_BasicMatch(t *testing.T) {
	cache := &Cache{
		Intents: []IntentVector{
			{
				IntentID: "greeting",
				Vector:   []float32{1, 0, 0},
			},
			{
				IntentID: "farewell",
				Vector:   []float32{0, 1, 0},
			},
		},
	}

	resolver := NewIntentResolver(cache, settings)

	query := []float32{0.9, 0.1, 0}

	match, err := resolver.Resolve(query)
	if err != nil {
		t.Fatal(err)
	}
	if match.IntentID != "greeting" {
		t.Errorf("expected greeting, got %s", match.IntentID)
	}

	if match.Score <= 0 {
		t.Errorf("expected positive score, got %f", match.Score)
	}
}

func TestIntentResolver_TieHandling(t *testing.T) {
	cache := &Cache{
		Intents: []IntentVector{
			{
				IntentID: "intent_a",
				Vector:   []float32{1, 0},
			},
			{
				IntentID: "intent_b",
				Vector:   []float32{1, 0},
			},
		},
	}

	resolver := NewIntentResolver(cache, settings)

	query := []float32{1, 0}

	match, err := resolver.Resolve(query)
	if err != nil {
		t.Fatal(err)
	}
	if match.Score != 1 {
		t.Errorf("expected score 1, got %f", match.Score)
	}

	// deterministic: first match wins
	if match.IntentID != "intent_a" {
		t.Errorf("expected intent_a, got %s", match.IntentID)
	}
}

func TestIntentResolver_LowSimilarity(t *testing.T) {
	cache := &Cache{
		Intents: []IntentVector{
			{
				IntentID: "alpha",
				Vector:   []float32{1, 0},
			},
		},
	}

	resolver := NewIntentResolver(cache, settings)

	query := []float32{-1, 0}

	match, err := resolver.Resolve(query)
	if err != nil {
		t.Fatal(err)
	}
	if match.Score != -1 {
		t.Errorf("expected -1 score, got %f", match.Score)
	}
}

func TestIntentResolver_MultiPhraseIntent(t *testing.T) {
	cache := &Cache{
		Intents: []IntentVector{
			{
				IntentID: "greeting",
				Phrase:   "hello",
				Vector:   []float32{1, 0},
			},
			{
				IntentID: "greeting",
				Phrase:   "hi there",
				Vector:   []float32{0.9, 0.1},
			},
			{
				IntentID: "farewell",
				Phrase:   "bye",
				Vector:   []float32{0, 1},
			},
		},
	}

	r := NewIntentResolver(cache, settings)

	query := []float32{1, 0}

	result, err := r.ResolveTopK(query)
	if err != nil {
		t.Fatal(err)
	}
	if result.Best.IntentID != "greeting" {
		t.Errorf("expected greeting, got %s", result.Best.IntentID)
	}
}
