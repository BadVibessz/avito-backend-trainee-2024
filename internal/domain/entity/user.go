package entity

import "time"

type User struct {
	ID             int       `db:"id"`
	Username       string    `db:"username"`
	IsAdmin        bool      `db:"is_admin"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
