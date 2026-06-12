package runtime

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

var testSettings = domain.Settings{
	SimilarityThreshold: 0.5,
}

func TestFlowRouter_StartFlow(t *testing.T) {

	cache := &Cache{
		Flows: []domain.Flow{
			{
				ID:        "support_flow",
				StartNode: "start",
				Nodes: []domain.FlowNode{
					{
						ID: "start",
						Response: domain.Response{
							Text: "Welcome",
						},
						IsTerminal: false,
					},
				},
			},
		},
	}

	router := NewFlowRouter(cache, testSettings)

	step, ok := router.StartFlow("support_flow")
	if !ok {
		t.Fatal("expected flow to start")
	}

	if step.NodeID != "start" {
		t.Errorf("expected start node")
	}
}

func TestFlowRouter_Transition(t *testing.T) {

	cache := &Cache{
		Flows: []domain.Flow{
			{
				ID:        "flow1",
				StartNode: "n1",
				Nodes: []domain.FlowNode{
					{
						ID:       "n1",
						Response: domain.Response{Text: "Node 1"},
					},
					{
						ID:       "n2",
						Response: domain.Response{Text: "Node 2"},
					},
				},
			},
		},
		Transitions: []TransitionVector{
			{
				FlowID:     "flow1",
				NodeID:     "n1",
				TargetNode: "n2",
				Vector:     []float32{1, 0}, // Semantic vector for "go"
			},
		},
	}

	r := NewFlowRouter(cache, testSettings)

	// User input vector that matches transition
	queryEmbedding := []float32{0.9, 0}

	step, ok := r.Step("flow1", "n1", queryEmbedding)
	if !ok {
		t.Fatal("expected transition success")
	}

	if step.NodeID != "n2" {
		t.Errorf("expected n2, got %s", step.NodeID)
	}
}

func TestFlowRouter_NoMatchFallback(t *testing.T) {

	cache := &Cache{
		Flows: []domain.Flow{
			{
				ID:        "flow1",
				StartNode: "n1",
				Nodes: []domain.FlowNode{
					{
						ID:       "n1",
						Response: domain.Response{Text: "Stay here"},
					},
				},
			},
		},
	}

	r := NewFlowRouter(cache, testSettings)

	// Query embedding that doesn't match anything
	queryEmbedding := []float32{0, 0}

	step, ok := r.Step("flow1", "n1", queryEmbedding)
	if !ok {
		t.Fatal("expected fallback success")
	}

	if step.NodeID != "n1" {
		t.Errorf("expected fallback to same node")
	}
}
