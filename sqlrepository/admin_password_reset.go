package sqlrepository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

func AdminPasswordReset(db *dbaccess.Queries, newPassword []byte) {
	slog.Debug("[ADMIN] create admin?")
	adminName := sql.NullString{String: "admin", Valid: true}
	_, err := db.GetUserByUsername(context.Background(), adminName)
	if err == nil {
		slog.Debug("[ADMIN] admin already present", "err", err)
		return
	}
	pwd, err := password.EncodePassword(newPassword)
	if err != nil {
		slog.Error("could not encode new password for admin")
		return
	}
	adminUser := dbaccess.CreateUserParams{
		Username:  adminName,
		Password:  pwd,
		CreatedAt: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
	}
	err = db.CreateUser(context.Background(), adminUser)
	if err != nil {
		slog.Error("could not create user", "err", err)
		return
	}
	slog.Info("[ADMIN] admin user created")
}
