package entity

import "time"

type Banner struct {
	ID        int   `db:"id"`
	TagIDs    []int `db:"tag_ids"`
	FeatureID int   `db:"feature_id"`
	Content
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
