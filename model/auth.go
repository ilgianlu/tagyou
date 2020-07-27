package model

import (
	"log"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type Auth struct {
	ID                   uint   `gorm:"primary_key"`
	ClientId             string `gorm:"unique_index:auth_cred_idx"`
	Username             string `gorm:"unique_index:auth_cred_idx"`
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
		log.Println("Error encrypting password", err)
		return err
	}
	a.Password = pwd
	return nil
}

func CheckAuth(db *gorm.DB, clientId string, username string, password string) (bool, string, string) {
	var auth Auth
	if db.First(&auth, "client_id = ? and username = ?", clientId, username).RecordNotFound() {
		return false, "", ""
	}

	if err := auth.checkPassword(password); err != nil {
		return false, "", ""
	}

	return true, auth.PublishAcl, auth.SubscribeAcl
}
