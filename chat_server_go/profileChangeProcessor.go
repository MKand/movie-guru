package main

import (
	"strings"
)

func processProfileChanges(userProfile *UserProfile, changes []*ProfileChangeRecommendation) (*UserProfile, error) {
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

func handleChange(category string, userProfile *UserProfile, change *ProfileChangeRecommendation) {
	if category == "ACTOR" {
		if change.Sentiment == POSITIVE {
			if !contains(userProfile.Likes.Actors, change.Item) {
				userProfile.Likes.Actors = append(userProfile.Likes.Actors, change.Item)
			}
			userProfile.Dislikes.Actors = removeItem(userProfile.Dislikes.Actors, change.Item)
		}
		if change.Sentiment == NEGATIVE {
			if !contains(userProfile.Dislikes.Actors, change.Item) {
				userProfile.Dislikes.Actors = append(userProfile.Dislikes.Actors, change.Item)
			}
			userProfile.Likes.Actors = removeItem(userProfile.Likes.Actors, change.Item)
		}
	}
	if category == "DIRECTOR" {
		if change.Sentiment == POSITIVE {
			if !contains(userProfile.Likes.Directors, change.Item) {
				userProfile.Likes.Directors = append(userProfile.Likes.Directors, change.Item)
			}
			userProfile.Dislikes.Directors = removeItem(userProfile.Dislikes.Directors, change.Item)
		}
		if change.Sentiment == NEGATIVE {
			if !contains(userProfile.Dislikes.Directors, change.Item) {
				userProfile.Dislikes.Directors = append(userProfile.Dislikes.Directors, change.Item)
			}
			userProfile.Likes.Directors = removeItem(userProfile.Likes.Directors, change.Item)
		}
	}
	if category == "GENRE" {
		if change.Sentiment == POSITIVE {
			if !contains(userProfile.Likes.Genres, change.Item) {
				userProfile.Likes.Genres = append(userProfile.Likes.Genres, change.Item)
			}
			userProfile.Dislikes.Genres = removeItem(userProfile.Dislikes.Genres, change.Item)
		}

		if change.Sentiment == NEGATIVE {
			if !contains(userProfile.Dislikes.Genres, change.Item) {
				userProfile.Dislikes.Genres = append(userProfile.Dislikes.Genres, change.Item)
			}
			userProfile.Likes.Genres = removeItem(userProfile.Likes.Genres, change.Item)
		}
	}
	if category == "OTHER" {
		if change.Sentiment == POSITIVE {
			if !contains(userProfile.Likes.Others, change.Item) {
				userProfile.Likes.Others = append(userProfile.Likes.Others, change.Item)
			}
			userProfile.Dislikes.Others = removeItem(userProfile.Dislikes.Others, change.Item)
		}
		if change.Sentiment == NEGATIVE {
			if !contains(userProfile.Dislikes.Others, change.Item) {
				userProfile.Dislikes.Others = append(userProfile.Dislikes.Others, change.Item)
			}
			userProfile.Likes.Others = removeItem(userProfile.Likes.Others, change.Item)
		}
	}
}
