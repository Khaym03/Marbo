package planner

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

func TestAnalyzeCoverage(t *testing.T) {
	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{ID: "i1", TrainingPhrases: make([]string, 3)},  // High Priority
			{ID: "i2", TrainingPhrases: make([]string, 12)}, // Low Priority
		},
	}
	report, _ := AnalyzeCoverage(kb)

	// Test Priority
	if report.Intents[0].Priority != PriorityHigh {
		t.Errorf("expected high priority for <5 phrases, got %s", report.Intents[0].Priority)
	}
	if report.Intents[1].Priority != PriorityLow {
		t.Errorf("expected low priority for >=10 phrases, got %s", report.Intents[1].Priority)
	}

	// Test Category Detection (Simple Check)
	kb = &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{ID: "i3", TrainingPhrases: []string{"qué requisitos necesito"}},
		},
	}
	report, _ = AnalyzeCoverage(kb)

	foundDirectQuestion := false
	for _, cat := range report.Intents[0].MissingCoverage {
		if cat == "direct_question" {
			foundDirectQuestion = true
		}
	}
	if foundDirectQuestion {
		t.Error("expected direct_question category to be detected, but it was listed as missing")
	}
}
