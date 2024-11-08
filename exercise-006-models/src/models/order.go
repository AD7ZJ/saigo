package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Order data struct
type Order struct {
	ID         int       `db:"order_id"`
	ProductID  int       `db:"product_id"`
	Quantity   int       `db:"quantity"`
	CustomerID int       `db:"customer_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// NewOrder creates a new order in the database.
func NewOrder(db *sqlx.DB, customerID int, productID int, quantity int) error {
	_, err := db.Exec(`INSERT INTO orders ( product_id, quantity, customer_id )
					   VALUES ($1, $2, $3)`, productID, quantity, customerID)

	if err != nil {
		return fmt.Errorf("failed to create new order: %w", err)
	}

	return nil
}

// UpdateOrder updates an existing order in the database.
func UpdateOrder(db *sqlx.DB, o *Order) error {
	result, err := db.Exec(`UPDATE orders SET product_id = $1, quantity = $2, updated_at = NOW()
					   WHERE order_id = $3`, o.ProductID, o.Quantity, o.ID)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("no order found with id %d", o.ID)
	}

	return nil
}

// DeleteOrder removes an order from the database by order ID.
func DeleteOrder(db *sqlx.DB, orderID int) error {
	result, err := db.Exec(`DELETE FROM orders WHERE order_id = $1`, orderID)

	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return fmt.Errorf("no order found with id %d", orderID)
	}

	return nil
}
