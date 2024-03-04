package jwt

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateToken(signingKey []byte, issuer string, hoursDuration int, userId int64) (string, error) {

	claims := TagyouClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hoursDuration) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func ParseToken(tokenString string, signingKey []byte) (*jwt.Token, TagyouClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	}

	claims := TagyouClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)

	return parsedToken, claims, err
}

func VerifyToken(tokenString string, signingKey []byte) VerificationResult {
	parsedToken, claims, err := ParseToken(tokenString, signingKey)

	if err != nil {
		return VerificationResult{Valid: false, Err: err}
	}

	return VerificationResult{Valid: parsedToken.Valid, UserId: claims.UserId}
}
