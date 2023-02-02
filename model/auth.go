package model

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"golang.org/x/crypto/bcrypt"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type Auth struct {
	ClientId             string
	Username             string
	Password             []byte `json:"-"`
	SubscribeAcl         string
	PublishAcl           string
	CreatedAt            int64
	InputPassword        string
	InputPasswordConfirm string
}

func (a *Auth) Validate() bool {
	if a.ClientId == "" {
		return false
	}
	if a.Username == "" {
		return false
	}
	return true
}

func (a *Auth) ValidPassword() bool {
	if len(a.InputPassword) > conf.PASSWORD_MIN_LENGTH && a.InputPassword == a.InputPasswordConfirm {
		return true
	}
	return false
}

func (a *Auth) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(a.Password, []byte(password))
}

func (a *Auth) SetPassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Err(err).Msg("Error encrypting password")
		return err
	}
	a.Password = pwd
	return nil
}
