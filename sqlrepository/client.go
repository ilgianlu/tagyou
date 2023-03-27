package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"gorm.io/gorm"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type Client struct {
	ID                   uint   `gorm:"primaryKey"`
	ClientId             string `gorm:"index:client_cred_idx,unique"`
	Username             string `gorm:"index:client_cred_idx,unique"`
	Password             []byte
	SubscribeAcl         string
	PublishAcl           string
	CreatedAt            int64
	InputPassword        string `gorm:"-" json:",omitempty"`
	InputPasswordConfirm string `gorm:"-" json:",omitempty"`
}

type ClientSqlRepository struct {
	Db *gorm.DB
}

func MappedClients(clients []Client) []model.Client {
	mClients := []model.Client{}

	for _, client := range clients {
		mClients = append(mClients, MappedClient(client))
	}

	return mClients
}

func MappedClient(client Client) model.Client {
	return model.Client{
		ClientId:     client.ClientId,
		Username:     client.Username,
		Password:     client.Password,
		SubscribeAcl: client.SubscribeAcl,
		PublishAcl:   client.PublishAcl,
		CreatedAt:    client.CreatedAt,
	}
}

func (ar ClientSqlRepository) Create(client model.Client) error {
	return ar.Db.Create(&client).Error
}

func (ar ClientSqlRepository) DeleteByClientIdUsername(clientId string, username string) error {
	return ar.Db.Where("client_id = ? and username = ?", clientId, username).Delete(&Client{}).Error
}

func (ar ClientSqlRepository) GetAll() []model.Client {
	clients := []Client{}
	ar.Db.Find(&clients)
	return MappedClients(clients)
}

func (ar ClientSqlRepository) GetByClientIdUsername(clientId string, username string) (model.Client, error) {
	var client Client
	if err := ar.Db.Where("client_id = ? and username = ?", clientId, username).First(&client).Error; err != nil {
		return model.Client{}, err
	}

	mClient := MappedClient(client)

	return mClient, nil
}
