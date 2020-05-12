package model

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	ClientId     string `gorm:"primary_key"`
	Username     string `gorm:"primary_key"`
	Password     []byte
	SubscribeAcl string
	PublishAcl   string
	CreatedAt    time.Time
}

func (a *Auth) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword(a.Password, []byte(password))
}

func (a *Auth) setPassword(password string) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("Error encrypting password", err)
		return err
	}
	a.Password = pwd
	return nil
}

func CheckAuth(db *gorm.DB, clientId string, username string, password string) bool {
	var auth Auth
	if db.First(&auth, "client_id = ? and username = ?", clientId, username).RecordNotFound() {
		return false
	}

	if err := auth.checkPassword(password); err != nil {
		return false
	}

	return true
}
