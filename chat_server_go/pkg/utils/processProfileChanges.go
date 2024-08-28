package utils

import (
	"strings"

	types "github.com/movie-guru/pkg/types"
)

func ProcessProfileChanges(userProfile *types.UserProfile, changes []*types.ProfileChangeRecommendation) (*types.UserProfile, error) {
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
			if !Contains(userProfile.Likes.Actors, change.Item) {
				userProfile.Likes.Actors = append(userProfile.Likes.Actors, change.Item)
			}
			userProfile.Dislikes.Actors = RemoveItem(userProfile.Dislikes.Actors, change.Item)
		}
		if change.Sentiment == types.NEGATIVE {
			if !Contains(userProfile.Dislikes.Actors, change.Item) {
				userProfile.Dislikes.Actors = append(userProfile.Dislikes.Actors, change.Item)
			}
			userProfile.Likes.Actors = RemoveItem(userProfile.Likes.Actors, change.Item)
		}
	}
	if category == "DIRECTOR" {
		if change.Sentiment == types.POSITIVE {
			if !Contains(userProfile.Likes.Directors, change.Item) {
				userProfile.Likes.Directors = append(userProfile.Likes.Directors, change.Item)
			}
			userProfile.Dislikes.Directors = RemoveItem(userProfile.Dislikes.Directors, change.Item)
		}
		if change.Sentiment == types.NEGATIVE {
			if !Contains(userProfile.Dislikes.Directors, change.Item) {
				userProfile.Dislikes.Directors = append(userProfile.Dislikes.Directors, change.Item)
			}
			userProfile.Likes.Directors = RemoveItem(userProfile.Likes.Directors, change.Item)
		}
	}
	if category == "GENRE" {
		if change.Sentiment == types.POSITIVE {
			if !Contains(userProfile.Likes.Genres, change.Item) {
				userProfile.Likes.Genres = append(userProfile.Likes.Genres, change.Item)
			}
			userProfile.Dislikes.Genres = RemoveItem(userProfile.Dislikes.Genres, change.Item)
		}

		if change.Sentiment == types.NEGATIVE {
			if !Contains(userProfile.Dislikes.Genres, change.Item) {
				userProfile.Dislikes.Genres = append(userProfile.Dislikes.Genres, change.Item)
			}
			userProfile.Likes.Genres = RemoveItem(userProfile.Likes.Genres, change.Item)
		}
	}
	if category == "OTHER" {
		if change.Sentiment == types.POSITIVE {
			if !Contains(userProfile.Likes.Others, change.Item) {
				userProfile.Likes.Others = append(userProfile.Likes.Others, change.Item)
			}
			userProfile.Dislikes.Others = RemoveItem(userProfile.Dislikes.Others, change.Item)
		}
		if change.Sentiment == types.NEGATIVE {
			if !Contains(userProfile.Dislikes.Others, change.Item) {
				userProfile.Dislikes.Others = append(userProfile.Dislikes.Others, change.Item)
			}
			userProfile.Likes.Others = RemoveItem(userProfile.Likes.Others, change.Item)
		}
	}
}
