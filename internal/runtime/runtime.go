// Package runtime provides the main execution engine for the conversation system, handling knowledge base caching, intent resolution, flow routing, and conversation state management.
package runtime

import (
	"encoding/json"
	"fmt"

	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/embedder"
)

type RuntimeResult struct {
	Type RuntimeResultType

	Response domain.Response

	IntentID domain.IntentID
	ZoneID   domain.ZoneID

	FlowID domain.FlowID
	NodeID domain.NodeID

	Extension *RuntimeExtension
}

func (r *RuntimeResult) String() string {
	data, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return "{}"
	}
	return string(data)
}

type RuntimeResultType string

const (
	ResultAnswer        RuntimeResultType = "answer"
	ResultFlow          RuntimeResultType = "flow"
	ResultFallback      RuntimeResultType = "fallback"
	ResultClarification RuntimeResultType = "clarification"
)

type RuntimeExtension struct {
	Confidence *ConfidenceData
	Clarify    *ClarificationData
	Trace      *TraceData
}

type ConfidenceData struct {
	Score float32
}

type ClarificationData struct {
	Candidates []domain.IntentID
}

// Placeholders for future phases
type TraceData struct{}

type ConversationState struct {
	ActiveFlowID  domain.FlowID
	CurrentNodeID domain.NodeID
}

type Runtime struct {
	embedder       embedder.Embedder
	intentResolver *IntentResolver
	flowRouter     *FlowRouter
	cache          *Cache
	state          *ConversationState
	settings       domain.Settings
}

func NewRuntime(embedder embedder.Embedder, cache *Cache, settings domain.Settings) *Runtime {
	return &Runtime{
		embedder:       embedder,
		intentResolver: NewIntentResolver(cache, settings),
		flowRouter:     NewFlowRouter(cache, settings),
		cache:          cache,
		settings:       settings,
	}
}

func (r *Runtime) Handle(text string) (RuntimeResult, error) {
	PrintStateDebug(r, "BEFORE HANDLE")
	// 1. Embed
	vector, err := r.embedder.Embed(text)
	if err != nil {
		return RuntimeResult{}, err
	}

	// 2. Handle active flow
	if r.state != nil && r.state.ActiveFlowID != "" {
		fmt.Println("FLOW PATH SELECTED")
		step, ok := r.flowRouter.Step(r.state.ActiveFlowID, r.state.CurrentNodeID, vector)
		if !ok {
			return RuntimeResult{}, fmt.Errorf("failed to step flow: %s", r.state.ActiveFlowID)
		}

		result := RuntimeResult{
			Type:     ResultFlow,
			Response: step.Response,
			FlowID:   step.FlowID,
			NodeID:   step.NodeID,
		}

		if step.IsTerminal {
			r.state = nil
		} else {
			r.state.CurrentNodeID = step.NodeID
		}
		PrintStateDebug(r, "AFTER HANDLE")

		return result, nil
	}

	// 3. Resolve Intent (only if no active flow)
	fmt.Println("INTENT RESOLUTION PATH SELECTED")
	ranking, err := r.intentResolver.ResolveTopK(vector)
	if err != nil {
		return RuntimeResult{}, err
	}

	decision := EvaluateConfidence(ranking, r.settings)

	debug := BuildConfidenceDebug(text, ranking.Candidates, r.settings, decision)
	PrintConfidenceDebug(debug)

	switch decision.Type {
	case DecisionFallback:
		PrintStateDebug(r, "AFTER HANDLE")
		return RuntimeResult{Type: ResultFallback}, nil
	case DecisionClarification:
		PrintStateDebug(r, "AFTER HANDLE")
		return RuntimeResult{
			Type: ResultClarification,
			Extension: &RuntimeExtension{
				Confidence: &ConfidenceData{Score: decision.Confidence},
				Clarify:    &ClarificationData{Candidates: decision.Candidates},
			},
		}, nil
	case DecisionExecute:
		// Proceed with execution
		intentID := decision.IntentID
		intentDef, ok := r.cache.IntentMap[intentID]
		if !ok {
			return RuntimeResult{}, fmt.Errorf("intent definition not found: %s", intentID)
		}

		extension := &RuntimeExtension{
			Confidence: &ConfidenceData{Score: decision.Confidence},
		}

		if !intentDef.RequiresFlow {
			PrintStateDebug(r, "AFTER HANDLE")
			return RuntimeResult{
				Type:      ResultAnswer,
				IntentID:  intentDef.ID,
				ZoneID:    intentDef.ZoneID,
				Response:  intentDef.Response,
				Extension: extension,
			}, nil
		}

		step, ok := r.flowRouter.StartFlow(intentDef.FlowID)
		if !ok {
			return RuntimeResult{}, fmt.Errorf("failed to start flow: %s", intentDef.FlowID)
		}

		r.state = &ConversationState{
			ActiveFlowID:  step.FlowID,
			CurrentNodeID: step.NodeID,
		}
		PrintStateDebug(r, "AFTER HANDLE")

		return RuntimeResult{
			Type:      ResultFlow,
			IntentID:  intentDef.ID,
			ZoneID:    intentDef.ZoneID,
			Response:  step.Response,
			FlowID:    step.FlowID,
			NodeID:    step.NodeID,
			Extension: extension,
		}, nil
	}
	PrintStateDebug(r, "AFTER HANDLE")
	return RuntimeResult{Type: ResultFallback}, nil
}
