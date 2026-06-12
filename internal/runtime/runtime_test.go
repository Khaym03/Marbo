package runtime

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

func TestRuntime_ConfidenceAndClarification(t *testing.T) {
	cache := &Cache{
		IntentMap: map[domain.IntentID]domain.Intent{
			"intent_a": {ID: "intent_a", RequiresFlow: false, Response: domain.Response{Text: "A"}},
			"intent_b": {ID: "intent_b", RequiresFlow: false, Response: domain.Response{Text: "B"}},
			"intent_c": {ID: "intent_c", RequiresFlow: false, Response: domain.Response{Text: "C"}},
		},
		Intents: []IntentVector{
			{IntentID: "intent_a", Vector: []float32{1, 0}},
			{IntentID: "intent_b", Vector: []float32{0, 1}},
			{IntentID: "intent_c", Vector: []float32{0, 0}},
		},
	}
	settings := domain.Settings{
		SimilarityThreshold:     0.5,
		AmbiguityThreshold:      0.1,
		MaxClarificationOptions: 2,
	}

	t.Run("Strong Match", func(t *testing.T) {
		mock := MockEmbedder{vector: []float32{0.9, 0}}
		rt := NewRuntime(mock, cache, settings)
		res, err := rt.Handle("query")
		if err != nil {
			t.Fatal(err)
		}
		if res.Type != ResultAnswer || res.IntentID != "intent_a" || res.Extension.Confidence.Score < 0.89 {
			t.Errorf("expected strong match, got %v", res)
		}
	})

	t.Run("Low Confidence", func(t *testing.T) {
		mock := MockEmbedder{vector: []float32{0.1, 0}}
		rt := NewRuntime(mock, cache, settings)
		res, err := rt.Handle("query")
		if err != nil {
			t.Fatal(err)
		}
		if res.Type != ResultFallback {
			t.Errorf("expected fallback, got %v", res)
		}
	})

	t.Run("Ambiguous Match", func(t *testing.T) {
		mock := MockEmbedder{vector: []float32{0.85, 0.80}}
		rt := NewRuntime(mock, cache, settings)
		res, err := rt.Handle("query")
		if err != nil {
			t.Fatal(err)
		}
		if res.Type != ResultClarification {
			t.Errorf("expected clarification, got %v", res)
		}
		if len(res.Extension.Clarify.Candidates) != 2 {
			t.Errorf("expected 2 candidates, got %d", len(res.Extension.Clarify.Candidates))
		}
	})

	t.Run("Active Flow Bypass", func(t *testing.T) {
		cache.Flows = []domain.Flow{{ID: "flow1", StartNode: "n1", Nodes: []domain.FlowNode{{ID: "n1", Response: domain.Response{Text: "Flow"}}}}}

		mock := MockEmbedder{vector: []float32{0.1, 0}}
		rt := NewRuntime(mock, cache, settings)
		rt.state = &ConversationState{ActiveFlowID: "flow1", CurrentNodeID: "n1"}

		res, err := rt.Handle("query")
		if err != nil {
			t.Fatal(err)
		}

		if res.Type != ResultFlow {
			t.Errorf("expected flow, got %v", res)
		}
	})
}
