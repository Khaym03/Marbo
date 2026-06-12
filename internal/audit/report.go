package audit

import "github.com/Khaym03/Marbo/internal/domain"

type AuditReport struct {
	Flows []FlowAudit
}

type FlowAudit struct {
	FlowID             string
	TransitionCount    int
	DuplicatePhrases   []DuplicatePhrase
	SimilarTransitions []SimilarTransition
	Warnings           []string
}

type DuplicatePhrase struct {
	Phrase      string
	TransitionA domain.NodeID
	TransitionB domain.NodeID
}

type SimilarTransition struct {
	PhraseA     string
	PhraseB     string
	Score       float32
	TransitionA domain.NodeID
	TransitionB domain.NodeID
}
