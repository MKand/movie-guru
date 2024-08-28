package utils

import (
	types "github.com/movie-guru/pkg/types"
)

func FilterRelevantContext(relevantMovies []string, fullContext []*types.MovieContext) []*types.MovieContext {
	relevantContext := make(
		[]*types.MovieContext,
		0,
		len(relevantMovies),
	)
	for _, m := range fullContext {
		for _, r := range relevantMovies {
			if r == m.Title {
				if m.Poster != "" {
					relevantContext = append(relevantContext, m)
				}
			}
		}
	}
	return relevantContext
}
