package domain

type ZoneID string
type IntentID string
type FlowID string
type NodeID string

type KnowledgeBase struct {
	Version  string   `json:"version"`
	Settings Settings `json:"settings"`
	Zones    []Zone   `json:"zones"`
	Intents  []Intent `json:"intents"`
	Flows    []Flow   `json:"flows"`
	Aliases  []Alias  `json:"aliases"`
}

type Settings struct {
	SimilarityThreshold     float32 `json:"similarity_threshold"`
	AmbiguityThreshold      float32 `json:"ambiguity_threshold"`
	MaxClarificationOptions int     `json:"max_clarification_options"`
	SessionTimeoutMinutes   int     `json:"session_timeout_minutes"`
}

type Zone struct {
	ID          ZoneID     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Map         *MapRegion `json:"map,omitempty"`
}

type MapRegion struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type IntentKind string

const (
	KindFAQ        IntentKind = "faq"
	KindProcedure  IntentKind = "procedure"
	KindDiagnostic IntentKind = "diagnostic"
)

type Intent struct {
	ID              IntentID   `json:"id"`
	Label           string     `json:"label"`
	ZoneID          ZoneID     `json:"zone_id"`
	Kind            IntentKind `json:"kind"`
	RequiresFlow    bool       `json:"requires_flow"`
	FlowID          FlowID     `json:"flow_id,omitempty"`
	Response        Response   `json:"response"`
	TrainingPhrases []string   `json:"training_phrases"`
}

type Flow struct {
	ID        FlowID     `json:"id"`
	Name      string     `json:"name"`
	StartNode NodeID     `json:"start_node_id"`
	Nodes     []FlowNode `json:"nodes"`
}

type FlowNode struct {
	ID          NodeID       `json:"id"`
	Response    Response     `json:"response"`
	IsTerminal  bool         `json:"is_terminal"`
	Transitions []Transition `json:"transitions"`
}

type Transition struct {
	TrainingPhrases []string `json:"phrases"`
	TargetNode      NodeID   `json:"target_node_id"`
}

type Response struct {
	Text        string   `json:"text"`
	ZoneID      ZoneID   `json:"zone_id,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

type Alias struct {
	Canonical string   `json:"canonical"`
	Variants  []string `json:"variants"`
}
