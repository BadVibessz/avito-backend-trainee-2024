package postgres

import (
	"context"
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
