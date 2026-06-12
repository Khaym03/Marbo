package runtime

import (
	"fmt"
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

type MockEmbedder struct {
	vector []float32
	err    error
}

func (m MockEmbedder) Embed(
	text string,
) ([]float32, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.vector != nil {
		return m.vector, nil
	}

	return []float32{
		float32(len(text)),
	}, nil
}

type FailingEmbedder struct{}

func (m FailingEmbedder) Embed(text string) ([]float32, error) {
	return nil, fmt.Errorf("embed failed")
}

func TestBuilder_Build_Intents(t *testing.T) {

	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{
				ID: "intent_1",
				TrainingPhrases: []string{
					"hola",
					"buenas",
				},
			},
		},
	}

	builder := NewBuilder(MockEmbedder{})

	cache, err := builder.Build(kb)
	if err != nil {
		t.Fatal(err)
	}

	if len(cache.Intents) != 2 {
		t.Fatalf("expected 2 intent vectors, got %d", len(cache.Intents))
	}

	if cache.Intents[0].IntentID != "intent_1" {
		t.Fatalf("wrong intent id")
	}
}

func TestBuilder_Build_Transitions(t *testing.T) {

	kb := &domain.KnowledgeBase{
		Flows: []domain.Flow{
			{
				ID: "flow_1",
				Nodes: []domain.FlowNode{
					{
						ID: "node_1",
						Transitions: []domain.Transition{
							{
								TargetNode: "node_2",
								TrainingPhrases: []string{
									"si",
									"ok",
									"continuar",
								},
							},
						},
					},
				},
			},
		},
	}

	builder := NewBuilder(MockEmbedder{})

	cache, err := builder.Build(kb)
	if err != nil {
		t.Fatal(err)
	}

	if len(cache.Transitions) != 3 {
		t.Fatalf(
			"expected 3 transitions, got %d",
			len(cache.Transitions),
		)
	}
}

func TestBuilder_Build_FieldMapping(t *testing.T) {

	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{
				ID:              "intent_x",
				TrainingPhrases: []string{"hola"},
			},
		},
	}

	builder := NewBuilder(MockEmbedder{})

	cache, err := builder.Build(kb)
	if err != nil {
		t.Fatal(err)
	}

	v := cache.Intents[0]

	if v.IntentID != "intent_x" {
		t.Fatalf("intent id mismatch")
	}

	if v.Phrase != "hola" {
		t.Fatalf("phrase mismatch")
	}
}

func TestBuilder_Build_EmbedError(t *testing.T) {

	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{
				ID:              "intent_1",
				TrainingPhrases: []string{"hola"},
			},
		},
	}

	builder := NewBuilder(FailingEmbedder{})

	_, err := builder.Build(kb)

	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
