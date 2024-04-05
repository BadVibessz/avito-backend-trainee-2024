package entity

import "time"

type Banner struct { // TODO: CHANGE POSTGRES BANNER SCHEMA
	ID        int
	Name      string
	TagIDs    []int
	FeatureID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TODO: сделать как массив или как one-to-many?
