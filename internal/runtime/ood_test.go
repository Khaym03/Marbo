package runtime

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

type TestMockEmbedder struct {
	vectorMap map[string][]float32
}

func (m *TestMockEmbedder) Embed(text string) ([]float32, error) {
	return m.vectorMap[text], nil
}

func TestOODDetection(t *testing.T) {
	// Case 1: "soy gay" (OOD, Expected: Fallback)
	// Case 2: "capital de francia" (OOD, Expected: Fallback)
	// Case 3: "inscribirme" (ID, Expected: Clarification)
	
	cache := &Cache{
		IntentMap: map[domain.IntentID]domain.Intent{
			"intent1": {ID: "intent1", Label: "Inscripción"},
			"intent2": {ID: "intent2", Label: "Otro"},
		},
		Intents: []IntentVector{
			{IntentID: "intent1", Phrase: "inscribirme", Vector: []float32{1.0, 0.0, 0.0}},
			{IntentID: "intent1", Phrase: "inscripcion", Vector: []float32{1.0, 0.1, 0.0}},
			{IntentID: "intent2", Phrase: "otro", Vector: []float32{0.0, 1.0, 0.0}},
			{IntentID: "intent2", Phrase: "test", Vector: []float32{0.0, 0.9, 0.0}},
		},
	}
	settings := domain.Settings{
		SimilarityThreshold:     0.75,
		AmbiguityThreshold:      0.05,
		MaxClarificationOptions: 2,
	}
	
	// Setup vectors to simulate ID and OOD
	mock := &TestMockEmbedder{
		vectorMap: map[string][]float32{
			// ID queries (high variance)
			"inscribirme": []float32{1.0, 0.0, 0.0},
			// OOD queries (flat/low variance)
			"soy gay":            []float32{0.5, 0.5, 0.0}, // Flat
			"capital de francia": []float32{0.5, 0.51, 0.0}, // Flat
		},
	}
	
	runtime := NewRuntime(mock, cache, settings)
	
	tests := []struct {
		query    string
		expected RuntimeResultType
	}{
		{"soy gay", ResultFallback},
		{"capital de francia", ResultFallback},
		{"inscribirme", ResultAnswer},
	}
	
	for _, tt := range tests {
		res, err := runtime.Handle(tt.query)
		if err != nil {
			t.Errorf("Error handling %s: %v", tt.query, err)
			continue
		}
		if res.Type != tt.expected {
			t.Errorf("Query '%s': expected %s, got %s", tt.query, tt.expected, res.Type)
		}
	}
}
