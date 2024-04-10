package request

import "github.com/go-playground/validator/v10"

type UpdateContentRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url" validate:"omitempty,url"`
}

type UpdateBannerRequest struct {
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	UpdateContentRequest
	IsActive bool `json:"is_active"`
}

func (br *UpdateBannerRequest) Validate(valid *validator.Validate) error { return valid.Struct(br) }
