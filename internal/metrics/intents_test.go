package metrics

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

func TestAnalyzeIntents(t *testing.T) {
	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{ID: "i1", TrainingPhrases: make([]string, 16)}, // Coverage 100
			{ID: "i2", TrainingPhrases: make([]string, 2)},  // Coverage 30
		},
	}
	metrics, _ := AnalyzeIntents(kb)
	if metrics[0].CoverageScore != 100 {
		t.Errorf("expected 100 coverage, got %.0f", metrics[0].CoverageScore)
	}
	if metrics[1].CoverageScore != 30 {
		t.Errorf("expected 30 coverage, got %.0f", metrics[1].CoverageScore)
	}
}
