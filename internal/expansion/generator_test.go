package expansion

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/planner"
)

// Helper to monkey patch GetSeeds
type MockGetSeeds func() []IntentSeedSet

var GetSeedsFunc MockGetSeeds

func TestGenerateExpansionPack(t *testing.T) {
	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{
			{ID: "INTENT_REQUISITOS_INGRESO", Label: "Requisitos", TrainingPhrases: []string{"phrase1"}},
		},
	}
	report := &planner.ExpansionReport{
		Intents: []planner.IntentExpansionReport{
			{
				IntentID:               "INTENT_REQUISITOS_INGRESO",
				CurrentPhraseCount:     1,
				RecommendedPhraseCount: 15,
				MissingCoverage:        []string{"conversational", "first_person"},
				Priority:               planner.PriorityHigh,
			},
		},
	}

	pack := GenerateExpansionPack(kb, report)

	if len(pack.Intents) != 1 {
		t.Errorf("expected 1 intent in pack, got %d", len(pack.Intents))
	}
	if len(pack.Intents[0].GeneratedPhrases) == 0 {
		t.Error("expected generated phrases, got none")
	}
}
