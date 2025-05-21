package types

// Metadata stores application metadata
type Metadata struct {
	TokenAudience        string
	HistoryLength        int
	MaxUserMessageLength int
	CorsOrigins          string
	FlowsURL             string
	FeedbackURL          string
	UseAuth              bool
	StrictCors           bool
	EnableMetrics        bool
}
