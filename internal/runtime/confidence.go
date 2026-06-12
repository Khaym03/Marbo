package runtime

import "github.com/Khaym03/Marbo/internal/domain"

type DecisionType string

const (
	DecisionExecute       DecisionType = "execute"
	DecisionClarification DecisionType = "clarification"
	DecisionFallback      DecisionType = "fallback"
)

type Decision struct {
	Type       DecisionType
	IntentID   domain.IntentID
	Confidence float32
	Candidates []domain.IntentID
}

func EvaluateConfidence(ranking IntentRanking, settings domain.Settings) Decision {
	if len(ranking.Candidates) == 0 {
		return Decision{Type: DecisionFallback}
	}

	best := ranking.Best
	if best.Score < settings.SimilarityThreshold {
		return Decision{Type: DecisionFallback}
	}

	if len(ranking.Candidates) == 1 {
		return Decision{Type: DecisionExecute, IntentID: best.IntentID, Confidence: best.Score}
	}

	second := ranking.Candidates[1]
	if (best.Score - second.Score) >= settings.AmbiguityThreshold {
		return Decision{Type: DecisionExecute, IntentID: best.IntentID, Confidence: best.Score}
	}

	candidates := []domain.IntentID{}
	limit := settings.MaxClarificationOptions
	for i := 0; i < len(ranking.Candidates) && i < limit; i++ {
		candidates = append(candidates, ranking.Candidates[i].IntentID)
	}

	return Decision{
		Type:       DecisionClarification,
		Candidates: candidates,
		Confidence: best.Score,
	}
}
