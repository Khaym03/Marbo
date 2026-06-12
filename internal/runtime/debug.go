package runtime

import (
	"fmt"

	"github.com/Khaym03/Marbo/internal/domain"
)

type ConfidenceDebug struct {
	Query string

	BestIntentID string
	BestScore    float32

	SecondIntentID string
	SecondScore    float32

	Gap float32

	SimilarityThreshold float32
	AmbiguityThreshold  float32

	Decision string

	Ranking []DebugCandidate
}

type DebugCandidate struct {
	IntentID string
	Score    float32
}

func BuildConfidenceDebug(
	query string,
	candidates []IntentScore,
	settings domain.Settings,
	decision Decision,
) ConfidenceDebug {
	debug := ConfidenceDebug{
		Query:               query,
		SimilarityThreshold: settings.SimilarityThreshold,
		AmbiguityThreshold:  settings.AmbiguityThreshold,
		Decision:            string(decision.Type),
	}

	if len(candidates) > 0 {
		debug.BestIntentID = string(candidates[0].IntentID)
		debug.BestScore = candidates[0].Score
	}

	if len(candidates) > 1 {
		debug.SecondIntentID = string(candidates[1].IntentID)
		debug.SecondScore = candidates[1].Score
		debug.Gap = debug.BestScore - debug.SecondScore
	} else {
		debug.Gap = 0
	}

	for _, cand := range candidates {
		debug.Ranking = append(debug.Ranking, DebugCandidate{
			IntentID: string(cand.IntentID),
			Score:    cand.Score,
		})
	}

	return debug
}

func PrintConfidenceDebug(debug ConfidenceDebug) {
	fmt.Println("=================================================")
	fmt.Println("CONFIDENCE DEBUG")
	fmt.Println("=================================================")
	fmt.Printf("Query:\n\n%s\n\n", debug.Query)
	fmt.Println("Ranking:")
	fmt.Println()
	for i, cand := range debug.Ranking {
		fmt.Printf("%d. %s -> %.4f\n", i+1, cand.IntentID, cand.Score)
	}
	fmt.Printf("\nBest Score:\n%.4f\n\n", debug.BestScore)
	fmt.Printf("Second Score:\n%.4f\n\n", debug.SecondScore)
	fmt.Printf("Gap:\n%.4f\n\n", debug.Gap)
	fmt.Printf("Similarity Threshold:\n%.4f\n\n", debug.SimilarityThreshold)
	fmt.Printf("Ambiguity Threshold:\n%.4f\n\n", debug.AmbiguityThreshold)
	fmt.Printf("Decision:\n%s\n", debug.Decision)
	fmt.Println("=================================================")
}
