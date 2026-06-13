package runtime

import (
	"github.com/Khaym03/Marbo/internal/domain"
)

type IntentVector struct {
	IntentID domain.IntentID `json:"intent_id"`
	Phrase   string          `json:"phrase"`
	Vector   []float32       `json:"vector"`
}

type TransitionVector struct {
	FlowID     domain.FlowID `json:"flow_id"`
	NodeID     domain.NodeID `json:"node_id"`
	TargetNode domain.NodeID `json:"target_node"`
	Phrase     string        `json:"phrase"`
	Vector     []float32     `json:"vector"`
}

type Cache struct {
	Intents     []IntentVector                    `json:"intents"`
	Transitions []TransitionVector                `json:"transitions"`
	Flows       []domain.Flow                     `json:"flows"`
	IntentMap   map[domain.IntentID]domain.Intent `json:"intent_map"`
}
