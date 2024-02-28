package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthClaims struct {
	Id         string `json:"id"`
	Type       JWTType `json:"type"`       // admin , collection
	Collection string `json:"collection"` // auth collection id
	jwt.RegisteredClaims
}

type JWTType string

const (
	JWTAdmin      JWTType = "admin"
	JWTCollection JWTType = "collection"
)

