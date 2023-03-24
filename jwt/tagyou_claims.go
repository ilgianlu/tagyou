package jwt

import jwt "github.com/golang-jwt/jwt/v5"

type TagyouClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}
