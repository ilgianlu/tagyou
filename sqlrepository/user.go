package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"gorm.io/gorm"
)

type User struct {
	Username  string `gorm:"primaryKey"`
	Password  []byte
	CreatedAt int64
}

type UserSqlRepository struct {
	Db *gorm.DB
}

func MappedUsers(clients []User) []model.User {
	mUsers := []model.User{}

	for _, client := range clients {
		mUsers = append(mUsers, MappedUser(client))
	}

	return mUsers
}

func MappedUser(user User) model.User {
	return model.User{
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

func (ar UserSqlRepository) GetAll() []model.User {
	users := []User{}
	ar.Db.Find(&users)
	return MappedUsers(users)
}

func (ar UserSqlRepository) GetByUsername(username string) (model.User, error) {
	var user User
	if err := ar.Db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}

	mClient := MappedUser(user)

	return mClient, nil
}
