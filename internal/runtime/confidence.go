package runtime

import (
	"github.com/Khaym03/Marbo/internal/domain"
)

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

	// OOD Detection: Flat ranking check (only if score is above threshold)
	if len(ranking.Candidates) > 1 {
		variance := calculateTop5Variance(ranking.Candidates)
		if variance < 0.0005 {
			return Decision{Type: DecisionFallback}
		}
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

func calculateTop5Variance(candidates []IntentScore) float32 {
	count := len(candidates)
	if count > 5 {
		count = 5
	}
	if count < 2 {
		return 0
	}

	var sum float32
	for i := 0; i < count; i++ {
		sum += candidates[i].Score
	}
	mean := sum / float32(count)
	var sqDiffSum float32
	for i := 0; i < count; i++ {
		sqDiffSum += (candidates[i].Score - mean) * (candidates[i].Score - mean)
	}
	return sqDiffSum / float32(count)
}
