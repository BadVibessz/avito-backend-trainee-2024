package entity

import "avito-backend-trainee-2024/internal/domain/entity"

func InitNilFieldsOfBanner(banner1 *entity.Banner, banner2 *entity.Banner) {
	if banner1.FeatureID == 0 {
		banner1.FeatureID = banner2.FeatureID
	}

	if banner1.TagIDs == nil {
		banner1.TagIDs = banner2.TagIDs
	}

	if banner1.Content.Title == "" {
		banner1.Content.Title = banner2.Content.Title
	}

	if banner1.Content.Text == "" {
		banner1.Content.Text = banner2.Content.Text
	}

	if banner1.Content.Url == "" {
		banner1.Content.Url = banner2.Content.Url
	}
}
