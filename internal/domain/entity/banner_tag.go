package entity

type BannerTag struct {
	ID       int `db:"id"`
	BannerID int `db:"banner_id"`
	TagID    int `db:"tag_id"`
}
