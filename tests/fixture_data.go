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
	}
)
