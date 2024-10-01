package types

type MovieFlowInput struct {
	History          []*SimpleMessage `json:"history"`
	UserPreferences  *UserProfile     `json:"userPreferences"`
	ContextDocuments []*MovieContext  `json:"contextDocuments"`
	UserMessage      string           `json:"userMessage"`
}

type MovieFlowOutput struct {
	Answer               string           `json:"answer"`
	RelevantMoviesTitles []*RelevantMovie `json:"relevantMovies"`
	WrongQuery           bool             `json:"wrongQuery,omitempty" `
	*ModelOutputMetadata
}

type RelevantMovie struct {
	Title  string `json:"title"`
	Reason string `json:"reason"`
}
