package gorth

import (
	"testing"
)

func TestSetupPostgres(t *testing.T) {
	err := Setup("postgres", "abcdefg")
	if err == nil {
		t.Error("should only work for mysql (at this stage)")
	}
}

func TestDBExists(t *testing.T) {
	// Setup("mysql", "tt:tt@/gorthtest")
	err := Setup("mysql", "abc")
	if err == nil {
		t.Error("should not be able to ping with bad connection details")
	}
}

func TestRegisterNoSetup(t *testing.T) {
	db = nil
	user := User{
		FirstName: "bob",
		LastName:  "bobb",
		Email:     "bob@bob.com",
		Password:  "password",
	}
	err := Register(&user)

	if err == nil {
		t.Error("register should fail without a DB")
	}
}
