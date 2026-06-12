package metrics

import (
	"strings"

	"github.com/Khaym03/Marbo/internal/domain"
)

func AnalyzeIntents(kb *domain.KnowledgeBase) ([]IntentMetrics, error) {
	metrics := []IntentMetrics{}
	for _, intent := range kb.Intents {
		count := len(intent.TrainingPhrases)
		var totalWords int
		for _, p := range intent.TrainingPhrases {
			totalWords += len(strings.Fields(p))
		}
		avgLen := 0.0
		if count > 0 {
			avgLen = float64(totalWords) / float64(count)
		}

		coverage := 0.0
		switch {
		case count >= 15:
			coverage = 100
		case count >= 10:
			coverage = 80
		case count >= 5:
			coverage = 60
		default:
			coverage = 30
		}

		// Simple risk heuristic: lower coverage + low diversity
		risk := 100.0 - coverage

		metrics = append(metrics, IntentMetrics{
			IntentID:            string(intent.ID),
			TrainingPhraseCount: count,
			AveragePhraseLength: avgLen,
			CoverageScore:       coverage,
			RiskScore:           risk,
		})
	}
	return metrics, nil
}
