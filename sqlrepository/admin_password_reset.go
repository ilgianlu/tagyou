package sqlrepository

import (
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/password"
	"gorm.io/gorm"
)

func AdminPasswordReset(db *gorm.DB, newPassword []byte) {
	slog.Debug("[ADMIN] create admin?")
	admin := User{}
	if err := db.Where("username = ?", "admin").First(&admin).Error; err == nil {
		slog.Debug("[ADMIN] admin already present", "err", err)
		return
	}
	pwd, err := password.EncodePassword(newPassword)
	if err != nil {
		slog.Error("could not encode new password for admin")
		return
	}
	admin.Username = "admin"
	admin.Password = pwd
	admin.CreatedAt = time.Now().Unix()
	db.Save(&admin)
	slog.Info("[ADMIN] admin user created")
}
