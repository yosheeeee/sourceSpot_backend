package models

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	ID   int64
	Name string
	Mail string
	jwt.RegisteredClaims
}
