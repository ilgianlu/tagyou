package sqlrepository

import (
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Auth{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&Retry{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&Retain{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&Session{}); err != nil {
		return
	}
	if err := db.AutoMigrate(&Subscription{}); err != nil {
		return
	}
}
