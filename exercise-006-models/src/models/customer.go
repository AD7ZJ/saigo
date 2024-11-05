package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Customer ...
type Customer struct {
	ID        int       `db:"customer_id"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	BirthDate time.Time `db:"birth_date"`
	Orders    []*Order
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Refresh loads the customer's details from the database based on customer_id.
func (c *Customer) Refresh(db *sqlx.DB) error {
	// Query to fetch customer details and populate them into the Customer struct
	err := db.QueryRowx(`SELECT customer_id, email, first_name, last_name, birth_date, created_at, updated_at 
						 FROM customers 
						 WHERE customer_id = $1`, c.ID).StructScan(c)

	if err != nil {
		return fmt.Errorf("failed to refresh customer data: %w", err)
	}

	// Query to fetch orders for the customer
	var orders []*Order
	err = db.Select(&orders, `SELECT order_id, product_id, quantity, customer_id, created_at, updated_at 
							  FROM orders 
							  WHERE customer_id = $1`, c.ID)

	if err != nil {
		return fmt.Errorf("failed to fetch customer orders: %w", err)
	}

	// Update the Orders field of the Customer struct
	c.Orders = orders

	return nil
}

// Create a new customer
func NewCustomer(db *sqlx.DB, email string, first_name string, last_name string, birth_date time.Time) (*Customer, error) {
	customer := &Customer{
		Email:     email,
		FirstName: first_name,
		LastName:  last_name,
		BirthDate: birth_date,
	}

	// add the customer to the database and request the newly created ID
	err := db.QueryRowx(`INSERT INTO customers (email, first_name, last_name, birth_date)
						 VALUES ($1, $2, $3, $4)
						 RETURNING customer_id`,
		email, first_name, last_name, birth_date).Scan(&customer.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to insert new customer: %w", err)
	}

	fmt.Printf("Created customer with ID: %d\n", customer.ID) // Debug log
	return customer, nil
}

// DeleteCustomer ...
func DeleteCustomer(db *sqlx.DB, id int) error {
	return nil
}

// UpdateCustomer ...
func UpdateCustomer(db *sqlx.DB, u *Customer) error {
	return nil
}

// FindCustomerByEmail ...
func FindCustomerByEmail(db *sqlx.DB, email string) (*Customer, error) {
	return nil, nil
}

// FindCustomerByID ...
func FindCustomerByID(db *sqlx.DB, id int) (*Customer, error) {
	return nil, nil
}

// AllCustomers ...
func AllCustomers(db *sqlx.DB) ([]*Customer, error) {
	return nil, nil
}
