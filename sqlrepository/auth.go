package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"gorm.io/gorm"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type Auth struct {
	ID                   uint   `gorm:"primaryKey"`
	ClientId             string `gorm:"index:auth_cred_idx,unique"`
	Username             string `gorm:"index:auth_cred_idx,unique"`
	Password             []byte
	SubscribeAcl         string
	PublishAcl           string
	CreatedAt            int64
	InputPassword        string `gorm:"-" json:",omitempty"`
	InputPasswordConfirm string `gorm:"-" json:",omitempty"`
}

type AuthSqlRepository struct {
	Db *gorm.DB
}

func MappedAuths(auths []Auth) []model.Auth {
	mAuths := []model.Auth{}

	for _, auth := range auths {
		mAuths = append(mAuths, MappedAuth(auth))
	}

	return mAuths
}

func MappedAuth(auth Auth) model.Auth {
	return model.Auth{
		ClientId:             auth.ClientId,
		Username:             auth.Username,
		Password:             auth.Password,
		SubscribeAcl:         auth.SubscribeAcl,
		PublishAcl:           auth.PublishAcl,
		CreatedAt:            auth.CreatedAt,
		InputPassword:        auth.InputPassword,
		InputPasswordConfirm: auth.InputPasswordConfirm,
	}
}

func (ar AuthSqlRepository) Create(auth model.Auth) error {
	return ar.Db.Create(&auth).Error
}

func (ar AuthSqlRepository) DeleteByClientIdUsername(clientId string, username string) error {
	return ar.Db.Where("client_id = ? and username = ?", clientId, username).Delete(&model.Auth{}).Error
}

func (ar AuthSqlRepository) GetAll() []model.Auth {
	auths := []Auth{}
	ar.Db.Find(&auths)
	return MappedAuths(auths)
}

func (ar AuthSqlRepository) GetByClientIdUsername(clientId string, username string) (model.Auth, error) {
	var auth Auth
	if err := ar.Db.Where("client_id = ? and username = ?", clientId, username).First(&auth).Error; err != nil {
		return model.Auth{}, err
	}

	mAuth := model.Auth{
		ClientId:             auth.ClientId,
		Username:             auth.Username,
		Password:             auth.Password,
		SubscribeAcl:         auth.SubscribeAcl,
		PublishAcl:           auth.PublishAcl,
		CreatedAt:            auth.CreatedAt,
		InputPassword:        auth.InputPassword,
		InputPasswordConfirm: auth.InputPasswordConfirm,
	}

	return mAuth, nil
}
