// Package validator provides functions to validate the structure and integrity of a knowledge base.
package validator

import (
	"fmt"
	"strings"

	"github.com/Khaym03/Marbo/internal/domain"
)

func Validate(kb *domain.KnowledgeBase) error {
	var errs []ValidationError

	// Rule 1: Unique Zone IDS
	zoneIDs := make(map[domain.ZoneID]bool)
	for _, zone := range kb.Zones {
		if zoneIDs[zone.ID] {
			errs = append(errs, ValidationError{"DUPLICATE_ZONE_ID", fmt.Sprintf("zone %s already exists", zone.ID)})
		}
		zoneIDs[zone.ID] = true
	}

	// Rule 2: Unique Intent IDS
	intentIDs := make(map[domain.IntentID]bool)
	for _, intent := range kb.Intents {
		if intentIDs[intent.ID] {
			errs = append(errs, ValidationError{"DUPLICATE_INTENT_ID", fmt.Sprintf("intent %s already exists", intent.ID)})
		}
		intentIDs[intent.ID] = true
	}

	// Rule 3: Unique Flow IDS
	flowIDs := make(map[domain.FlowID]bool)
	for _, flow := range kb.Flows {
		if flowIDs[flow.ID] {
			errs = append(errs, ValidationError{"DUPLICATE_FLOW_ID", fmt.Sprintf("flow %s already exists", flow.ID)})
		}
		flowIDs[flow.ID] = true
	}

	// Rule 4, 5, 6, 10
	for _, intent := range kb.Intents {
		// Rule 4: Valid Intent Zone
		if !zoneIDs[intent.ZoneID] {
			errs = append(errs, ValidationError{"INVALID_ZONE_REFERENCE", fmt.Sprintf("intent %s references unknown zone %s", intent.ID, intent.ZoneID)})
		}

		// Rule 5: Flow Reference
		if intent.RequiresFlow {
			if intent.FlowID == "" || !flowIDs[intent.FlowID] {
				errs = append(errs, ValidationError{"MISSING_FLOW_REFERENCE", fmt.Sprintf("intent %s references missing or empty flow %s", intent.ID, intent.FlowID)})
			}
		}

		// Rule 6: Static Response
		if !intent.RequiresFlow && intent.Response.Text == "" {
			errs = append(errs, ValidationError{"INVALID_RESPONSE", fmt.Sprintf("intent %s requires static response but text is empty", intent.ID)})
		}

		// Rule 10: Training Phrases
		if len(intent.TrainingPhrases) == 0 {
			errs = append(errs, ValidationError{"EMPTY_TRAINING_PHRASES", fmt.Sprintf("intent %s has no training phrases", intent.ID)})
		}
		for _, phrase := range intent.TrainingPhrases {
			if strings.TrimSpace(phrase) == "" {
				errs = append(errs, ValidationError{"INVALID_TRAINING_PHRASE", fmt.Sprintf("intent %s has an empty training phrase", intent.ID)})
			}
		}

		// Rule 12: Intent Label
		if strings.TrimSpace(intent.Label) == "" {
			errs = append(errs, ValidationError{"MISSING_INTENT_LABEL", fmt.Sprintf("intent %s has missing or empty label", intent.ID)})
		}
	}

	// Rule 7, 8, 9, 11
	for _, flow := range kb.Flows {
		nodeIDs := make(map[domain.NodeID]bool)
		fmt.Printf("\nVALIDATING FLOW: %s\n", flow.ID)
		fmt.Printf("START NODE: %s\n", flow.StartNode)
		fmt.Printf("AVAILABLE NODES:\n")
		for _, node := range flow.Nodes {
			fmt.Printf("* %s\n", node.ID)
			// Rule 8: Unique Node IDS
			if nodeIDs[node.ID] {
				errs = append(errs, ValidationError{"DUPLICATE_NODE_ID", fmt.Sprintf("flow %s has duplicate node %s", flow.ID, node.ID)})
			}
			nodeIDs[node.ID] = true
		}

		// Rule 7: Flow Start Node exists
		if _, ok := nodeIDs[flow.StartNode]; !ok {
			errs = append(errs, ValidationError{"INVALID_START_NODE", fmt.Sprintf("flow %s has invalid start node %s", flow.ID, flow.StartNode)})
		}

		// Rule 9: Valid Transition Targets
		for _, node := range flow.Nodes {
			for _, trans := range node.Transitions {
				fmt.Printf("\nCHECKING TRANSITION\nFROM: %s\nTO: %s\n", node.ID, trans.TargetNode)
				fmt.Printf("AVAILABLE TARGETS:\n")
				for n := range nodeIDs {
					fmt.Printf("* %s\n", n)
				}
				if _, ok := nodeIDs[trans.TargetNode]; !ok {
					errs = append(errs, ValidationError{"INVALID_TRANSITION_TARGET", fmt.Sprintf("flow %s has transition to missing node %s", flow.ID, trans.TargetNode)})
				}
			}
		}

		// Rule 11: Reachability
		if _, ok := nodeIDs[flow.StartNode]; ok {
			visited := make(map[domain.NodeID]bool)
			queue := []domain.NodeID{flow.StartNode}
			visited[flow.StartNode] = true
			for len(queue) > 0 {
				currID := queue[0]
				queue = queue[1:]

				var currNode *domain.FlowNode
				for i := range flow.Nodes {
					if flow.Nodes[i].ID == currID {
						currNode = &flow.Nodes[i]
						break
					}
				}
				if currNode == nil {
					continue
				}
				for _, trans := range currNode.Transitions {
					if !visited[trans.TargetNode] {
						visited[trans.TargetNode] = true
						queue = append(queue, trans.TargetNode)
					}
				}
			}
			for _, node := range flow.Nodes {
				if !visited[node.ID] {
					errs = append(errs, ValidationError{"UNREACHABLE_NODE", fmt.Sprintf("flow %s has unreachable node %s", flow.ID, node.ID)})
				}
			}
		}
	}

	if len(errs) > 0 {
		return &ValidationErrors{Errors: errs}
	}
	return nil
}
