package audit

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

type MockEmbedder struct{}

func (m MockEmbedder) Embed(text string) ([]float32, error) {
	// Simple mock: return vectors that are identical for identical text
	// and very similar for similar text (based on length)
	if text == "dup" {
		return []float32{1.0, 0.0}, nil
	}
	if text == "sim1" {
		return []float32{1.0, 0.0}, nil
	}
	if text == "sim2" {
		return []float32{0.999, 0.001}, nil // Dot product ≈ 1.0, >= 0.92
	}
	return []float32{0.0, 0.0}, nil
}

func TestTransitionAuditor(t *testing.T) {
	kb := &domain.KnowledgeBase{
		Flows: []domain.Flow{{
			ID: "test_flow",
			Nodes: []domain.FlowNode{
				{
					ID: "node1",
					Transitions: []domain.Transition{
						{TargetNode: "t1", TrainingPhrases: []string{"dup", "sim1"}},
						{TargetNode: "t2", TrainingPhrases: []string{"dup", "sim2"}},
					},
				},
				{
					ID: "node2",
					Transitions: []domain.Transition{
						{TargetNode: "t3", TrainingPhrases: []string{}},
					},
				},
				{
					ID:          "node3",
					IsTerminal:  false,
					Transitions: []domain.Transition{{TargetNode: "t4", TrainingPhrases: []string{"phrase"}}},
				},
			},
		}},
	}

	auditor := TransitionAuditor{Embedder: MockEmbedder{}}
	report, err := auditor.Audit(kb)
	if err != nil {
		t.Fatal(err)
	}

	f := report.Flows[0]

	// Check Rule 1: Duplicate phrase "dup"
	foundDup := false
	for _, d := range f.DuplicatePhrases {
		if d.Phrase == "dup" {
			foundDup = true
		}
	}
	if !foundDup {
		t.Error("expected to find duplicate phrase 'dup'")
	}

	// Check Rule 2: Similar phrases "sim1" and "sim2" (dot product of [0.9, 0.1] and [0.91, 0.09] is very high)
	foundSim := false
	for _, s := range f.SimilarTransitions {
		if s.PhraseA == "sim1" && s.PhraseB == "sim2" {
			foundSim = true
		}
	}
	if !foundSim {
		t.Error("expected to find similar phrases 'sim1' and 'sim2'")
	}

	// Check Rule 4: Missing phrases
	foundMissing := false
	for _, w := range f.Warnings {
		if w == "Transition node2 -> t3 has no phrases." {
			foundMissing = true
		}
	}
	if !foundMissing {
		t.Error("expected missing phrases warning")
	}

	// Check Rule 5: One transition and not terminal
	foundSingle := false
	for _, w := range f.Warnings {
		if w == "Node node3 has only one outgoing transition and is not terminal." {
			foundSingle = true
		}
	}
	if !foundSingle {
		t.Error("expected single transition node warning")
	}
}
