package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthClaims struct {
	Id         string       `json:"id"`
	Type       JWTClaimType `json:"type"`       // admin , collection
	Collection string       `json:"collection"` // auth collection id
	jwt.RegisteredClaims
}

type JWTClaimType string

const (
	JwtTypeAdmin      JWTClaimType = "admin"
	JwtTypeCollection JWTClaimType = "collection"
)
