package sqlrepository

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/password"
	"gorm.io/gorm"
)

func AdminPasswordReset(db *gorm.DB, newPassword []byte) {
	log.Debug().Msg("[ADMIN] create admin?")
	admin := User{}
	if err := db.Where("username = ?", "admin").First(&admin).Error; err == nil {
		log.Debug().Err(err).Msg("[ADMIN] admin already present")
		return
	}
	pwd, err := password.EncodePassword(newPassword)
	if err != nil {
		log.Fatal().Msg("could not encode new password for admin")
	}
	admin.Username = "admin"
	admin.Password = pwd
	admin.CreatedAt = time.Now().Unix()
	db.Save(&admin)
	log.Info().Msg("[ADMIN] admin user created")
}
