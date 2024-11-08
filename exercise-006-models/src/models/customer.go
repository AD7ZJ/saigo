package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Customer data struct
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
	err = queryOrders(db, c)

	if err != nil {
		return fmt.Errorf("failed to fetch orders for customer_id %d: %w", c.ID, err)
	}

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

	return customer, nil
}

// Remove a customer from the table by customer_id
func DeleteCustomer(db *sqlx.DB, id int) error {
	// first delete any orders associated with this customer
	_, err := db.Exec(`DELETE FROM orders WHERE customer_id = $1;`, id)

	if err != nil {
		return fmt.Errorf("failed to delete orders for customer %d: %w", id, err)
	}

	// now delete the customer
	_, err = db.Exec(`DELETE FROM customers WHERE customer_id = $1;`, id)

	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// update the customers and orders tables with data from the passed in Customer struct
func UpdateCustomer(db *sqlx.DB, u *Customer) error {
	// Update the customers table
	_, err := db.Exec(`UPDATE customers SET 
						email = $1, 
						first_name = $2, 
						last_name = $3, 
						birth_date = $4, 
						updated_at = NOW()
						WHERE customer_id = $5`, u.Email, u.FirstName, u.LastName, u.BirthDate, u.ID)

	if err != nil {
		return fmt.Errorf("failed to update customers for customer_id %d: %w", u.ID, err)
	}

	// get the current orders for this customer
	var currentOrders []*Order
	err = db.Select(&currentOrders, `SELECT order_id, product_id, quantity, customer_id, created_at, updated_at 
									FROM orders 
									WHERE customer_id = $1`, u.ID)

	if err != nil {
		return fmt.Errorf("failed to fetch orders for customer_id %d: %w", u.ID, err)
	}

	// If an order exists in the database but not in the incoming customer order list, delete it from the database
	for _, order := range currentOrders {
		if !orderExists(u.Orders, order.ID) {
			// Delete order from database, it no longer exists
			DeleteOrder(db, order.ID)
		}
	}

	// update existing orders and insert new ones
	for _, order := range u.Orders {
		if orderExists(currentOrders, order.ID) {
			// Update existing order
			_, err = db.Exec(`UPDATE orders SET 
							product_id = $1, 
							quantity = $2, 
							updated_at = NOW() 
							WHERE order_id = $3`, order.ProductID, order.Quantity, order.ID)

			if err != nil {
				return fmt.Errorf("failed to update order ID %d for customer ID %d: %w", order.ID, u.ID, err)
			}
		} else {
			err = NewOrder(db, u.ID, order.ProductID, order.Quantity)

			if err != nil {
				return fmt.Errorf("while updating customer id %d, failed to insert new order: %w", u.ID, err)
			}
		}
	}

	return nil
}

// helper function for UpdateCustomer()
func orderExists(orderList []*Order, orderId int) bool {
	for _, order := range orderList {
		// skip over any order id's of 0
		if order.ID != 0 && order.ID == orderId {
			return true
		}
	}
	return false
}

// helper function used in several places to populate a customer's order list from the database
func queryOrders(db *sqlx.DB, c *Customer) error {
	// Query to fetch orders for the customer
	var orders []*Order
	err := db.Select(&orders, `SELECT order_id, product_id, quantity, customer_id, created_at, updated_at 
								FROM orders 
								WHERE customer_id = $1 ORDER BY order_id`, c.ID)

	if err != nil {
		return err
	}
	// reference this new order slice in the customer struct
	c.Orders = orders

	return nil
}

// Returns the customer record based on email
func FindCustomerByEmail(db *sqlx.DB, email string) (*Customer, error) {
	var customer Customer

	// Query to fetch customer details and populate them into the Customer struct
	err := db.QueryRowx(`SELECT customer_id, email, first_name, last_name, birth_date, created_at, updated_at 
						 FROM customers 
						 WHERE email = $1`, email).StructScan(&customer)

	if err != nil {
		return nil, fmt.Errorf("failed to query customer data: %w", err)
	}

	// Query to fetch orders for the customer
	err = queryOrders(db, &customer)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for customer_id %d: %w", customer.ID, err)
	}

	return &customer, nil
}

// Returns the customer record based on customer_id
func FindCustomerByID(db *sqlx.DB, id int) (*Customer, error) {
	var customer Customer

	// Query to fetch customer details and populate them into the Customer struct
	err := db.QueryRowx(`SELECT customer_id, email, first_name, last_name, birth_date, created_at, updated_at 
						 FROM customers 
						 WHERE customer_id = $1`, id).StructScan(&customer)

	if err != nil {
		return nil, fmt.Errorf("failed to query customer data: %w", err)
	}

	// Query to fetch orders for the customer
	err = queryOrders(db, &customer)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for customer_id %d: %w", customer.ID, err)
	}

	return &customer, nil
}

// Returns a slice containing all customers in the database
func AllCustomers(db *sqlx.DB) ([]*Customer, error) {
	var customers []*Customer

	// Query to fetch all customers and populate them into the slice
	err := db.Select(&customers, `SELECT customer_id, email, first_name, last_name, birth_date, created_at, updated_at 
								  FROM customers ORDER BY customer_id`)

	if err != nil {
		return nil, fmt.Errorf("failed to query all customers: %w", err)
	}

	// get the orders for each customer
	for i, _ := range customers {
		err = queryOrders(db, customers[i])

		if err != nil {
			return nil, fmt.Errorf("failed to fetch orders for customer_id %d: %w", customers[i].ID, err)
		}
	}

	return customers, nil
}
