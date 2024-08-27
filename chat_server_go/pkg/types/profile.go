package types

type ProfileAgentInput struct {
	Query        string `json:"query"`
	AgentMessage string `json:"agentMessage"`
}

type MovieFeatureCategory string
type ProfileAction string
type Sentiment string

const (
	OTHER    MovieFeatureCategory = "OTHER"
	ACTOR    MovieFeatureCategory = "ACTOR"
	DIRECTOR MovieFeatureCategory = "DIRECTOR"
	GENRE    MovieFeatureCategory = "GENRE"
)

const (
	UNSURE ProfileAction = "UNSURE"
	ADD    ProfileAction = "ADD"
	REMOVE ProfileAction = "REMOVE"
)

const (
	POSITIVE Sentiment = "POSITIVE"
	NEGATIVE Sentiment = "NEGATIVE"
)

type ProfileChangeRecommendation struct {
	Item     string               `json:"item"`
	Reason   string               `json:"reason"`
	Category MovieFeatureCategory `json:"category"`
	Sentiment
}

type UserProfileAgentOutput struct {
	ProfileChangeRecommendations []*ProfileChangeRecommendation `json:"profileChangeRecommendations"`
	ChangesMade                  bool                           `json:"changesMade,omitempty"`
	*ModelOutputMetadata
}

func NewUserProfileAgentOuput() *UserProfileAgentOutput {
	return &UserProfileAgentOutput{
		ProfileChangeRecommendations: make([]*ProfileChangeRecommendation, 5),
		ChangesMade:                  false,
		ModelOutputMetadata: &ModelOutputMetadata{
			Justification: "",
			SafetyIssue:   false,
		},
	}
}

type UserProfileOutput struct {
	UserProfile *UserProfile `json:"userProfile"`
	ChangesMade bool
	*ModelOutputMetadata
}

type UserProfile struct {
	Likes    ProfileCategories `json:"likes, omitempty"`
	Dislikes ProfileCategories `json:"dislikes, omitempty"`
}
type ProfileCategories struct {
	Actors    []string `json:"actors, omitempty"`
	Directors []string `json:"director, omitempty"`
	Genres    []string `json:"genres, omitempty"`
	Others    []string `json:"other, omitempty"`
}

func NewUserProfile() *UserProfile {
	return &UserProfile{
		Likes: ProfileCategories{
			Actors:    make([]string, 0),
			Directors: make([]string, 0),
			Genres:    make([]string, 0),
			Others:    make([]string, 0),
		},
		Dislikes: ProfileCategories{
			Actors:    make([]string, 0),
			Directors: make([]string, 0),
			Genres:    make([]string, 0),
			Others:    make([]string, 0),
		},
	}
}
