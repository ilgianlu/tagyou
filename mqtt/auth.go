package mqtt

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	clientId     string
	username     string
	password     []byte
	subscribeAcl string
	publishAcl   string
	createdAt    time.Time
}

func (a *Auth) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword(a.password, []byte(password))
}

func (a *Auth) setPassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("Error encrypting password", err)
		return err
	}
	a.password = pwd
	return nil
}
