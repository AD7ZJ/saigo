package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Product data type
type Product struct {
	ID      int    `db:"product_id"`
	Product string `db:"product_name"`
}

func FindProduct(db *sqlx.DB, productName string) (*Product, error) {
	var product Product

	// Query to fetch product details and populate them into the Product struct
	// the ILIKE % $1 % allows the search to be case insensitive and match part of the name - the string doesn't have to be an exact match.
	err := db.QueryRowx(`SELECT product_id, product_name 
						 FROM products 
						 WHERE product_name ILIKE '%' || $1 || '%'`, productName).StructScan(&product)

	if err != nil {
		return nil, fmt.Errorf("failed to query product data: %w", err)
	}

	return &product, nil
}
