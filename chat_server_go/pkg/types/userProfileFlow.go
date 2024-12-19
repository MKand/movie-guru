package types

type UserProfileFlowInput struct {
	Query        string `json:"query"`
	AgentMessage string `json:"agentMessage"`
}

type UserProfileFlowOutput struct {
	ProfileChangeRecommendations []*ProfileChangeRecommendation `json:"profileChangeRecommendations"`
	*ModelOutputMetadata         `json:"modelOutputMetadata"`
}

type ProfileChangeRecommendation struct {
	Item     string               `json:"item"`
	Reason   string               `json:"reason"`
	Category MovieFeatureCategory `json:"category"`
	Sentiment
}

type UserProfileOutput struct {
	UserProfile *UserProfile `json:"userProfile"`
	*ModelOutputMetadata
}

func NewUserProfileFlowOuput() *UserProfileFlowOutput {
	return &UserProfileFlowOutput{
		ProfileChangeRecommendations: make([]*ProfileChangeRecommendation, 5),
		ModelOutputMetadata: &ModelOutputMetadata{
			Justification: "",
			SafetyIssue:   false,
		},
	}
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

type UserProfile struct {
	Likes    ProfileCategories `json:"likes, omitempty"`
	Dislikes ProfileCategories `json:"dislikes, omitempty"`
}
type ProfileCategories struct {
	Actors    []string `json:"actors, omitempty"`
	Directors []string `json:"directors, omitempty"`
	Genres    []string `json:"genres, omitempty"`
	Others    []string `json:"others, omitempty"`
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
