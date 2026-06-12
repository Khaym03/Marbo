// Package planner provides functions to analyze and report on knowledge base coverage.
package planner

import (
	"slices"
	"strings"

	"github.com/Khaym03/Marbo/internal/domain"
)

func AnalyzeCoverage(kb *domain.KnowledgeBase) (*ExpansionReport, error) {
	report := &ExpansionReport{}

	for _, intent := range kb.Intents {
		count := len(intent.TrainingPhrases)
		var priority ExpansionPriority
		recommended := 15

		if count < 5 {
			priority = PriorityHigh
		} else if count < 10 {
			priority = PriorityMedium
		} else {
			priority = PriorityLow
			recommended = count
		}

		categories := detectCategories(intent.TrainingPhrases)
		allCategories := []string{"direct_question", "first_person", "conversational", "abbreviated", "problem_statement", "procedural_request"}
		missing := []string{}
		for _, cat := range allCategories {
			if !contains(categories, cat) {
				missing = append(missing, cat)
			}
		}

		report.Intents = append(report.Intents, IntentExpansionReport{
			IntentID:               string(intent.ID),
			CurrentPhraseCount:     count,
			RecommendedPhraseCount: recommended,
			MissingCoverage:        missing,
			Priority:               priority,
		})
	}
	return report, nil
}

func detectCategories(phrases []string) []string {
	detected := make(map[string]bool)
	for _, p := range phrases {
		p = strings.ToLower(p)
		if strings.HasPrefix(p, "¿") || strings.Contains(p, "qué") || strings.Contains(p, "cuál") || strings.Contains(p, "cómo") {
			detected["direct_question"] = true
		}
		if strings.Contains(p, "yo") || strings.Contains(p, "mi") || strings.Contains(p, "tengo") {
			detected["first_person"] = true
		}
		if len(strings.Fields(p)) > 5 {
			detected["conversational"] = true
		}
		if len(strings.Fields(p)) < 3 {
			detected["abbreviated"] = true
		}
		if strings.Contains(p, "problema") || strings.Contains(p, "no puedo") {
			detected["problem_statement"] = true
		}
		if strings.Contains(p, "inscribirme") || strings.Contains(p, "quiero") {
			detected["procedural_request"] = true
		}
	}
	res := []string{}
	for cat := range detected {
		res = append(res, cat)
	}
	return res
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
