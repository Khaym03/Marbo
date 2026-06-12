// Package audit provides functions to audit knowledge base flows and transitions.
package audit

import (
	"fmt"

	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/embedder"
	"github.com/Khaym03/Marbo/internal/similarity"
)

type TransitionAuditor struct {
	Embedder embedder.Embedder
}

func (a *TransitionAuditor) Audit(kb *domain.KnowledgeBase) (*AuditReport, error) {
	report := &AuditReport{}

	for _, flow := range kb.Flows {
		flowAudit := FlowAudit{
			FlowID: string(flow.ID),
		}

		for _, node := range flow.Nodes {
			flowAudit.TransitionCount += len(node.Transitions)

			// Rule 5: Detect nodes with only one outgoing transition and not terminal.
			if len(node.Transitions) == 1 && !node.IsTerminal {
				flowAudit.Warnings = append(flowAudit.Warnings, fmt.Sprintf("Node %s has only one outgoing transition and is not terminal.", node.ID))
			}

			// Pre-embed phrases and organize them
			type nodePhrase struct {
				phrase       string
				transitionID domain.NodeID
				vector       []float32
			}
			phrases := []nodePhrase{}

			for _, trans := range node.Transitions {
				if len(trans.TrainingPhrases) == 0 {
					flowAudit.Warnings = append(flowAudit.Warnings, fmt.Sprintf("Transition %s -> %s has no phrases.", node.ID, trans.TargetNode))
					continue
				}

				for _, p := range trans.TrainingPhrases {
					vec, err := a.Embedder.Embed(p)
					if err != nil {
						return nil, err
					}
					phrases = append(phrases, nodePhrase{p, trans.TargetNode, vec})
				}
			}

			// Rule 1: Duplicate phrases in same node
			for i := 0; i < len(phrases); i++ {
				for j := i + 1; j < len(phrases); j++ {
					if phrases[i].phrase == phrases[j].phrase {
						flowAudit.DuplicatePhrases = append(flowAudit.DuplicatePhrases, DuplicatePhrase{
							Phrase:      phrases[i].phrase,
							TransitionA: phrases[i].transitionID,
							TransitionB: phrases[j].transitionID,
						})
					}
				}
			}

			// Rule 2 & 3: Similar phrases
			overlapCount := 0
			for i := 0; i < len(phrases); i++ {
				for j := i + 1; j < len(phrases); j++ {
					score, err := similarity.DotProduct(phrases[i].vector, phrases[j].vector)
					if err != nil {
						return nil, err
					}

					// Rule 2: High similarity
					if score >= 0.92 {
						flowAudit.SimilarTransitions = append(flowAudit.SimilarTransitions, SimilarTransition{
							PhraseA:     phrases[i].phrase,
							PhraseB:     phrases[j].phrase,
							Score:       score,
							TransitionA: phrases[i].transitionID,
							TransitionB: phrases[j].transitionID,
						})
					}

					// Rule 3: Weak branch separation (different branches, score >= 0.90)
					if phrases[i].transitionID != phrases[j].transitionID && score >= 0.90 {
						overlapCount++
					}
				}
			}

			if overlapCount > 0 {
				flowAudit.Warnings = append(flowAudit.Warnings, fmt.Sprintf("Node %s may contain semantically overlapping branches.", node.ID))
			}
		}

		if len(flowAudit.Warnings) == 0 && len(flowAudit.DuplicatePhrases) == 0 && len(flowAudit.SimilarTransitions) == 0 {
			// Optional: add a "None" marker if preferred, but let's just leave it empty.
		}

		report.Flows = append(report.Flows, flowAudit)
	}

	return report, nil
}
