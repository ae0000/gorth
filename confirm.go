package gorth

import (
	"errors"
	"time"
)

const (
	confirmationTokenMaxDuration = "120h" // 120 hours, ie.5 days
)

// Confirm that the user has access to their email address
func Confirm(user *User, confirmationToken string) error {
	// Check the tokens actually match
	if confirmationToken != user.Token {
		return errors.New("the token does not match for this user")
	}

	duration, err := time.ParseDuration(confirmationTokenMaxDuration)
	if err != nil {
		return err
	}

	if user.TokenTimestamp.Add(duration).Before(time.Now()) {
		return errors.New("duration for confirming has passed")
	}

	// User is confirmed
	user.Status = StatusActive
	return updateStatus(user.ID, user.Status)
}
