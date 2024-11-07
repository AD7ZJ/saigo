package models

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// global db object to avoid creating it in every test
var db *sqlx.DB

// Special go function that allows us to call a setup and teardown function before all the tests in this file are run
func TestMain(m *testing.M) {
	SuiteSetup()
	code := m.Run()
	SuiteTeardown()
	os.Exit(code)
}

// Any setup that needs to happen prior to running the tests goes in here
func SuiteSetup() {
	// Set up the database connection
	var err error
	db, err = sqlx.Connect("postgres", "dbname=exercise_006 sslmode=disable")
	if err != nil {
		log.Fatalln("SuiteSetup: Failed to connect to the database:", err)
	}

	// start with an empty database to ensure consistency between test runs (FIXME: need to truncate other tables as well)
	_, err = db.Exec("TRUNCATE TABLE customers RESTART IDENTITY CASCADE")

	if err != nil {
		log.Fatalln("SuiteSetup: Failed to truncate the customers table:", err)
	}
}

// any cleanup activity that should happen after running the tests goes in here
func SuiteTeardown() {
	db.Close()
}

// Tests the creation of a new customer.
func TestNewCustomer(t *testing.T) {
	// Define test data
	email := "test@domain.com"
	firstName := "Test"
	lastName := "User"
	birthDate := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Call NewCustomer
	customer, err := NewCustomer(db, email, firstName, lastName, birthDate)

	// verify all of the returned paremeters are as expected
	assert.NoError(t, err, "Expected no error when inserting customer")
	assert.NotNil(t, customer, "Expected customer struct to be non-nil")
	assert.NotZero(t, customer.ID, "Expected customer ID to be set after insertion")
	assert.Equal(t, "test@domain.com", customer.Email)
	assert.Equal(t, "Test", customer.FirstName)
	assert.Equal(t, "User", customer.LastName)
}

// TestRefresh tests the Refresh function of the Customer struct
func TestRefreshCustomer(t *testing.T) {
	// Create a customer
	email := "IAmATestUser@domain.com"
	firstName := "TestFirstName"
	lastName := "TestLastName"
	birthDate := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	// NewCustomer() will add this test user to the database
	customer, err := NewCustomer(db, email, firstName, lastName, birthDate)
	assert.NoError(t, err, "Expected no error when creating the test customer data")

	// NewOrder() will add an order for this test user
	err = NewOrder(db, customer.ID, 1, 100)
	assert.NoError(t, err, "Expected no error when creating an order for the test customer data")

	// Now call Refresh and we'll verify the data populated in the Customer struct is correct.
	err = customer.Refresh(db)
	assert.NoError(t, err, "Expected no error when refreshing customer data")

	// check if the customer fields were loaded correctly
	assert.Equal(t, email, customer.Email, "Email does not match")
	assert.Equal(t, firstName, customer.FirstName, "First name does not match")
	assert.Equal(t, lastName, customer.LastName, "Last name does not match")
	assert.Equal(t, birthDate, customer.BirthDate.UTC(), "Birth date does not match")
	assert.False(t, customer.CreatedAt.IsZero(), "Created at should be populated")
	assert.False(t, customer.UpdatedAt.IsZero(), "Updated at should be populated")

	// verify the orders got populated correctly
	assert.NotNil(t, customer.Orders[0], "Orders should be populated")
	assert.NotZero(t, customer.Orders[0].ID, "Order ID should be populated")
	assert.NotZero(t, customer.Orders[0].ProductID, "Product ID should be populated")
	assert.NotZero(t, customer.Orders[0].Quantity, "Quantity should be populated")
	assert.Equal(t, customer.ID, customer.Orders[0].CustomerID, "Customer ID in order should match the customer ID")
	assert.False(t, customer.Orders[0].CreatedAt.IsZero(), "Order CreatedAt should be populated")
	assert.False(t, customer.Orders[0].UpdatedAt.IsZero(), "Order UpdatedAt should be populated")
}

// Test the UpdateCustomer function can correctly update a customer when the orders are completely different
func TestUpdateCustomerVerifyDBOrdersDeleted(t *testing.T) {
	// Create a customer
	email := "IAmATestUserForTestUpdate@domain.com"
	firstName := "TestFirstName"
	lastName := "TestLastName"
	birthDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	// add this test user to the database
	customer, err := NewCustomer(db, email, firstName, lastName, birthDate)
	assert.NoError(t, err, "Expected no error when creating the test customer data")

	// add some orders
	for i := 1; i < 5; i++ {
		err = NewOrder(db, customer.ID, i, 100)
		assert.NoError(t, err, "Expected no error when creating an order for the test customer data")
	}

	// Initialize some orders for a customer struct to make sure they would get correctly entered in the database
	order1 := &Order{
		ID:         0,
		ProductID:  4,
		Quantity:   10,
		CustomerID: customer.ID,
	}

	order2 := &Order{
		ID:         0,
		ProductID:  4,
		Quantity:   10,
		CustomerID: customer.ID,
	}

	customer.Email = "IAmAnUpdatedTestUserForTestUpdate@domain.com"
	customer.FirstName = "UpdatedFirstName"
	customer.LastName = "UpdatedLastName"
	customer.BirthDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	customer.Orders = []*Order{order1, order2}

	// Here's the test - call UpdateCustomer
	err = UpdateCustomer(db, customer)
	assert.NoError(t, err, "Expected no error when calling UpdateCustomer")

	// now call refresh and let's see if it all lines up
	err = customer.Refresh(db)
	assert.NoError(t, err, "Expected no error when calling Refresh")

	// check if the customer fields were loaded correctly
	assert.Equal(t, "IAmAnUpdatedTestUserForTestUpdate@domain.com", customer.Email, "Email does not match")
	assert.Equal(t, "UpdatedFirstName", customer.FirstName, "First name does not match")
	assert.Equal(t, "UpdatedLastName", customer.LastName, "Last name does not match")
	assert.Equal(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), customer.BirthDate.UTC(), "Birth date does not match")
	assert.False(t, customer.CreatedAt.IsZero(), "Created at should be populated")
	assert.False(t, customer.UpdatedAt.IsZero(), "Updated at should be populated")

	// verify the orders got populated correctly
	assert.NotNil(t, customer.Orders[0], "Orders should be populated")
	assert.Equal(t, 2, len(customer.Orders), "Expected 2 orders for this customer")
	assert.NotZero(t, customer.Orders[0].ID, "Order ID should be populated")
	assert.Equal(t, 4, customer.Orders[0].ProductID, "Product ID not as expected")
	assert.Equal(t, 10, customer.Orders[0].Quantity, "Product quantity not as expected")
	assert.Equal(t, customer.ID, customer.Orders[0].CustomerID, "Order customer ID does not match")
	assert.False(t, customer.Orders[0].CreatedAt.IsZero(), "Order CreatedAt should be populated")
	assert.False(t, customer.Orders[0].UpdatedAt.IsZero(), "Order UpdatedAt should be populated")

	assert.NotNil(t, customer.Orders[1], "Orders should be populated")
	assert.NotZero(t, customer.Orders[1].ID, "Order ID should be populated")
	assert.Equal(t, 4, customer.Orders[1].ProductID, "Product ID not as expected")
	assert.Equal(t, 10, customer.Orders[1].Quantity, "Product quantity not as expected")
	assert.Equal(t, customer.ID, customer.Orders[1].CustomerID, "Order customer ID does not match")
	assert.False(t, customer.Orders[1].CreatedAt.IsZero(), "Order CreatedAt should be populated")
	assert.False(t, customer.Orders[1].UpdatedAt.IsZero(), "Order UpdatedAt should be populated")
}

// test that UpdateCustomer can correctly update a customer, making simple updates to the order table
func TestUpdateCustomerVerifyDBOrdersUpdatedAndInserted(t *testing.T) {
	// Create a customer
	email := "IAmATestUserForTestUpdate1@domain.com"
	firstName := "TestFirstName"
	lastName := "TestLastName"
	birthDate := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	// add this test user to the database
	customer, err := NewCustomer(db, email, firstName, lastName, birthDate)
	assert.NoError(t, err, "Expected no error when creating the test customer data")

	// add some orders
	for i := 1; i < 5; i++ {
		err = NewOrder(db, customer.ID, i, 100)
		assert.NoError(t, err, "Expected no error when creating an order for the test customer data")
	}

	// call refresh to sync the customer struct with the database
	err = customer.Refresh(db)
	assert.NoError(t, err, "Expected no error when calling Refresh")

	// Initialize an additional order and we'll make sure it gets correctly entered in the database
	order1 := &Order{
		ID:         0,
		ProductID:  2,
		Quantity:   35,
		CustomerID: customer.ID,
	}

	customer.Email = "IAmAnUpdatedTestUserForTestUpdate1@domain.com"
	customer.FirstName = "UpdatedFirstName"
	customer.LastName = "UpdatedLastName"
	customer.BirthDate = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	customer.Orders = append(customer.Orders, order1)
	customer.Orders[0].ProductID = 4 // was 1, now set to 4 and we'll make sure it got updated
	customer.Orders[0].Quantity = 50 // was 100, now set to 50 and we'll make sure it got updated

	// Here's the test - call UpdateCustomer
	err = UpdateCustomer(db, customer)
	assert.NoError(t, err, "Expected no error when calling UpdateCustomer")

	// now call refresh and let's see if it all lines up
	err = customer.Refresh(db)
	assert.NoError(t, err, "Expected no error when calling Refresh")

	// check if the customer fields were loaded correctly
	assert.Equal(t, "IAmAnUpdatedTestUserForTestUpdate1@domain.com", customer.Email, "Email does not match")
	assert.Equal(t, "UpdatedFirstName", customer.FirstName, "First name does not match")
	assert.Equal(t, "UpdatedLastName", customer.LastName, "Last name does not match")
	assert.Equal(t, time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), customer.BirthDate.UTC(), "Birth date does not match")
	assert.False(t, customer.CreatedAt.IsZero(), "Created at should be populated")
	assert.False(t, customer.UpdatedAt.IsZero(), "Updated at should be populated")

	// we should have a total of 5 orders now for this customer
	assert.Equal(t, 5, len(customer.Orders), "Expected 5 orders for this customer")

	// verify the first order got updated correctly
	assert.NotNil(t, customer.Orders[0], "Orders should be populated")
	assert.NotZero(t, customer.Orders[0].ID, "Order ID should be populated")
	assert.Equal(t, 4, customer.Orders[0].ProductID, "Product ID not as expected")
	assert.Equal(t, 50, customer.Orders[0].Quantity, "Product quantity not as expected")
	assert.Equal(t, customer.ID, customer.Orders[0].CustomerID, "Order customer ID does not match")
	assert.False(t, customer.Orders[0].CreatedAt.IsZero(), "Order CreatedAt should be populated")
	assert.False(t, customer.Orders[0].UpdatedAt.IsZero(), "Order UpdatedAt should be populated")

	// verify the new order got added
	assert.NotNil(t, customer.Orders[4], "Orders should be populated")
	assert.NotZero(t, customer.Orders[4].ID, "Order ID should be populated")
	assert.Equal(t, 2, customer.Orders[4].ProductID, "Product ID not as expected")
	assert.Equal(t, 35, customer.Orders[4].Quantity, "Product quantity not as expected")
	assert.Equal(t, customer.ID, customer.Orders[4].CustomerID, "Order customer ID does not match")
	assert.False(t, customer.Orders[4].CreatedAt.IsZero(), "Order CreatedAt should be populated")
	assert.False(t, customer.Orders[4].UpdatedAt.IsZero(), "Order UpdatedAt should be populated")
}

func TestCustomer(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}
