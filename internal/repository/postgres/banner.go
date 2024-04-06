package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"math"

	"avito-backend-trainee-2024/internal/domain/entity"
)

type BannerRepo struct {
	DB *sqlx.DB
}

func NewBannerRepo(db *sqlx.DB) *BannerRepo {
	return &BannerRepo{
		DB: db,
	}
}

func (br *BannerRepo) GetAllBanners(ctx context.Context, limit, offset int) ([]*entity.Banner, error) {
	var query string

	if limit == math.MaxInt64 {
		query = fmt.Sprintf("SELECT * FROM banner ORDER BY feature_id OFFSET %v", offset)
	} else {
		query = fmt.Sprintf("SELECT * FROM banner ORDER BY feature_id LIMIT %v OFFSET %v", limit, offset)
	}

	rows, err := br.DB.QueryxContext(ctx, query)
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

	return banners, nil
}

func (br *BannerRepo) CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error) {
	// execute in transaction
	tx, err := br.DB.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	// firstly add content to Content table
	rows, err := tx.NamedQuery(`INSERT INTO content (title, text, url) VALUES (:title, :text, :url) RETURNING *`, &banner.Content)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	var content entity.Content

	if rows.Next() {
		if err = rows.StructScan(&content); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}

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
RETURNING (id, name, feature_id, is_active, created_at, updated_at)`, // todo: missing destination name row in *entity.Banner???????
		content.ID)

	rows, err = tx.NamedQuery(query, &banner)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}

		return nil, err
	}

	if rows.Next() {
		if err = rows.StructScan(&banner); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}

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
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return nil, rollbackErr
			}

			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	banner.Content = content

	return &banner, nil
}
