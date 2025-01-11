package auth

import "github.com/golang-jwt/jwt/v5"

// this interface will be avaialble to all structs inside this auth package

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}