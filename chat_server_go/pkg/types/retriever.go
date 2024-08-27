package types

type MovieContext struct {
	Title          string   `json:"title"`
	RuntimeMinutes int      `json:"runtime_minutes"`
	Genres         []string `json:"genres"`
	Rating         float32  `json:"rating"`
	Plot           string   `json:"plot"`
	Released       int      `json:"released"`
	Director       string   `json:"director"`
	Actors         []string `json:"actors"`
	Poster         string   `json:"poster"`
}
