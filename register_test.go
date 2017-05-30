package gorth

import (
	"fmt"
	"os"
	"testing"

	"github.com/ae0000/gorandom"
)

// gorth env:
/*
	export GORTH_TEST_DB_USER=tt
	export GORTH_TEST_DB_USER=tt
	export GORTH_TEST_DB_NAME=gorth_test
*/

func initTests() {
	if db != nil {
		return
	}

	databaseUser := os.Getenv("GORTH_TEST_DB_USER")
	databasePassword := os.Getenv("GORTH_TEST_DB_PASSWORD")
	databaseName := os.Getenv("GORTH_TEST_DB_NAME")

	if len(databasePassword) > 0 {
		databasePassword = fmt.Sprintf(":%s", databasePassword)
	}
	err := Setup(
		"mysql",
		fmt.Sprintf("%s%s@/%s", databaseUser, databasePassword, databaseName))
	if err != nil {
		panic(fmt.Sprintf("setup failed : %s", err))
	}

	_, err = db.Exec("DROP TABLE Users")
	if err != nil {
		fmt.Printf("setup failed to truncate Users table : %s", err)
	}

	CreateUsersTable()
}

func TestRegister(t *testing.T) {
	initTests()

	user := User{
		FirstName: "bob",
		LastName:  "bobb",
		Email:     gorandom.Email(),
		Password:  "password",
	}
	err := Register(&user)

	if err != nil {
		t.Errorf("register failed with err: %s", err)
	}

	if user.ID == 0 {
		t.Error("user register failed to return a user.ID")
	}
}

func TestFilter(t *testing.T) {
	initTests()

	user := User{
		FirstName: "  bob  ",
		LastName:  "bobb",
		Email:     gorandom.Email(),
		Password:  "password",
	}
	err := Register(&user)

	if err != nil {
		t.Errorf("register failed with err: %s", err)
	}

	// Firstname should now be "bob"
	if user.FirstName != "bob" {
		t.Error("user filter failed to trim...")
	}
}
func TestRegisterValidation(t *testing.T) {
	initTests()

	// Test small password
	user := User{
		FirstName: "bob",
		LastName:  "bobb",
		Email:     gorandom.Email(),
		Password:  "pass",
	}
	err := Register(&user)

	if err == nil {
		t.Error("register should of caught small password")
	}

	// Test bad email
	user = User{
		FirstName: "bob",
		LastName:  "bobb",
		Email:     "bob@",
		Password:  "password",
	}
	err = Register(&user)

	if err == nil {
		t.Error("register should of caught bad email")
	}

	// Test multiple emails
	user = User{
		FirstName: "bob",
		LastName:  "bobb",
		Email:     gorandom.Email(),
		Password:  "password",
	}

	err = Register(&user)

	if err != nil {
		t.Error("register should of been fine! ", err)
	}

	// Now register again with same email
	err = Register(&user)

	if err == nil {
		t.Error("register should of caught duplicate email")
	}
}
