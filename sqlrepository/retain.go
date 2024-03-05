package sqlrepository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

type RetainSqlRepository struct {
	Db *dbaccess.Queries
}

func mappingRetain(ret dbaccess.Retain) model.Retain {
	return model.Retain{
		ClientID:           ret.ClientID.String,
		Topic:              ret.Topic.String,
		ApplicationMessage: ret.ApplicationMessage,
		CreatedAt:          ret.CreatedAt.Int64,
	}
}

func mappingRetains(rets []dbaccess.Retain) []model.Retain {
	retains := []model.Retain{}
	for _, ret := range rets {
		retains = append(retains, mappingRetain(ret))
	}
	return retains
}

func (r RetainSqlRepository) FindRetains(topics []string) []model.Retain {
	nullTpcs := []sql.NullString{}
	for _, tpc := range topics {
		nullTpcs = append(nullTpcs, sql.NullString{String: tpc, Valid: true})
	}
	retains, err := r.Db.GetRetains(context.Background(), nullTpcs)
	if err != nil {
		slog.Error("could not query for retains", "err", err)
		return []model.Retain{}
	}
	return mappingRetains(retains)
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
