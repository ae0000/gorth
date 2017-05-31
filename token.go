package gorth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// How long is the token valid for
var (
	TokenDuration = time.Hour * 24
)

// Claims is used as a custom claim for the JWT process
type Claims struct {
	UserID int64  `json:"userid"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// Token will return a valid (JWT) token what is valid for TokenDuration. The
// secretSignKey should be a constant
func Token(user *User, secretSignKey []byte) (string, error) {
	// 	token := jwt.New(jwt.SigningMethodHS256)
	// FooFoo

	// Create the Claims
	claims := &Claims{
		user.ID,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenDuration).Unix(),
			Issuer:    "gorth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretSignKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken received from client
func ValidateToken(tokenString string, secretSignKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretSignKey, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// fmt.Printf("%v %v", claims.Email, claims.StandardClaims.ExpiresAt)
		return claims.Email, nil
	}
	return "", err
}
