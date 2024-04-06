package sqlx

import "github.com/jmoiron/sqlx"

func Transaction(tx *sqlx.Tx, fun func(tx *sqlx.Tx)) error {
	return nil // TODO:
}
