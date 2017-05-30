package gorth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Register the account
func Register(user *User) error {
	// Makre sure Setup() has been called (we have a db)
	if db == nil {
		return errors.New("you need to call Setup() before registering")
	}

	// Filter out whitespace
	user.Filter()

	// Validate email and password
	err := user.Validate()
	if err != nil {
		return err
	}

	// Make sure user email is unique
	if !emailUnique(user.Email) {
		return errors.New("that email is already in use")
	}

	// Convert password
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.DefaultStatus()
	user.DefaultPhoto()
	user.DefaultRole()

	// Set the confirmation token up
	user.TokenConfirm()

	// All good, save to the DB
	return insertUser(user)
}

// hashPassword takes the users password and turns it into a hash via bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
