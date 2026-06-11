package domain

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(path string) (*KnowledgeBase, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var kb KnowledgeBase

	if err := json.Unmarshal(data, &kb); err != nil {
		return nil, err
	}

	// Instrumentation
	fmt.Printf("--- LOADER DIAGNOSTICS ---\n")
	fmt.Printf("Number of flows loaded: %d\n", len(kb.Flows))
	for _, flow := range kb.Flows {
		fmt.Printf("\nFLOW: %s\n", flow.ID)
		fmt.Printf("StartNode=%s\n", flow.StartNode)
		fmt.Printf("Nodes=%d\n", len(flow.Nodes))
		for _, node := range flow.Nodes {
			fmt.Printf("\nNODE: %s\n", node.ID)
			fmt.Printf("Transitions=%d\n", len(node.Transitions))
		}
	}
	fmt.Printf("\n--------------------------\n")

	// Validation happens later in the pipeline
	return &kb, nil
}
