package sqlrepository

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/password"
	"gorm.io/gorm"
)

func AdminPasswordReset(db *gorm.DB, newPassword []byte) {
	log.Info().Msg("[ADMIN] reset admin password")
	admin := User{}
	if err := db.Debug().Where("username = ?", "admin").First(&admin).Error; err == nil {
		log.Error().Err(err)
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
}
