package runtime

import (
	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/embedder"
)

type Builder struct {
	embedder embedder.Embedder
}

func NewBuilder(
	embedder embedder.Embedder,
) *Builder {
	return &Builder{
		embedder: embedder,
	}
}

func (b *Builder) Build(
	kb *domain.KnowledgeBase,
) (*Cache, error) {

	cache := &Cache{
		Intents:     make([]IntentVector, 0),
		Transitions: make([]TransitionVector, 0),
		Flows:       kb.Flows,
		IntentMap:   make(map[domain.IntentID]domain.Intent),
	}

	for _, intent := range kb.Intents {
		cache.IntentMap[intent.ID] = intent

		for _, phrase := range intent.TrainingPhrases {

			vector, err := b.embedder.Embed(phrase)
			if err != nil {
				return nil, err
			}

			cache.Intents = append(cache.Intents, IntentVector{
				IntentID: intent.ID,
				Phrase:   phrase,
				Vector:   vector,
			})
		}
	}

	// -------------------------
	// FLOWS -> TRANSITIONS
	// -------------------------
	for _, flow := range kb.Flows {

		for _, node := range flow.Nodes {

			for _, transition := range node.Transitions {

				for _, phrase := range transition.TrainingPhrases {

					vector, err := b.embedder.Embed(phrase)
					if err != nil {
						return nil, err
					}

					cache.Transitions = append(cache.Transitions, TransitionVector{
						FlowID:     flow.ID,
						NodeID:     node.ID,
						TargetNode: transition.TargetNode,
						Phrase:     phrase,
						Vector:     vector,
					})
				}
			}
		}
	}

	return cache, nil
}
