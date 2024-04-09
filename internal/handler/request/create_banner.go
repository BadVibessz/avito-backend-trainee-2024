package request

import "github.com/go-playground/validator/v10"

type CreateContentRequest struct {
	Title string `json:"title" validate:"required,min=1"`
	Text  string `json:"text" validate:"required,min=1"`
	Url   string `json:"url" validate:"required,min=1,url"`
}

type CreateBannerRequest struct {
	TagIDs    []int `json:"tag_ids" validate:"required,min=1"`
	FeatureID int   `json:"feature_id" validate:"required,min=0"`
	CreateContentRequest
	IsActive bool `json:"is_active"`
}

func (br *CreateBannerRequest) Validate(valid *validator.Validate) error { return valid.Struct(br) }
