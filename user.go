package gorth

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ae0000/gorandom"
)

// Statuses for managing what users can and can't do``
const (
	StatusActive              = "active"
	StatusDeleted             = "deleted"
	StatusPendingConfirmation = "pendingconfirmation"
	StatusLocked              = "locked"
)

// User roles
const (
	RoleSuper  = "super"
	RoleAdmin  = "admin"
	RoleNormal = "normal"
)

// MinPasswordLength for how many chars are needed in a valid password
// NOTE: leaving as a var in case people want to edit it
var MinPasswordLength = 8

// Email validation regex
var emailPattern = regexp.MustCompile(`^(((([a-zA-Z]|\d|[!#\$%&'\*\+\-\/=\?\^_` + "`" + `{\|}~]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])+(\.([a-zA-Z]|\d|[!#\$%&'\*\+\-\/=\?\^_` + "`" + `{\|}~]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])+)*)|((\x22)((((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(([\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])|(\([\x01-\x09\x0b\x0c\x0d-\x7f]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}]))))*(((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(\x22)))@((([a-zA-Z]|\d|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])|(([a-zA-Z]|\d|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])([a-zA-Z]|\d|-|\.|_|~|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])*([a-zA-Z]|\d|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])))\.)+(([a-zA-Z]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])|(([a-zA-Z]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])([a-zA-Z]|\d|-|\.|_|~|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])*([a-zA-Z]|[\x{00A0}-\x{D7FF}\x{F900}-\x{FDCF}\x{FDF0}-\x{FFEF}])))\.?$`)

// User is used to hold and store the authentication details
type User struct {
	ID             int64
	FirstName      string
	LastName       string
	Email          string
	Password       string
	Photo          string
	Status         string
	Role           string
	Token          string
	TokenTimestamp time.Time
	LastLogin      time.Time
	Created        time.Time
}

// DefaultStatus sets the default status for a user (if not already set)
func (u *User) DefaultStatus() {
	if len(u.Status) == 0 {
		u.Status = StatusPendingConfirmation
	}
}

// DefaultPhoto sets the users photo to a gravagtar based on their email (if not
// already set)
func (u *User) DefaultPhoto() {
	if len(u.Photo) == 0 {
		u.Photo = gravatar(u.Email)
	}
}

// DefaultRole assigns a role if there is none given
func (u *User) DefaultRole() {
	if len(u.Role) == 0 {
		u.Role = RoleNormal
	}
}

// TokenConfirm sets up a token so that the user can confirm their email address
func (u *User) TokenConfirm() {
	// Create a token and send it to the user
	u.Token = gorandom.AlphaNumeric(30)
	u.TokenTimestamp = time.Now()
}

// Filter will remove any whitespace etc.
func (u *User) Filter() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
}

// Validate that the user contains valid info
func (u *User) Validate() error {
	// Validate email
	if !emailPattern.MatchString(u.Email) {
		return fmt.Errorf("email '%s' is not valid", u.Email)
	}

	// Validate password length
	if len(u.Password) < MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters in length", MinPasswordLength)
	}

	return nil
}

// gravatar returns a gravatar for the given email address with a size of 100,
// rating of "general" and a default image using "identicons"
func gravatar(email string) string {
	hasher := md5.New()
	hasher.Write([]byte(email))

	URL, _ := url.Parse("https://www.gravatar.com/")

	URL.Path = fmt.Sprintf("avatar/%v.jpg", hex.EncodeToString(hasher.Sum(nil)))
	params := url.Values{}
	params.Add("s", fmt.Sprintf("%v", 100))
	params.Add("r", fmt.Sprintf("%v", "g"))
	params.Add("d", fmt.Sprintf("%v", "identicon"))
	URL.RawQuery = params.Encode()
	return URL.String()
}

// emailUnique returns true if the email does not exist already in the db
func emailUnique(email string) bool {
	var id int
	err := db.QueryRow("SELECT ID FROM Users WHERE Email = ? LIMIT 1", email).Scan(&id)

	return id == 0 && err == sql.ErrNoRows
}

// GetUserByEmail returns the user specified by the email
func GetUserByEmail(email string) (*User, error) {
	u := User{}
	selectSQL := `
	SELECT ID, FirstName, LastName, Email, Password, Photo, Status,	Role, Token, 
		TokenTimestamp, LastLogin, Created 
	FROM Users
	WHERE Email = ? AND Status != "deleted"`

	err := db.QueryRow(selectSQL, email).Scan(&u.ID, &u.FirstName, &u.LastName,
		&u.Email, &u.Password, &u.Photo, &u.Status, &u.Role, &u.Token,
		&u.TokenTimestamp, &u.LastLogin, &u.Created)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
