package types

type SearchType string

const (
	Vector  SearchType = "Vector"
	Keyword SearchType = "Keyword"
	Mixed   SearchType = "Mixed"
)

type UserQuery struct {
	VectorQuery    string     `json:"vectorQuery"`
	KeywordQuery   string     `json:"keywordQuery"`
	Total          int        `json:"totalDocuments"`
	SearchCategory SearchType `json:"searchCategory"`
}

type MovieContextSchema struct {
	Title          string   `json:"title"`
	RuntimeMinutes int      `json:"runtime_minutes"`
	Genres         []string `json:"genres"`
	Rating         float32  `json:"rating"`
	Released       int      `json:"released"`
	Director       string   `json:"director"`
	Actors         []string `json:"actors"`
	Poster         string   `json:"poster"`
	Tconst         string   `json:"tconst"`
}