package runtime

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

type HighScoreEmbedder struct{}

func (m HighScoreEmbedder) Embed(text string) ([]float32, error) {
	return []float32{1.0, 1.0}, nil
}

func setupRuntime() *Runtime {
	cache := &Cache{
		IntentMap: map[domain.IntentID]domain.Intent{
			"intent_flow": {ID: "intent_flow", RequiresFlow: true, FlowID: "flow1"},
		},
		Flows: []domain.Flow{{
			ID:        "flow1",
			StartNode: "n1",
			Nodes: []domain.FlowNode{
				{ID: "n1", Response: domain.Response{Text: "Node 1"}, Transitions: []domain.Transition{{TargetNode: "n3"}}},
				{ID: "n3", Response: domain.Response{Text: "Node 3"}, IsTerminal: false},
				{ID: "n2", Response: domain.Response{Text: "Node 2"}, IsTerminal: true},
			},
		}},
		Transitions: []TransitionVector{
			{FlowID: "flow1", NodeID: "n1", TargetNode: "n3", Vector: []float32{1.0, 1.0}},
		},
	}
	settings := domain.Settings{SimilarityThreshold: 0.5, AmbiguityThreshold: 0.1}
	return NewRuntime(HighScoreEmbedder{}, cache, settings)
}

func TestActiveFlowSkipsIntentResolution(t *testing.T) {
	rt := setupRuntime()
	rt.state = &ConversationState{ActiveFlowID: "flow1", CurrentNodeID: "n1"}

	// Reset call count before action
	rt.intentResolver.CallCount = 0

	_, err := rt.Handle("query")
	if err != nil {
		t.Fatal(err)
	}

	if rt.intentResolver.CallCount != 0 {
		t.Errorf("expected intent resolution to NOT be called, but it was called %d times", rt.intentResolver.CallCount)
	}
}

func TestStatePersistsAcrossFlowSteps(t *testing.T) {
	rt := setupRuntime()
	// Start flow
	rt.state = &ConversationState{ActiveFlowID: "flow1", CurrentNodeID: "n1"}

	// Perform step (n1 -> n3)
	res, err := rt.Handle("something")
	if err != nil {
		t.Fatal(err)
	}

	if res.Type != ResultFlow {
		t.Fatalf("expected result type ResultFlow, got %s", res.Type)
	}

	if rt.state == nil {
		t.Fatal("expected state to persist, but it was nil")
	}

	if rt.state.ActiveFlowID != "flow1" || rt.state.CurrentNodeID != "n3" {
		t.Errorf("expected state flow1/n3, got %s/%s", rt.state.ActiveFlowID, rt.state.CurrentNodeID)
	}
}

func TestTerminalNodeClearsState(t *testing.T) {
	rt := setupRuntime()
	// Set to n2 (terminal)
	rt.state = &ConversationState{ActiveFlowID: "flow1", CurrentNodeID: "n2"}

	// Trigger step (n2 is terminal, so it should clear state)
	_, err := rt.Handle("something")
	if err != nil {
		t.Fatal(err)
	}

	if rt.state != nil {
		t.Error("expected state to be cleared (nil) on terminal node, but it was not")
	}
}

func TestFlowConversationNeverReturnsClarification(t *testing.T) {
	rt := setupRuntime()
	rt.state = &ConversationState{ActiveFlowID: "flow1", CurrentNodeID: "n1"}

	res, err := rt.Handle("query")
	if err != nil {
		t.Fatal(err)
	}

	if res.Type == ResultClarification {
		t.Error("expected flow result, got ResultClarification during active flow")
	}
}

func TestFlowConversationNeverInvokesConfidenceEngine(t *testing.T) {
	rt := setupRuntime()
	rt.state = &ConversationState{ActiveFlowID: "flow1", CurrentNodeID: "n1"}

	res, err := rt.Handle("query")
	if err != nil {
		t.Fatal(err)
	}

	if res.Type != ResultFlow {
		t.Errorf("expected flow result, got %s - check if confidence engine was triggered", res.Type)
	}
}
