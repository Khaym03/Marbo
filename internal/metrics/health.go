package metrics

import (
	"github.com/Khaym03/Marbo/internal/audit"
	"github.com/Khaym03/Marbo/internal/domain"
)

func BuildHealthReport(
	kb *domain.KnowledgeBase,
	auditReport *audit.AuditReport,
) (*KBHealthReport, error) {
	intentMetrics, _ := AnalyzeIntents(kb)
	flowMetrics, _ := AnalyzeFlows(kb)

	avgCoverage := 0.0
	for _, m := range intentMetrics {
		avgCoverage += m.CoverageScore
	}
	if len(intentMetrics) > 0 {
		avgCoverage /= float64(len(intentMetrics))
	}

	avgFlowRisk := 0.0
	for _, m := range flowMetrics {
		avgFlowRisk += m.RiskScore
	}
	if len(flowMetrics) > 0 {
		avgFlowRisk /= float64(len(flowMetrics))
	}

	overall := (avgCoverage + (100 - avgFlowRisk)) / 2.0

	var highestRiskIntent string
	maxIntentRisk := -1.0
	for _, m := range intentMetrics {
		if m.RiskScore > maxIntentRisk {
			maxIntentRisk = m.RiskScore
			highestRiskIntent = m.IntentID
		}
	}

	var highestRiskFlow string
	maxFlowRisk := -1.0
	for _, m := range flowMetrics {
		if m.RiskScore > maxFlowRisk {
			maxFlowRisk = m.RiskScore
			highestRiskFlow = m.FlowID
		}
	}

	return &KBHealthReport{
		OverallScore:      overall,
		IntentMetrics:     intentMetrics,
		FlowMetrics:       flowMetrics,
		HighestRiskIntent: highestRiskIntent,
		HighestRiskFlow:   highestRiskFlow,
	}, nil
}
