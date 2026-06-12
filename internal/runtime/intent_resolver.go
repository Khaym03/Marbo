package runtime

import (
	"sort"

	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/similarity"
)

type IntentResolver struct {
	cache     *Cache
	settings  domain.Settings
	CallCount int // For testing purposes
}

type IntentMatch struct {
	IntentID domain.IntentID
	Score    float32
}

type IntentScore struct {
	IntentID domain.IntentID
	Score    float32
}

type IntentRanking struct {
	Best       IntentMatch
	Candidates []IntentScore
	Ambiguous  bool
}

func NewIntentResolver(cache *Cache, settings domain.Settings) *IntentResolver {
	return &IntentResolver{
		cache:    cache,
		settings: settings,
	}
}

func (r *IntentResolver) Resolve(query []float32) (IntentMatch, error) {
	var (
		bestID    domain.IntentID
		bestScore float32 = -1.0
	)

	for _, intent := range r.cache.Intents {
		score, err := similarity.DotProduct(query, intent.Vector)
		if err != nil {
			return IntentMatch{}, err
		}

		if score > bestScore {
			bestScore = score
			bestID = intent.IntentID
		}
	}

	return IntentMatch{
		IntentID: bestID,
		Score:    bestScore,
	}, nil
}

func (r *IntentResolver) ResolveTopK(query []float32) (IntentRanking, error) {
	r.CallCount++
	// intentID → best score seen across phrases
	bestScores := make(map[domain.IntentID]float32)

	// 1. Score ALL phrase vectors
	for _, iv := range r.cache.Intents {

		score, err := similarity.DotProduct(query, iv.Vector)
		if err != nil {
			return IntentRanking{}, err
		}

		if existing, ok := bestScores[iv.IntentID]; !ok || score > existing {
			bestScores[iv.IntentID] = score
		}
	}

	// 2. Flatten into sortable list
	scores := make([]IntentScore, 0, len(bestScores))
	for id, score := range bestScores {
		scores = append(scores, IntentScore{
			IntentID: id,
			Score:    score,
		})
	}

	// 3. Sort descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if len(scores) == 0 {
		return IntentRanking{}, nil
	}

	best := scores[0]

	second := IntentScore{}
	if len(scores) > 1 {
		second = scores[1]
	}

	settings := r.settings

	ambiguityGap := best.Score - second.Score

	ambiguous :=
		best.Score < settings.SimilarityThreshold ||
			ambiguityGap < settings.AmbiguityThreshold

	return IntentRanking{
		Best: IntentMatch{
			IntentID: best.IntentID,
			Score:    best.Score,
		},
		Candidates: scores,
		Ambiguous:  ambiguous,
	}, nil
}
