// Package metrics provides functions to analyze and report on knowledge base metrics.
package metrics

import "github.com/Khaym03/Marbo/internal/domain"

func AnalyzeFlows(kb *domain.KnowledgeBase) ([]FlowMetrics, error) {
	metrics := []FlowMetrics{}
	for _, flow := range kb.Flows {
		nodeCount := len(flow.Nodes)
		transCount := 0
		for _, n := range flow.Nodes {
			transCount += len(n.Transitions)
		}

		maxDepth := computeMaxDepth(flow)

		complexity := float64((nodeCount * 2) + transCount + maxDepth)

		// Risk heuristic: low nodes/transitions/depth + potential for being broken
		risk := 0.0
		if nodeCount < 3 {
			risk += 50
		}
		if transCount < 2 {
			risk += 30
		}
		if risk > 100 {
			risk = 100
		}

		metrics = append(metrics, FlowMetrics{
			FlowID:          string(flow.ID),
			NodeCount:       nodeCount,
			TransitionCount: transCount,
			MaxDepth:        maxDepth,
			ComplexityScore: complexity,
			RiskScore:       risk,
		})
	}
	return metrics, nil
}

func computeMaxDepth(flow domain.Flow) int {
	// Simple BFS-like approach or just fixed-point
	depths := make(map[domain.NodeID]int)

	// Initialize
	for _, node := range flow.Nodes {
		depths[node.ID] = 1
	}

	changed := true
	for changed {
		changed = false
		for _, node := range flow.Nodes {
			for _, trans := range node.Transitions {
				if depths[trans.TargetNode] < depths[node.ID]+1 {
					depths[trans.TargetNode] = depths[node.ID] + 1
					changed = true
				}
			}
		}
	}

	maxD := 0
	for _, d := range depths {
		if d > maxD {
			maxD = d
		}
	}
	return maxD
}
