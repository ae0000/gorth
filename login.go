package gorth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Login is used to log the user in
// TODO add rate limited failed logins
func Login(email, password string) (*User, error) {
	// This user has not been authenticated yet
	noAuthUser := User{
		Email:    email,
		Password: password,
	}

	// Filter and validate the email and password
	noAuthUser.Filter()
	err := noAuthUser.Validate()
	if err != nil {
		return nil, err
	}

	// Get the user with that email
	dataUser, err := GetUserByEmail(noAuthUser.Email)
	if err != nil {
		return nil, err
	}

	// Compare the passowrd and hash
	err = bcrypt.CompareHashAndPassword(
		[]byte(dataUser.Password),
		[]byte(noAuthUser.Password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	// User is authenticated, update their LastLogin value and return the user
	// NOTE: we dont change the LastLogin on the current user object as we want
	// to provide the last login besides this one back
	err = updateLastLogin(dataUser.ID)

	return dataUser, nil
}
