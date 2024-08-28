package standaloneWeb

import (
	db "github.com/movie-guru/pkg/db"
	standaloneWrappers "github.com/movie-guru/pkg/standaloneWrappers"
)

type Dependencies struct {
	QueryTransformAgent *standaloneWrappers.QueryTransformAgent
	PrefAgent           *standaloneWrappers.ProfileAgent
	MovieAgent          *standaloneWrappers.MovieAgent
	Retriever           *standaloneWrappers.MovieRetriever
	DB                  *db.MovieAgentDB
}
