package tag

import (
	"avito-backend-trainee-2024/internal/domain/entity"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
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
	IDsStr := ""

	for i, id := range IDs {
		IDsStr += strconv.Itoa(id)

		if i != len(IDs)-1 {
			IDsStr += ","
		}
	}

	rows, err := r.DB.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM tag WHERE id in (%v) ORDER BY id", IDsStr))
	if err != nil {
		return nil, err
	}

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
