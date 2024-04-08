package feature

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"context"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) GetFeatureByID(ctx context.Context, id int) (*entity.Feature, error) {
	row := r.DB.QueryRowxContext(ctx, "SELECT * FROM feature WHERE id = $1", id)

	if err := row.Err(); err != nil {
		return nil, err
	}

	var feature entity.Feature

	if err := row.StructScan(&feature); err != nil {
		return nil, err
	}

	return &feature, nil
}
