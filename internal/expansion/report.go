package expansion

import (
	"encoding/json"
	"os"
	"time"
)

type ExpansionPack struct {
	GeneratedAt time.Time             `json:"generated_at"`
	Intents     []IntentExpansionPack `json:"intents"`
}

type IntentExpansionPack struct {
	IntentID            string   `json:"intent_id"`
	ExistingPhraseCount int      `json:"existing_phrase_count"`
	TargetPhraseCount   int      `json:"target_phrase_count"`
	GeneratedPhrases    []string `json:"generated_phrases"`
}

func SaveExpansionPack(path string, pack *ExpansionPack) error {
	data, err := json.MarshalIndent(pack, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
