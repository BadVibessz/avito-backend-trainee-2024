package tests

import (
	"avito-backend-trainee-2024/internal/domain/entity"
)

var (
	users = []entity.User{
		{
			Username:       "user",
			IsAdmin:        false,
			HashedPassword: "12345678",
		},
		{
			Username:       "admin",
			IsAdmin:        true,
			HashedPassword: "12345678",
		},
	}

	banners = []entity.Banner{
		{
			TagIDs:    []int{1, 2},
			FeatureID: 1,
			Content: entity.Content{
				Title: "title",
				Text:  "text",
				Url:   "http://url.com",
			},
			IsActive: true,
		},
		{
			TagIDs:    []int{1},
			FeatureID: 2,
			Content: entity.Content{
				Title: "title2",
				Text:  "text2",
				Url:   "http://url2.com",
			},
			IsActive: false,
		},
	}
)
