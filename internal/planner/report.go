package planner

type ExpansionPriority string

const (
	PriorityLow    ExpansionPriority = "low"
	PriorityMedium ExpansionPriority = "medium"
	PriorityHigh   ExpansionPriority = "high"
)

type ExpansionReport struct {
	Intents []IntentExpansionReport
}

type IntentExpansionReport struct {
	IntentID               string
	CurrentPhraseCount     int
	RecommendedPhraseCount int
	MissingCoverage        []string
	Priority               ExpansionPriority
}
