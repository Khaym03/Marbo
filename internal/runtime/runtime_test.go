package runtime

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Khaym03/Marbo/internal/domain"
)

func TestRuntime_ClarificationEnrichment(t *testing.T) {
	cache := &Cache{
		IntentMap: map[domain.IntentID]domain.Intent{
			"intent1": {ID: "intent1", Label: "Label 1"},
			"intent2": {ID: "intent2", Label: "Label 2"},
		},
		Intents: []IntentVector{
			{IntentID: "intent1", Vector: []float32{1, 0}},
			{IntentID: "intent2", Vector: []float32{1, 0}},
		},
	}
	settings := domain.Settings{
		SimilarityThreshold:     0.9,
		AmbiguityThreshold:      0.1,
		MaxClarificationOptions: 2, // Set to 2
	}
	mock := &ConversationMockEmbedder{
		vector: func(count int) []float32 {
			return []float32{1, 0}
		},
	}
	runtime := NewRuntime(mock, cache, settings)
	result, err := runtime.Handle("test")
	if err != nil {
		t.Fatal(err)
	}

	if result.Type != ResultClarification {
		t.Fatalf("expected clarification, got %s", result.Type)
	}

	if len(result.Extension.Clarify.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(result.Extension.Clarify.Options))
	}

	for _, opt := range result.Extension.Clarify.Options {
		if opt.IntentID == "intent1" && opt.Label != "Label 1" {
			t.Errorf("expected Label 1, got %s", opt.Label)
		}
		if opt.IntentID == "intent2" && opt.Label != "Label 2" {
			t.Errorf("expected Label 2, got %s", opt.Label)
		}
	}
}

func TestRuntime_JSONSerialization(t *testing.T) {
	res := RuntimeResult{
		Type: ResultAnswer,
		Response: domain.Response{
			Text: "Answer",
		},
		IntentID: "intent1",
		Extension: &RuntimeExtension{
			Confidence: &ConfidenceData{Score: 0.9},
			Clarify: &ClarificationData{
				Options: []ClarificationOption{{IntentID: "intent1", Label: "Label 1"}},
			},
		},
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	jsonStr := string(data)

	// Verify field names (must be lowercase as per requirements)
	expected := []string{
		"\"type\"",
		"\"response\"",
		"\"intent_id\"",
		"\"extension\"",
		"\"confidence\"",
		"\"clarify\"",
		"\"options\"",
		"\"label\"",
	}
	for _, e := range expected {
		if !strings.Contains(jsonStr, e) {
			t.Errorf("JSON missing expected field %s: %s", e, jsonStr)
		}
	}
}

func TestExampleSerialization(t *testing.T) {
	examples := map[string]RuntimeResult{
		"ResultAnswer": {
			Type:     ResultAnswer,
			IntentID: "intent_a",
			Response: domain.Response{Text: "Hello!"},
			Extension: &RuntimeExtension{
				Confidence: &ConfidenceData{Score: 0.95},
			},
		},
		"ResultFlow": {
			Type:     ResultFlow,
			IntentID: "intent_b",
			FlowID:   "flow_1",
			NodeID:   "node_1",
			Response: domain.Response{Text: "Starting flow..."},
		},
		"ResultClarification": {
			Type: ResultClarification,
			Extension: &RuntimeExtension{
				Confidence: &ConfidenceData{Score: 0.5},
				Clarify: &ClarificationData{
					Options: []ClarificationOption{
						{IntentID: "intent_1", Label: "Requisitos de inscripción"},
						{IntentID: "intent_2", Label: "Fechas de examen"},
					},
				},
			},
		},
		"ResultFallback": {
			Type: ResultFallback,
		},
	}

	for name, res := range examples {
		data, _ := json.MarshalIndent(res, "", "  ")
		t.Logf("--- %s ---\n%s\n", name, string(data))
	}
}
