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
	Password             []byte `json:"-"`
	SubscribeAcl         string
	PublishAcl           string
	CreatedAt            int64
	InputPassword        string `gorm:"-" json:",omitempty"`
	InputPasswordConfirm string `gorm:"-" json:",omitempty"`
}

type AuthRepository struct {
	Db *gorm.DB
}

func (ar AuthRepository) GetByClientIdUsername(clientId string, username string) (model.Auth, error) {
	var auth Auth
	if err := ar.Db.Where("client_id = ? and username = ?", clientId, username).First(&auth).Error; err != nil {
		return model.Auth{}, err
	}

	return model.Auth(auth), nil
}
