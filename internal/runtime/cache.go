package runtime

import (
	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/embedder"
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

func PopulateCache(kb *domain.KnowledgeBase, emb embedder.Embedder) (*Cache, error) {
	cache := &Cache{
		Flows:     kb.Flows,
		IntentMap: make(map[domain.IntentID]domain.Intent),
	}

	for _, intent := range kb.Intents {
		cache.IntentMap[intent.ID] = intent
		for _, phrase := range intent.TrainingPhrases {
			vec, err := emb.Embed(phrase)
			if err != nil {
				return nil, err
			}
			cache.Intents = append(cache.Intents, IntentVector{
				IntentID: intent.ID,
				Phrase:   phrase,
				Vector:   vec,
			})
		}
	}

	for _, flow := range kb.Flows {
		for _, node := range flow.Nodes {
			for _, trans := range node.Transitions {
				for _, phrase := range trans.TrainingPhrases {
					vec, err := emb.Embed(phrase)
					if err != nil {
						return nil, err
					}
					cache.Transitions = append(cache.Transitions, TransitionVector{
						FlowID:     flow.ID,
						NodeID:     node.ID,
						TargetNode: trans.TargetNode,
						Phrase:     phrase,
						Vector:     vec,
					})
				}
			}
		}
	}

	return cache, nil
}
