package response

import "time"

type GetContentResponse struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type GetAdminBannerResponse struct {
	ID        int   `json:"banner_id"`
	TagIDs    []int `json:"tag_ids"`
	FeatureID int   `json:"feature_id"`
	GetContentResponse
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
