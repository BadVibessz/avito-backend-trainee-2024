package mapper

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"avito-backend-trainee-2024/internal/handler/request"
	"avito-backend-trainee-2024/internal/handler/response"
)

func MapBannerToAdminBannerResponse(banner *entity.Banner) response.GetAdminBannerResponse {
	return response.GetAdminBannerResponse{
		ID:        banner.ID,
		TagIDs:    banner.TagIDs,
		FeatureID: banner.FeatureID,
		GetContentResponse: response.GetContentResponse{
			Title: banner.Content.Title,
			Text:  banner.Content.Text,
			Url:   banner.Content.Url,
		},
		IsActive:  banner.IsActive,
		CreatedAt: banner.CreatedAt,
		UpdatedAt: banner.UpdatedAt,
	}
}

func MapBannerToUserBannerResponse(banner *entity.Banner) response.GetUserBannerResponse {
	return response.GetUserBannerResponse{
		GetContentResponse: response.GetContentResponse{
			Title: banner.Content.Title,
			Text:  banner.Content.Text,
			Url:   banner.Content.Url,
		}}
}

func MapBannerToCreateBannerResponse(banner *entity.Banner) response.CreateBannerResponse {
	return response.CreateBannerResponse{ID: banner.ID}
}

func MapCreateBannerRequestToEntity(req *request.CreateBannerRequest) entity.Banner {
	return entity.Banner{
		TagIDs:    req.TagIDs,
		FeatureID: req.FeatureID,
		Content: entity.Content{
			Title: req.CreateContentRequest.Title,
			Text:  req.CreateContentRequest.Text,
			Url:   req.CreateContentRequest.Url,
		},
		IsActive: req.IsActive,
	}
}

func MapUpdateBannerRequestToEntity(req *request.UpdateBannerRequest) entity.Banner {
	return entity.Banner{
		TagIDs:    req.TagIDs,
		FeatureID: req.FeatureID,
		Content: entity.Content{
			Title: req.UpdateContentRequest.Title,
			Text:  req.UpdateContentRequest.Text,
			Url:   req.UpdateContentRequest.Url,
		},
		IsActive: req.IsActive,
	}
}
