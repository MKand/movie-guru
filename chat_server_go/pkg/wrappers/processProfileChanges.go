package wrappers

import (
	"strings"

	types "github.com/movie-guru/pkg/types"
	"github.com/movie-guru/pkg/utils"
)

func processProfileChanges(userProfile *types.UserProfile, changes []*types.ProfileChangeRecommendation) (*types.UserProfile, error) {
	for _, change := range changes {
		change.Item = strings.ToLower(change.Item)

		switch change.Category {
		case "ACTOR":
			handleChange("ACTOR", userProfile, change)
		case "DIRECTOR":
			handleChange("DIRECTOR", userProfile, change)
		case "GENRE":
			handleChange("GENRE", userProfile, change)
		case "OTHER":
			handleChange("OTHER", userProfile, change)
		}
	}
	return userProfile, nil
}

func handleChange(category string, userProfile *types.UserProfile, change *types.ProfileChangeRecommendation) {
	if category == "ACTOR" {
		if change.Sentiment == types.POSITIVE {
			if !utils.Contains(userProfile.Likes.Actors, change.Item) {
				userProfile.Likes.Actors = append(userProfile.Likes.Actors, change.Item)
			}
			userProfile.Dislikes.Actors = utils.RemoveItem(userProfile.Dislikes.Actors, change.Item)
		}
		if change.Sentiment == types.NEGATIVE {
			if !utils.Contains(userProfile.Dislikes.Actors, change.Item) {
				userProfile.Dislikes.Actors = append(userProfile.Dislikes.Actors, change.Item)
			}
			userProfile.Likes.Actors = utils.RemoveItem(userProfile.Likes.Actors, change.Item)
		}
	}
	if category == "DIRECTOR" {
		if change.Sentiment == types.POSITIVE {
			if !utils.Contains(userProfile.Likes.Directors, change.Item) {
				userProfile.Likes.Directors = append(userProfile.Likes.Directors, change.Item)
			}
			userProfile.Dislikes.Directors = utils.RemoveItem(userProfile.Dislikes.Directors, change.Item)
		}
		if change.Sentiment == types.NEGATIVE {
			if !utils.Contains(userProfile.Dislikes.Directors, change.Item) {
				userProfile.Dislikes.Directors = append(userProfile.Dislikes.Directors, change.Item)
			}
			userProfile.Likes.Directors = utils.RemoveItem(userProfile.Likes.Directors, change.Item)
		}
	}
	if category == "GENRE" {
		if change.Sentiment == types.POSITIVE {
			if !utils.Contains(userProfile.Likes.Genres, change.Item) {
				userProfile.Likes.Genres = append(userProfile.Likes.Genres, change.Item)
			}
			userProfile.Dislikes.Genres = utils.RemoveItem(userProfile.Dislikes.Genres, change.Item)
		}

		if change.Sentiment == types.NEGATIVE {
			if !utils.Contains(userProfile.Dislikes.Genres, change.Item) {
				userProfile.Dislikes.Genres = append(userProfile.Dislikes.Genres, change.Item)
			}
			userProfile.Likes.Genres = utils.RemoveItem(userProfile.Likes.Genres, change.Item)
		}
	}
	if category == "OTHER" {
		if change.Sentiment == types.POSITIVE {
			if !utils.Contains(userProfile.Likes.Others, change.Item) {
				userProfile.Likes.Others = append(userProfile.Likes.Others, change.Item)
			}
			userProfile.Dislikes.Others = utils.RemoveItem(userProfile.Dislikes.Others, change.Item)
		}
		if change.Sentiment == types.NEGATIVE {
			if !utils.Contains(userProfile.Dislikes.Others, change.Item) {
				userProfile.Dislikes.Others = append(userProfile.Dislikes.Others, change.Item)
			}
			userProfile.Likes.Others = utils.RemoveItem(userProfile.Likes.Others, change.Item)
		}
	}
}
