package sqlrepository

import (
	"context"
	"database/sql"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type ClientSqlRepository struct {
	Db *dbaccess.Queries
}

func MappedClients(clients []dbaccess.Client) []model.Client {
	mClients := []model.Client{}

	for _, client := range clients {
		mClients = append(mClients, MappedClient(client))
	}

	return mClients
}

func MappedClient(client dbaccess.Client) model.Client {
	return model.Client{
		ID:           uint(client.ID),
		ClientId:     client.ClientID.String,
		Username:     client.Username.String,
		Password:     client.Password,
		SubscribeAcl: client.SubscribeAcl.String,
		PublishAcl:   client.PublishAcl.String,
		CreatedAt:    client.CreatedAt.Int64,
	}
}

func (ar ClientSqlRepository) Create(client model.Client) error {
	createClientParams := dbaccess.CreateClientParams{
		ClientID:     sql.NullString{Valid: true, String: client.ClientId},
		Username:     sql.NullString{Valid: true, String: client.Username},
		Password:     client.Password,
		SubscribeAcl: sql.NullString{Valid: true, String: client.SubscribeAcl},
		PublishAcl:   sql.NullString{Valid: true, String: client.PublishAcl},
	}
	return ar.Db.CreateClient(context.Background(), createClientParams)
}

func (ar ClientSqlRepository) DeleteById(id uint) error {
	return ar.Db.DeleteClientById(context.Background(), int64(id))
}

func (ar ClientSqlRepository) GetAll() []model.Client {
	clients, err := ar.Db.GetAllClients(context.Background())
	if err != nil {
		return []model.Client{}
	}
	return MappedClients(clients)
}

func (ar ClientSqlRepository) GetByClientIdUsername(clientId string, username string) (model.Client, error) {
	params := dbaccess.GetClientByClientIdUsernameParams{
		ClientID: sql.NullString{String: clientId, Valid: true},
		Username: sql.NullString{String: username, Valid: true},
	}
	client, err := ar.Db.GetClientByClientIdUsername(context.Background(), params)
	if err != nil {
		return model.Client{}, err
	}
	mClient := MappedClient(client)
	return mClient, nil
}

func (ar ClientSqlRepository) GetById(id uint) (model.Client, error) {
	client, err := ar.Db.GetClientById(context.Background(), int64(id))
	if err != nil {
		return model.Client{}, err
	}
	mClient := MappedClient(client)
	return mClient, nil
}
