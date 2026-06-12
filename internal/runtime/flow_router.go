package runtime

import (
	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/similarity"
)

type FlowStep struct {
	FlowID     domain.FlowID
	NodeID     domain.NodeID
	Response   domain.Response
	IsTerminal bool
}
type FlowRouter struct {
	cache    *Cache
	settings domain.Settings
}

func NewFlowRouter(cache *Cache, settings domain.Settings) *FlowRouter {
	return &FlowRouter{
		cache:    cache,
		settings: settings,
	}
}

func (r *FlowRouter) getFlow(flowID domain.FlowID) *domain.Flow {
	for i := range r.cache.Flows {
		if r.cache.Flows[i].ID == flowID {
			return &r.cache.Flows[i]
		}
	}
	return nil
}

func getNode(flow *domain.Flow, nodeID domain.NodeID) *domain.FlowNode {
	for i := range flow.Nodes {
		if flow.Nodes[i].ID == nodeID {
			return &flow.Nodes[i]
		}
	}
	return nil
}

func (r *FlowRouter) StartFlow(flowID domain.FlowID) (*FlowStep, bool) {
	flow := r.getFlow(flowID)
	if flow == nil {
		return nil, false
	}

	node := getNode(flow, flow.StartNode)
	if node == nil {
		return nil, false
	}

	return &FlowStep{
		FlowID:     flow.ID,
		NodeID:     node.ID,
		Response:   node.Response,
		IsTerminal: node.IsTerminal,
	}, true
}

func (r *FlowRouter) Step(flowID domain.FlowID, currentNodeID domain.NodeID, queryEmbedding []float32) (*FlowStep, bool) {

	flow := r.getFlow(flowID)
	if flow == nil {
		return nil, false
	}

	node := getNode(flow, currentNodeID)
	if node == nil {
		return nil, false
	}

	// 1. Filter transition vectors for current node
	var bestTarget domain.NodeID
	var bestScore float32 = -1.0

	// 2. Score transitions
	for _, tv := range r.cache.Transitions {
		if tv.FlowID == flowID && tv.NodeID == currentNodeID {
			score, err := similarity.DotProduct(queryEmbedding, tv.Vector)
			if err != nil {
				continue
			}

			if score > bestScore {
				bestScore = score
				bestTarget = tv.TargetNode
			}
		}
	}

	// 3. Threshold check
	if bestScore < r.settings.SimilarityThreshold {
		// Fallback: stay in same node
		return &FlowStep{
			FlowID:     flow.ID,
			NodeID:     node.ID,
			Response:   node.Response,
			IsTerminal: node.IsTerminal,
		}, true
	}

	// 4. Move to next node
	next := getNode(flow, bestTarget)
	if next == nil {
		return nil, false
	}

	return &FlowStep{
		FlowID:     flow.ID,
		NodeID:     next.ID,
		Response:   next.Response,
		IsTerminal: next.IsTerminal,
	}, true
}
