package models

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	ID   int
	Name string
	Mail string
	jwt.RegisteredClaims
}
