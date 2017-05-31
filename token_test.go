package gorth

import (
	"fmt"
	"testing"

	"time"

	"github.com/ae0000/gorandom"
)

func TestToken(t *testing.T) {
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

	// Get token
	secretKey := []byte("SECRET")
	token, err := Token(&user, secretKey)

	if err != nil {
		t.Errorf("could not get token, error: %s", err)
	}

	if len(token) == 0 {
		t.Error("token is no good")
	}

	// Now try and validate the token
	email, err := ValidateToken(token, secretKey)
	if err != nil {
		t.Errorf("token should of been valid, err: %s", err)
	}
	if email != user.Email {
		t.Errorf("users email in token is wrong, we got %s: ", email)
	}

	// Try a different key
	// Now try and validate the token
	email, err = ValidateToken(token, []byte("ABCDE"))
	if err == nil {
		t.Error("token should not of been valid, bad key")
	}
	if len(email) > 0 {
		t.Error("should not of gotten a valid email")
	}

	// Edit the token
	token = fmt.Sprintf("AAAAA%s", token[5:])
	email, err = ValidateToken(token, secretKey)
	if err == nil {
		t.Error("token should not of been valid")
	}
	if len(email) > 0 {
		t.Error("should not of gotten a valid email")
	}
}

func TestTokenExpiry(t *testing.T) {
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

	// Set the token duration to 2 seconds
	TokenDuration = time.Millisecond * 500

	// Get token
	secretKey := []byte("SECRET")
	token, err := Token(&user, secretKey)

	if err != nil {
		t.Errorf("could not get token, error: %s", err)
	}

	if len(token) == 0 {
		t.Error("token is no good")
	}

	// Now try and validate the token
	email, err := ValidateToken(token, secretKey)
	if err != nil {
		t.Errorf("token should of been valid, err: %s", err)
	}
	if email != user.Email {
		t.Errorf("users email in token is wrong, we got %s: ", email)
	}

	// Wait for two seconds so the token expires
	time.Sleep(time.Second * 1)
	email, err = ValidateToken(token, secretKey)
	if err == nil {
		t.Error("token should of been invalid (expired)")
	}
	if email == user.Email {
		t.Error("token should of been invalid (expired)")
	}
}
