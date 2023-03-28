package password

import (
	"github.com/ilgianlu/tagyou/conf"
	"golang.org/x/crypto/bcrypt"
)

const DEFAULT_COST = 10

func ValidPassword(password []byte) bool {
	return len(password) > conf.PASSWORD_MIN_LENGTH
}

func CheckPassword(currentpassword []byte, guessedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(currentpassword, guessedPassword)
}

func EncodePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), DEFAULT_COST)
}
