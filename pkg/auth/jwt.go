package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenLifespan for responses
const TokenLifespan = time.Hour * 1

// SigningMethod is the method for authenticating the token
var SigningMethod = jwt.SigningMethodHS256

// TokenSecret is the in memory reference to the secret from env
var TokenSecret []byte

// GenerateToken creates a token for a valid authenticate request
func GenerateToken() (string, int64, error) {
	now := time.Now()
	now = now.UTC()

	expiration := now.Add(TokenLifespan).Unix()
	claims := &jwt.StandardClaims{
		Issuer:    "auth.evanmoncuso.com",
		Audience:  "*",
		ExpiresAt: expiration,
		IssuedAt:  now.Unix(),
	}

	token := jwt.NewWithClaims(SigningMethod, claims)
	tokenString, err := token.SignedString(TokenSecret)

	return tokenString, expiration, err
}

// ParseToken returns the claims of the token
// TODO: if you start adding information into the token, return that as well
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return TokenSecret, nil
	})
}

// ValidateToken is the checker for whather a token is valid
func ValidateToken(tokenString string) (bool, error) {
	_, err := ParseToken(tokenString)

	// I'm not interested in getting any of the information off the jwt, just verification and expiration

	if err != nil {
		return false, err
	}

	return true, nil
}
