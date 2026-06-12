package expansion

import (
	"fmt"
	"time"

	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/planner"
)

func GenerateExpansionPack(kb *domain.KnowledgeBase, report *planner.ExpansionReport) *ExpansionPack {
	pack := &ExpansionPack{
		GeneratedAt: time.Now(),
	}

	seeds := GetSeeds()

	templates := map[string][]string{
		"direct_question":    {"cómo hago para %s", "¿qué necesito para %s?", "¿cuál es el procedimiento para %s?"},
		"first_person":       {"yo quiero %s", "necesito %s", "tengo dudas sobre %s"},
		"conversational":     {"ayuda con %s", "por favor explícame sobre %s", "hablemos sobre %s"},
		"abbreviated":        {"%s breve", "%s rápido", "%s info"},
		"problem_statement":  {"tengo problemas con %s", "no entiendo %s", "ayuda con mi caso de %s"},
		"procedural_request": {"quiero realizar %s", "solicitud de %s", "proceder con %s"},
	}

	for _, intentReport := range report.Intents {
		if intentReport.Priority == planner.PriorityLow {
			continue
		}

		var intentSeeds []string
		for _, seedSet := range seeds {
			if string(seedSet.IntentID) == intentReport.IntentID {
				intentSeeds = seedSet.Concepts
				break
			}
		}

		// Fallback if no seeds found, but skip generation
		if len(intentSeeds) == 0 {
			continue
		}

		generated := []string{}
		for _, missing := range intentReport.MissingCoverage {
			if tmpl, ok := templates[missing]; ok {
				// Deterministically pick a concept from the seed set for each template
				// Simple approach: rotate through concepts
				for idx, t := range tmpl {
					concept := intentSeeds[idx%len(intentSeeds)]
					generated = append(generated, fmt.Sprintf(t, concept))
				}
			}
		}

		pack.Intents = append(pack.Intents, IntentExpansionPack{
			IntentID:            intentReport.IntentID,
			ExistingPhraseCount: intentReport.CurrentPhraseCount,
			TargetPhraseCount:   intentReport.RecommendedPhraseCount,
			GeneratedPhrases:    generated,
		})
	}
	return pack
}
