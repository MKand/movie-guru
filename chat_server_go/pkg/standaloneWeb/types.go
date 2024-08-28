package standaloneWeb

import (
	db "github.com/movie-guru/pkg/db"
	standaloneWrappers "github.com/movie-guru/pkg/standaloneWrappers"
)

type Dependencies struct {
	QueryTransformFlow *standaloneWrappers.QueryTransformFlow
	UserProfileFlow    *standaloneWrappers.ProfileFlow
	MovieFlow          *standaloneWrappers.MovieFlow
	MovieRetrieverFlow *standaloneWrappers.MovieRetrieverFlow
	DB                 *db.MovieDB
}
