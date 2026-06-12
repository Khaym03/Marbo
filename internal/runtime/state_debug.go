package runtime

import "fmt"

type StateSnapshot struct {
	HasState      bool
	ActiveFlowID  string
	CurrentNodeID string
}

func (r *Runtime) DebugState() StateSnapshot {
	if r.state == nil {
		return StateSnapshot{HasState: false}
	}
	return StateSnapshot{
		HasState:      true,
		ActiveFlowID:  string(r.state.ActiveFlowID),
		CurrentNodeID: string(r.state.CurrentNodeID),
	}
}

func PrintStateDebug(r *Runtime, prefix string) {
	fmt.Printf("=================================================\n")
	fmt.Printf("STATE %s\n", prefix)
	fmt.Printf("=================================================\n")
	if r.state == nil {
		fmt.Println("Flow:")
		fmt.Println("<none>")
		fmt.Println("\nNode:")
		fmt.Println("<none>")
	} else {
		fmt.Printf("Flow:\n%s\n\nNode:\n%s\n", r.state.ActiveFlowID, r.state.CurrentNodeID)
	}
	fmt.Printf("=================================================\n")
}
