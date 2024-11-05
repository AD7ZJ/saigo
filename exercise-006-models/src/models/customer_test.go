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
func TestRefresh(t *testing.T) {
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

func TestCustomer(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
}
