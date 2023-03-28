package jwt

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func TestCreateVerify(t *testing.T) {
	key := []byte("fantastic key")
	tokenString, err := CreateToken(key, "tagyou.at", 48, 1010)

	if err != nil {
		t.Errorf("unexpected token creation error %s", err)
	}

	_, claims, err := ParseToken(tokenString, key)

	if err != nil {
		t.Errorf("unexpected token parse error %s", err)
	}

	expTime, _ := claims.GetExpirationTime()
	expectedExpTime := jwt.NewNumericDate(time.Now().Add(48 * time.Hour))
	if expTime.Compare(expectedExpTime.Time) == 1 {
		t.Errorf("expiration time greater than expected %d > %d", expTime.Unix(), expectedExpTime.Unix())
	}

	res := VerifyToken(tokenString, key)

	if !res.Valid {
		t.Errorf("expected valid token, received false")
	}

	if res.UserId != 1010 {
		t.Errorf("expected user id %d, received %d", 1010, res.UserId)
	}
}
