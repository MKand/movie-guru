package web

import (
	"github.com/movie-guru/pkg/db"
	"github.com/movie-guru/pkg/types"
	"github.com/movie-guru/pkg/wrappers"
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
	QueryTransformAgent *wrappers.QueryTransformAgent
	PrefAgent           *wrappers.ProfileAgent
	MovieAgent          *wrappers.MovieAgent
	Retriever           *wrappers.MovieRetriever
	DB                  *db.MovieAgentDB
}
