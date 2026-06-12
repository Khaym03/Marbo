package validator

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

func TestValidate(t *testing.T) {
	t.Run("Valid KB", func(t *testing.T) {
		kb := &domain.KnowledgeBase{
			Zones: []domain.Zone{{ID: "Z1"}},
			Intents: []domain.Intent{
				{ID: "I1", ZoneID: "Z1", RequiresFlow: false, Response: domain.Response{Text: "OK"}, TrainingPhrases: []string{"phrase"}},
			},
		}
		if err := Validate(kb); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})

	t.Run("Duplicate IDs", func(t *testing.T) {
		kb := &domain.KnowledgeBase{
			Zones:   []domain.Zone{{ID: "Z1"}, {ID: "Z1"}},
			Intents: []domain.Intent{{ID: "I1"}, {ID: "I1"}},
			Flows:   []domain.Flow{{ID: "F1"}, {ID: "F1"}},
		}
		err := Validate(kb)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Invalid References", func(t *testing.T) {
		kb := &domain.KnowledgeBase{
			Zones: []domain.Zone{{ID: "Z1"}},
			Intents: []domain.Intent{
				{ID: "I1", ZoneID: "Z2", TrainingPhrases: []string{"p"}},
				{ID: "I2", ZoneID: "Z1", RequiresFlow: true, FlowID: "F2", TrainingPhrases: []string{"p"}},
			},
		}
		err := Validate(kb)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Flow Nodes", func(t *testing.T) {
		kb := &domain.KnowledgeBase{
			Zones: []domain.Zone{{ID: "Z1"}},
			Flows: []domain.Flow{
				{
					ID:        "F1",
					StartNode: "N99",
					Nodes: []domain.FlowNode{
						{ID: "N1"},
						{ID: "N1"},
						{ID: "N2", Transitions: []domain.Transition{{TargetNode: "N99"}}},
					},
				},
			},
		}
		err := Validate(kb)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Unreachable Node", func(t *testing.T) {
		kb := &domain.KnowledgeBase{
			Zones: []domain.Zone{{ID: "Z1"}},
			Flows: []domain.Flow{
				{
					ID:        "F1",
					StartNode: "N1",
					Nodes: []domain.FlowNode{
						{ID: "N1"},
						{ID: "N2"},
					},
				},
			},
		}
		err := Validate(kb)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
