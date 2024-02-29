package sqlrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
	"github.com/ilgianlu/tagyou/topic"
)

type RetainSqlRepository struct {
	Db *dbaccess.Queries
}

func mapped(ret dbaccess.Retain) model.Retain {
	return model.Retain{
		ClientID:           ret.ClientID.String,
		Topic:              ret.Topic.String,
		ApplicationMessage: ret.ApplicationMessage,
		CreatedAt:          ret.CreatedAt.Int64,
	}
}

func (r RetainSqlRepository) FindRetains(subscribedTopic string) []model.Retain {
	allRetains, err := r.Db.GetAllRetains(context.Background())
	if err != nil {
		return []model.Retain{}
	}

	retains := []model.Retain{}
	for _, ret := range allRetains {
		if topic.Match(ret.Topic.String, subscribedTopic) {
			retains = append(retains, mapped(ret))
		}
	}

	return retains
}

func (r RetainSqlRepository) Create(retain model.Retain) error {
	params := dbaccess.CreateRetainParams{
		ClientID:           sql.NullString{String: retain.ClientID, Valid: true},
		Topic:              sql.NullString{String: retain.Topic, Valid: true},
		ApplicationMessage: retain.ApplicationMessage,
		CreatedAt:          sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
	}
	return r.Db.CreateRetain(context.Background(), params)
}

func (r RetainSqlRepository) Delete(retain model.Retain) error {
	params := dbaccess.DeleteRetainByClientIdTopicParams{
		ClientID: sql.NullString{String: retain.ClientID, Valid: true},
		Topic:    sql.NullString{String: retain.Topic, Valid: true},
	}
	return r.Db.DeleteRetainByClientIdTopic(context.Background(), params)
}
