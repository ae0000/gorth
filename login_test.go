package gorth

import (
	"fmt"
	"testing"

	"github.com/ae0000/gorandom"
)

func TestLogin(t *testing.T) {
	initTests()

	password := "password123"

	user := User{
		FirstName: "bob",
		LastName:  "bobb",
		Email:     gorandom.Email(),
		Password:  password,
	}
	err := Register(&user)

	if err != nil {
		t.Errorf("register failed with err: %s", err)
	}

	// Now try and login
	_, err = Login(user.Email, password)
	if err != nil {
		t.Errorf("login is not working, error: %s", err)
	}

	// Now try with a dodgy password
	_, err = Login(user.Email, "dodgy")
	if err == nil {
		t.Error("login is not working, logged in with dodgy password")
	}

	// Try and login with a dodgy email
	_, err = Login(fmt.Sprintf("%sabc", user.Email), password)
	if err == nil {
		t.Error("login is not working, logged in with dodgy email")
	}

}
