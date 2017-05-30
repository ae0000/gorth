package gorth

import (
	"testing"
	"time"

	"github.com/ae0000/gorandom"
)

func TestConfirm(t *testing.T) {
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

	// The user should be in a state of pendingconfitmation
	if user.Status != StatusPendingConfirmation {
		t.Error("user should have a status of StatusPendingConfirmation")
	}

	// Confirm with a dodgy token
	err = Confirm(&user, "dodgytoken")
	if err == nil {
		t.Error("confirm is completely broken")
	}

	// Try and confirm
	err = Confirm(&user, user.Token)
	if err != nil {
		t.Errorf("user should of been confirmed.., got error %s", err)
	}

	if user.Status != StatusActive {
		t.Error("user status should be StatusActive")
	}
}

func TestConfirmTimesot(t *testing.T) {
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

	// The user should be in a state of pendingconformation
	if user.Status != StatusPendingConfirmation {
		t.Error("user should have a status of StatusPendingConfirmation")
	}

	// Set the TokenTimestamp to more than confirmationTokenMaxDuration
	duration, err := time.ParseDuration(confirmationTokenMaxDuration)
	if err != nil {
		t.Error(err)
	}

	user.TokenTimestamp = user.TokenTimestamp.Add(duration * -2)

	// Try and confirm
	err = Confirm(&user, user.Token)
	if err == nil {
		t.Error("user should not of been confirmed")
	}

	if user.Status == StatusActive {
		t.Error("user status should be StatusPendingConfirmation")
	}
}
