package persistence

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

type UserSqlRepository struct {
	Db *dbaccess.Queries
}

func MappedUsers(clients []dbaccess.User) []model.User {
	mUsers := []model.User{}

	for _, client := range clients {
		mUsers = append(mUsers, MappedUser(client))
	}

	return mUsers
}

func MappedUser(user dbaccess.User) model.User {
	return model.User{
		ID:        user.ID,
		Username:  user.Username.String,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Int64,
	}
}

func (ar UserSqlRepository) GetAll() []model.User {
	users, err := ar.Db.GetAllUsers(context.Background())
	if err != nil {
		slog.Error("could not query for users", "err", err)
		return []model.User{}
	}
	return MappedUsers(users)
}

func (ar UserSqlRepository) GetById(id int64) (model.User, error) {
	user, err := ar.Db.GetUserById(context.Background(), id)
	if err != nil {
		return model.User{}, err
	}
	mClient := MappedUser(user)
	return mClient, nil
}

func (ar UserSqlRepository) GetByUsername(username string) (model.User, error) {
	user, err := ar.Db.GetUserByUsername(context.Background(), sql.NullString{String: username, Valid: true})
	if err != nil {
		return model.User{}, err
	}
	mClient := MappedUser(user)
	return mClient, nil
}

func (ar UserSqlRepository) Create(user model.User) error {
	return ar.Db.CreateUser(context.Background(), dbaccess.CreateUserParams{
		Username:  sql.NullString{String: user.Username, Valid: true},
		Password:  user.Password,
		CreatedAt: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
	})
}

func (ar UserSqlRepository) DeleteById(id int64) error {
	return ar.Db.DeleteUserById(context.Background(), id)
}
