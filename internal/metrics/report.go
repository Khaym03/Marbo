package metrics

type KBHealthReport struct {
	OverallScore        float64
	IntentMetrics       []IntentMetrics
	FlowMetrics         []FlowMetrics
	HighestRiskIntent   string
	HighestRiskFlow     string
	Warnings            []string
}

type IntentMetrics struct {
	IntentID            string
	TrainingPhraseCount int
	AveragePhraseLength float64
	CoverageScore       float64
	RiskScore           float64
}

type FlowMetrics struct {
	FlowID          string
	NodeCount       int
	TransitionCount int
	MaxDepth        int
	ComplexityScore float64
	RiskScore       float64
}
