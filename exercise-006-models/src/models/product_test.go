package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Because this test file is in the package models, it gets access to global variables declared anywhere in the package.
// The following variables are declared in customer_test.go, as well as a SuiteSetup() and SuiteTeardown() function that
// sets up the database and gets it into a known state before tests are run.

// global db object to avoid creating it in every test
// var db *sqlx.DB

func TestFindProduct(t *testing.T) {
	// Call FindProduct()
	product, err := FindProduct(db, "kayak")

	// verify all of the returned paremeters are as expected
	assert.NoError(t, err, "Expected no error when searching for a product")
	assert.NotNil(t, product, "Expected product struct to be non-nil")

	// make sure it's case insensitive
	product, err = FindProduct(db, "KaYak")

	// verify all of the returned paremeters are as expected
	assert.NoError(t, err, "Expected no error when searching for a product")
	assert.NotNil(t, product, "Expected product struct to be non-nil")

	// make sure it matches a partial string
	product, err = FindProduct(db, "Paddl")

	// verify all of the returned paremeters are as expected
	assert.NoError(t, err, "Expected no error when searching for a product")
	assert.NotNil(t, product, "Expected product struct to be non-nil")
	assert.Equal(t, 3, product.ID)
	assert.Equal(t, "paddle", product.Product)
}
