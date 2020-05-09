package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Auth{})
	db.AutoMigrate(&Retry{})
	db.AutoMigrate(&Retain{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Subscription{})
}
