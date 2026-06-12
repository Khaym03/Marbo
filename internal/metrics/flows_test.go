package metrics

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

func TestAnalyzeFlows(t *testing.T) {
	kb := &domain.KnowledgeBase{
		Flows: []domain.Flow{{
			ID: "f1",
			Nodes: []domain.FlowNode{
				{ID: "n1", Transitions: []domain.Transition{{TargetNode: "n2"}}},
				{ID: "n2", Transitions: []domain.Transition{}},
			},
		}},
	}
	metrics, _ := AnalyzeFlows(kb)
	if metrics[0].NodeCount != 2 {
		t.Errorf("expected 2 nodes, got %d", metrics[0].NodeCount)
	}
	if metrics[0].TransitionCount != 1 {
		t.Errorf("expected 1 transition, got %d", metrics[0].TransitionCount)
	}
	if metrics[0].MaxDepth != 2 {
		t.Errorf("expected depth 2, got %d", metrics[0].MaxDepth)
	}
}
