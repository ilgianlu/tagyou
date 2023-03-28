package jwt

import jwt "github.com/golang-jwt/jwt/v5"

type TagyouClaims struct {
	UserId uint `json:"userId"`
	jwt.RegisteredClaims
}
