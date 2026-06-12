package runtime

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

// Re-define here to ensure it uses the desired behavior
type ConversationMockEmbedder struct {
	vector func(count int) []float32
	count  int
}

func (m *ConversationMockEmbedder) Embed(text string) ([]float32, error) {
	vec := m.vector(m.count)
	m.count++
	return vec, nil
}

func TestRuntime_ConversationState(t *testing.T) {
	// Setup: Flow with 2 nodes (node1 -> node2)
	cache := &Cache{
		IntentMap: map[domain.IntentID]domain.Intent{
			"start_flow": {
				ID:           "start_flow",
				RequiresFlow: true,
				FlowID:       "flow1",
			},
		},
		Intents: []IntentVector{
			{
				IntentID: "start_flow",
				Vector:   []float32{0, 1},
			},
		},
		Flows: []domain.Flow{
			{
				ID:        "flow1",
				StartNode: "node1",
				Nodes: []domain.FlowNode{
					{
						ID:       "node1",
						Response: domain.Response{Text: "Step 1"},
					},
					{
						ID:         "node2",
						Response:   domain.Response{Text: "Step 2 (Terminal)"},
						IsTerminal: true,
					},
				},
			},
		},
		Transitions: []TransitionVector{
			{
				FlowID:     "flow1",
				NodeID:     "node1",
				TargetNode: "node2",
				Vector:     []float32{1, 0},
			},
		},
	}

	mock := &ConversationMockEmbedder{
		vector: func(count int) []float32 {
			if count == 0 || count == 2 {
				return []float32{0, 1} // start_flow
			}
			return []float32{1, 0} // transition
		},
	}
	settings := domain.Settings{SimilarityThreshold: 0.1}

	runtime := NewRuntime(mock, cache, settings)

	// 1. Message 1: Start Flow
	t.Log("Message 1: Start Flow")
	result, err := runtime.Handle("start")
	if err != nil {
		t.Fatal(err)
	}

	if result.Type != ResultFlow || result.FlowID != "flow1" || result.NodeID != "node1" {
		t.Errorf("expected to start flow1 at node1, got %v", result)
	}
	if result.Response.Text != "Step 1" {
		t.Errorf("expected 'Step 1', got '%s'", result.Response.Text)
	}

	// 2. Message 2: Continue Flow
	t.Log("Message 2: Continue Flow")
	result, err = runtime.Handle("continue")
	if err != nil {
		t.Fatal(err)
	}

	// Flow ended, Type is ResultFlow, checking if it was terminal
	if result.Type != ResultFlow || result.NodeID != "node2" {
		t.Errorf("expected continued flow to node2, got %v", result)
	}
	if result.Response.Text != "Step 2 (Terminal)" {
		t.Errorf("expected 'Step 2 (Terminal)', got '%s'", result.Response.Text)
	}

	// 3. Message 3: After Terminal
	t.Log("Message 3: After Terminal")
	result, err = runtime.Handle("start again")
	if err != nil {
		t.Fatal(err)
	}

	if result.Type != ResultFlow {
		t.Errorf("expected to start a new flow, got %v", result)
	}
}

func TestRuntime_NonFlowIntent(t *testing.T) {
	cache := &Cache{
		IntentMap: map[domain.IntentID]domain.Intent{
			"info": {
				ID:           "info",
				RequiresFlow: false,
				Response:     domain.Response{Text: "Information"},
			},
		},
		Intents: []IntentVector{
			{
				IntentID: "info",
				Vector:   []float32{1, 0},
			},
		},
	}

	mock := &ConversationMockEmbedder{
		vector: func(count int) []float32 {
			return []float32{1, 0}
		},
	}
	settings := domain.Settings{SimilarityThreshold: 0.1}

	runtime := NewRuntime(mock, cache, settings)

	result, err := runtime.Handle("info")
	if err != nil {
		t.Fatal(err)
	}

	if result.Type != ResultAnswer {
		t.Errorf("expected answer, got %v", result)
	}
	if runtime.state != nil {
		t.Errorf("expected state to be nil, got %v", runtime.state)
	}
}
