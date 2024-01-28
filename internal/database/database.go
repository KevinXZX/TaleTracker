package database

import "taletracker.com/internal/model"

// TODO: Temporary in-memory database for testing purposes
// UserID -> Tales
var db = map[string][]model.Tale{
	"1": []model.Tale{
		{
			ID:     1,
			Title:  "The Tale of the First Tale",
			Blurb:  "This is the first tale. Lorem ipsum wodkwoqdko qdokqokdk dqkdkqkd qdkoqokdqokdkq dkqdkqkd qdokqkodqk dkqodkqodq dkqdkoq",
			Author: "John Doe",
			Review: model.Review{
				Rating: 2,
			},
		},
		{
			ID:     2,
			Title:  "The Tale of the second Tale",
			Blurb:  "This is the second tale.",
			Author: "John Doe",
			Review: model.Review{
				Rating: 3,
			},
		},
		{
			ID:     3,
			Title:  "The Tale of the third Tale",
			Blurb:  "This is the third tale.",
			Author: "John Doe",
			Review: model.Review{
				Rating: 4,
			},
		},
		{
			ID:     4,
			Title:  "The Tale of the Fourth Tale",
			Blurb:  "This is the Fourth tale.",
			Author: "John Doe",
			Review: model.Review{
				Rating: 5,
			},
		},
	},
}

func GetTales(userID string) []model.Tale {
	return db[userID]
}
