package sqlrepository

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

type SubscriptionSqlRepository struct {
	Db *dbaccess.Queries
}

func mappingSubscription(sub dbaccess.Subscription) model.Subscription {
	return model.Subscription{
		ClientId:          sub.ClientID.String,
		Topic:             sub.Topic.String,
		RetainHandling:    uint8(sub.RetainHandling.Int64),
		RetainAsPublished: uint8(sub.RetainAsPublished.Int64),
		NoLocal:           uint8(sub.Qos.Int64),
		ProtocolVersion:   uint8(sub.ProtocolVersion.Int64),
		Enabled:           sub.Enabled.Int64 == 1,
		CreatedAt:         sub.CreatedAt.Int64,
		SessionID:         sub.SessionID.Int64,
		Shared:            sub.Shared.Int64 == 1,
		ShareName:         sub.ShareName.String,
	}
}

func mappingSubscriptions(subscriptions []dbaccess.Subscription) []model.Subscription {
	subs := []model.Subscription{}
	for _, sub := range subscriptions {
		subs = append(subs, mappingSubscription(sub))
	}
	return subs
}

func (s SubscriptionSqlRepository) CreateOne(sub model.Subscription) error {
	enbld := 0
	if sub.Enabled {
		enbld = 1
	}

	shrd := 0
	if sub.Shared {
		shrd = 1
	}

	_, err := s.Db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		ClientID:          sql.NullString{String: sub.ClientId, Valid: true},
		Topic:             sql.NullString{String: sub.Topic, Valid: true},
		RetainHandling:    sql.NullInt64{Int64: int64(sub.RetainHandling), Valid: true},
		RetainAsPublished: sql.NullInt64{Int64: int64(sub.RetainAsPublished), Valid: true},
		NoLocal:           sql.NullInt64{Int64: int64(sub.Qos), Valid: true},
		ProtocolVersion:   sql.NullInt64{Int64: int64(sub.ProtocolVersion), Valid: true},
		Enabled:           sql.NullInt64{Int64: int64(enbld), Valid: true},
		CreatedAt:         sql.NullInt64{Int64: sub.CreatedAt, Valid: true},
		SessionID:         sql.NullInt64{Int64: int64(sub.SessionID), Valid: true},
		Shared:            sql.NullInt64{Int64: int64(shrd), Valid: true},
		ShareName:         sql.NullString{String: sub.ShareName, Valid: true},
	})

	if err != nil {
		slog.Error("could not create subscription", "err", err)
		return err
	}
	return nil
}

func (s SubscriptionSqlRepository) FindToUnsubscribe(shareName string, topic string, clientId string) (model.Subscription, error) {
	sub, err := s.Db.GetSubscriptionToUnsubscribe(context.Background(), dbaccess.GetSubscriptionToUnsubscribeParams{
		ShareName: sql.NullString{String: shareName, Valid: true},
		Topic:     sql.NullString{String: topic, Valid: true},
		ClientID:  sql.NullString{String: clientId, Valid: true},
	})
	return mappingSubscription(sub), err
}

func (s SubscriptionSqlRepository) FindSubscriptions(topics []string, shared bool) []model.Subscription {
	nullTpcs := []sql.NullString{}
	for _, tpc := range topics {
		nullTpcs = append(nullTpcs, sql.NullString{String: tpc, Valid: true})
	}
	shrd := sql.NullInt64{Int64: 0, Valid: true}
	if shared {
		shrd = sql.NullInt64{Int64: 1, Valid: true}
	}
	subs, err := s.Db.GetSubscriptions(context.Background(), dbaccess.GetSubscriptionsParams{
		Topics: nullTpcs,
		Shared: shrd,
	})
	if err != nil {
		slog.Error("could not query for subscriptions", "err", err)
		return []model.Subscription{}
	}
	return mappingSubscriptions(subs)
}

func (s SubscriptionSqlRepository) FindOrderedSubscriptions(topics []string, shared bool) []model.Subscription {
	nullTpcs := []sql.NullString{}
	for _, tpc := range topics {
		nullTpcs = append(nullTpcs, sql.NullString{String: tpc, Valid: true})
	}
	shrd := sql.NullInt64{Int64: 0, Valid: true}
	if shared {
		shrd = sql.NullInt64{Int64: 1, Valid: true}
	}
	subs, err := s.Db.GetSubscriptionsOrdered(context.Background(), dbaccess.GetSubscriptionsOrderedParams{
		Topics: nullTpcs,
		Shared: shrd,
	})
	if err != nil {
		slog.Error("could not query for subscriptions", "err", err)
		return []model.Subscription{}
	}
	return mappingSubscriptions(subs)
}

func (s SubscriptionSqlRepository) DeleteByClientIdTopicShareName(clientId string, topic string, shareName string) error {
	return s.Db.DeleteSubscriptionByClientIdTopicShareName(context.Background(), dbaccess.DeleteSubscriptionByClientIdTopicShareNameParams{
		Topic:     sql.NullString{String: topic, Valid: true},
		ClientID:  sql.NullString{String: clientId, Valid: true},
		ShareName: sql.NullString{String: shareName, Valid: true},
	})
}

func (s SubscriptionSqlRepository) GetAll() []model.Subscription {
	subs, err := s.Db.GetAllSubscriptions(context.Background())
	if err != nil {
		slog.Error("could not query for subscriptions", "err", err)
		return []model.Subscription{}
	}
	return mappingSubscriptions(subs)
}
