package models

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// global db object to avoid creating it in every test (defined in customer_test.go)
//var db *sqlx.DB

// this customer is guaranteed to be present in the database (defined in customer_test.go)
//var testCustomer *Customer

func TestNewAndDeleteOrder(t *testing.T) {
	// Attempt to create a new order
	err := NewOrder(db, testCustomer.ID, 1, 10)
	assert.NoError(t, err, "Expected no error when creating a new order")

	// refresh the customer - this syncs the testCustomer struct with the database
	err = testCustomer.Refresh(db)
	assert.NoError(t, err, "Expected no error when refreshing customer")

	// that customer should now have at least 4 orders
	assert.GreaterOrEqual(t, len(testCustomer.Orders), 4, "Customer should have at least 4 orders")

	// now delete the order
	lastOrder := testCustomer.Orders[len(testCustomer.Orders)-1]
	err = DeleteOrder(db, lastOrder.ID)
	assert.NoError(t, err, "Expected no error when calling DeleteOrder()")

	// verify the order is gone
	err = UpdateOrder(db, lastOrder)
	assert.Error(t, err, "Expected an error when calling UpdateOrder() on a non existent order")
}

func TestUpdateOrder(t *testing.T) {
	// refresh the customer - this syncs the testCustomer struct with the database
	err := testCustomer.Refresh(db)
	assert.NoError(t, err, "Expected no error when refreshing customer")

	// make some changes to the order
	testCustomer.Orders[0].ProductID = 4
	testCustomer.Orders[0].Quantity = 1000

	// call update
	UpdateOrder(db, testCustomer.Orders[0])

	// sync the testCustomer struct with the database
	err = testCustomer.Refresh(db)
	assert.NoError(t, err, "Expected no error when refreshing customer")

	// Verify the updated values in the database
	assert.Equal(t, 4, testCustomer.Orders[0].ProductID, "Product ID should be updated")
	assert.Equal(t, 1000, testCustomer.Orders[0].Quantity, "Quantity should be updated")
}
