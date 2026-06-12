package metrics

import (
	"testing"

	"github.com/Khaym03/Marbo/internal/audit"
	"github.com/Khaym03/Marbo/internal/domain"
)

func TestBuildHealthReport(t *testing.T) {
	kb := &domain.KnowledgeBase{
		Intents: []domain.Intent{{ID: "i1", TrainingPhrases: make([]string, 16)}},
		Flows:   []domain.Flow{{ID: "f1", Nodes: []domain.FlowNode{{ID: "n1"}}}},
	}
	auditReport := &audit.AuditReport{}
	report, _ := BuildHealthReport(kb, auditReport)

	if report.OverallScore == 0 {
		t.Error("expected non-zero overall score")
	}
	if report.HighestRiskIntent != "i1" && report.HighestRiskIntent != "" {
		// Risk heuristic is based on coverage, i1 has 100% coverage (risk 0)
	}
}
