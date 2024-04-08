package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"

	"avito-backend-trainee-2024/internal/domain/entity"

	sliceutils "avito-backend-trainee-2024/pkg/utils/slice"
)

type BannerRepo struct {
	DB *sqlx.DB
}

func NewBannerRepo(db *sqlx.DB) *BannerRepo {
	return &BannerRepo{
		DB: db,
	}
}

// general method for fetching banners with where condition
func (br *BannerRepo) getBannersWhere(ctx context.Context, condition string, offset, limit int) ([]*entity.Banner, error) {
	var query string

	// TODO: maybe without join but two sql queries?
	if limit == math.MaxInt64 {
		query = fmt.Sprintf(`SELECT banner.id, name, feature_id, is_active, created_at, updated_at, title, text, url, c.content_id 
FROM banner JOIN public.content c ON c.content_id = banner.content_id %v
ORDER BY banner.feature_id OFFSET %v`, condition, offset)
	} else {
		query = fmt.Sprintf(`SELECT banner.id, name, feature_id, is_active, created_at, updated_at, title, text, url, c.content_id 
FROM banner JOIN public.content c ON c.content_id = banner.content_id %v
ORDER BY banner.feature_id LIMIT %V OFFSET %v`, condition, limit, offset)

	}

	// execute in transaction
	tx, err := br.DB.BeginTxx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	rows, err := tx.Queryx(query)
	if err != nil {
		return nil, err
	}

	var banners []*entity.Banner

	for rows.Next() {
		var banner entity.Banner

		if err = rows.StructScan(&banner); err != nil {
			return nil, err
		}

		banners = append(banners, &banner)
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// find tagIds associated with each banner
	rows, err = tx.Queryx("SELECT * FROM banner_tag") // todo: performance bottleneck (understand how to use offset and limit here for optimization)
	if err != nil {
		return nil, err
	}

	var bannerTagsMap = make(map[int][]entity.BannerTag)

	for rows.Next() {
		var bannerTag entity.BannerTag

		if err = rows.StructScan(&bannerTag); err != nil {
			return nil, err
		}

		_, exists := bannerTagsMap[bannerTag.BannerID]
		if !exists {
			bannerTagsMap[bannerTag.BannerID] = []entity.BannerTag{bannerTag}
		} else {
			bannerTagsMap[bannerTag.BannerID] = append(bannerTagsMap[bannerTag.BannerID], bannerTag)
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	for _, banner := range banners {
		_, exists := bannerTagsMap[banner.ID]
		if exists {
			for _, tag := range bannerTagsMap[banner.ID] {
				banner.TagIDs = append(banner.TagIDs, tag.TagID)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (br *BannerRepo) GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error) {
	return br.getBannersWhere(ctx, "", offset, limit)
}

func (br *BannerRepo) GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error) {
	banners, err := br.getBannersWhere(ctx, fmt.Sprintf("WHERE banner.feature_id = %v", featureID), 0, math.MaxInt64)
	if err != nil {
		return nil, err
	}

	for _, banner := range banners {
		if sliceutils.Equals(banner.TagIDs, tagIDs) {
			return banner, nil
		}
	}

	return nil, ErrNoSuchBanner
}

func (br *BannerRepo) CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error) {
	// execute in transaction
	tx, err := br.DB.BeginTxx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	// firstly add content to Content table
	rows, err := tx.NamedQuery(`INSERT INTO content (title, text, url) VALUES (:title, :text, :url) RETURNING *`, &banner.Content)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if rows.Next() {
		if err = rows.StructScan(&content); err != nil {
			return nil, err
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// then insert new banner into banner table
	query := fmt.Sprintf(`INSERT INTO banner (name, feature_id, content_id) 
VALUES (:name, :feature_id, %v) 
RETURNING id, name, feature_id, is_active, created_at, updated_at`, // todo: missing destination name row in *entity.Banner???????
		content.ID)

	rows, err = tx.NamedQuery(query, &banner)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		if err = rows.StructScan(&banner); err != nil {
			return nil, err
		}
	}

	// close rows
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// for each tag id in entity.banner create new row (banner.ID, tag.ID) in BannerTag table
	for _, tag := range banner.TagIDs {
		_, err = tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO banner_tag (banner_id, tag_id) VALUES (%v, %v)", banner.ID, tag))
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	banner.Content = content

	return &banner, nil
}

func (br *BannerRepo) DeleteBanner(ctx context.Context, id int) (*entity.Banner, error) {
	row, err := br.DB.QueryxContext(ctx, "DELETE FROM banner WHERE id = $1 RETURNING *", id)
	if err != nil {
		return nil, err
	}

	var banner entity.Banner

	if row.Next() {
		if err = row.StructScan(&banner); err != nil {
			return nil, err
		}
	}

	return &banner, nil
}
