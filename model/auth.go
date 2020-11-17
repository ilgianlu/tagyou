package model

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type Auth struct {
	ID                   uint   `gorm:"primaryKey"`
	ClientId             string `gorm:"index:auth_cred_idx,unique"`
	Username             string `gorm:"index:auth_cred_idx,unique"`
	Password             []byte `json:"-"`
	SubscribeAcl         string
	PublishAcl           string
	CreatedAt            time.Time
	InputPassword        string `gorm:"-" json:",omitempty"`
	InputPasswordConfirm string `gorm:"-" json:",omitempty"`
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

func (a *Auth) checkPassword(password string) error {
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

func CheckAuth(db *gorm.DB, clientId string, username string, password string) (bool, string, string) {
	var auth Auth
	if err := db.Where("client_id = ? and username = ?", clientId, username).First(&auth).Error; err != nil {
		return false, "", ""
	}

	if err := auth.checkPassword(password); err != nil {
		return false, "", ""
	}

	return true, auth.PublishAcl, auth.SubscribeAcl
}
