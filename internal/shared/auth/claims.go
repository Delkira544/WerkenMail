package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Role Role `json:"role"`
	jwt.RegisteredClaims
}
