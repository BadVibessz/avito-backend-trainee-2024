package tag

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"

	"avito-backend-trainee-2024/internal/domain/entity"
)

type Repo struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

func (r *Repo) GetTagsWithIDs(ctx context.Context, IDs []int) ([]*entity.Tag, error) {
	idsStr := ""

	for i, id := range IDs {
		idsStr += strconv.Itoa(id)

		if i != len(IDs)-1 {
			idsStr += ","
		}
	}

	rows, err := r.DB.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM tag WHERE id in (%v) ORDER BY id", idsStr))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []*entity.Tag

	for rows.Next() {
		var tag entity.Tag

		if err = rows.StructScan(&tag); err != nil {
			return nil, err
		}

		tags = append(tags, &tag)
	}

	return tags, nil
}

func (r *Repo) GetTagByID(ctx context.Context, id int) (*entity.Tag, error) {
	rows, err := r.DB.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM tag WHERE id = %v", id))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tag entity.Tag

	if rows.Next() {
		if err = rows.StructScan(&tag); err != nil {
			return nil, err
		}
	}

	return &tag, nil
}
