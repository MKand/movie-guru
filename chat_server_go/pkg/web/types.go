package web

import (
	"github.com/movie-guru/pkg/db"
	"github.com/movie-guru/pkg/types"
	wrappers "github.com/movie-guru/pkg/wrappers"
)

type SessionInfo struct {
	ID            string
	User          string
	Authenticated bool
}

type LoginBody struct {
	InviteCode string `json:"inviteCode" omitempty`
}

type PrefBody struct {
	Content *types.UserProfile `json:"content"`
}

type ChatRequest struct {
	Content string `json:"content"`
}

type Dependencies struct {
	QueryTransformFlowClient  *wrappers.QueryTransformFlowClient
	UserProfileFlowClient     *wrappers.UserProfileFlowClient
	MovieFlowClient           *wrappers.MovieFlowClient
	MovieRetrieverFlowClient  *wrappers.MovieRetrieverFlowClient
	ResponseQualityFlowClient *wrappers.ResponseQualityFlowClient
	DB                        *db.MovieDB
}
