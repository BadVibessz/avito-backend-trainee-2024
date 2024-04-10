package user

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

func (r *Repo) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	rows, err := r.DB.NamedQueryContext(ctx,
		`INSERT INTO users (username, is_admin, hashed_password) VALUES (:username, :is_admin, :hashed_password) RETURNING *`,
		&user)
	if err != nil {
		return nil, err
	}

	var created entity.User

	if rows.Next() {
		if err = rows.StructScan(&created); err != nil {
			return nil, err
		}
	}

	return &created, nil
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := r.DB.QueryRowxContext(ctx, "SELECT * FROM users where username = $1", username)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.User

	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) CheckUniqueConstraints(ctx context.Context, username string) error {
	if user, err := r.GetUserByUsername(ctx, username); err == nil || user != nil {
		return ErrUsernameExists
	}

	return nil
}
