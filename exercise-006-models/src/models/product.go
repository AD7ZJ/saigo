package models

import "github.com/jmoiron/sqlx"

// Order ...
type Product struct {
	ID      int    `db:"product_id"`
	Product string `db:"product_name"`
}

func FindProduct(db *sqlx.DB, product string) (*Product, error) {
	return nil, nil
}
