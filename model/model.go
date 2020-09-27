package model

import (
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Auth{})
	db.AutoMigrate(&Retry{})
	db.AutoMigrate(&Retain{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Subscription{})
}
