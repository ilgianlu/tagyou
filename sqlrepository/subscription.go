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

func (s SubscriptionSqlRepository) CreateOne(sub model.Subscription) error {
	return s.Db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		ClientID:          sql.NullString{String: sub.ClientId},
		Topic:             sql.NullString{String: sub.Topic},
		RetainHandling:    sql.NullInt64{Int64: int64(sub.RetainHandling)},
		RetainAsPublished: sql.NullInt64{Int64: int64(sub.RetainAsPublished)},
		NoLocal:           sql.NullInt64{Int64: int64(sub.Qos)},
		ProtocolVersion:   sql.NullInt64{Int64: int64(sub.ProtocolVersion)},
		Enabled:           sql.NullInt64{Int64: sub.Enabled},
		CreatedAt:         sql.NullInt64{Int64: sub.CreatedAt},
		SessionID:         sql.NullInt64{Int64: int64(sub.SessionID)},
		Shared:            sql.NullInt64{Int64: sub.Shared},
		ShareName:         sql.NullString{String: sub.ShareName},
	})
}

func (s SubscriptionSqlRepository) FindToUnsubscribe(shareName string, topic string, clientId string) (model.Subscription, error) {
	var sub model.Subscription
	if err := s.Db.Where("share_name = ? and topic = ? and client_id = ?", shareName, topic, clientId).First(&sub).Error; err != nil {
		return sub, err
	}
	return sub, nil
}

func (s SubscriptionSqlRepository) FindSubscriptions(topics []string, shared bool) []model.Subscription {
	subs := []model.Subscription{}
	if err := s.Db.Where("topic IN (?)", topics).Where("shared = ?", shared).Find(&subs).Error; err != nil {
		slog.Error("could not query for subscriptions", "err", err)
	}
	return subs
}

func (s SubscriptionSqlRepository) FindOrderedSubscriptions(topics []string, shared bool, orderField string) []model.Subscription {
	subs := []model.Subscription{}
	if err := s.Db.Where("topic IN (?)", topics).Where("shared = ?", shared).Order(orderField).Find(&subs).Error; err != nil {
		slog.Error("could not query for subscriptions", "err", err)
	}
	return subs
}

func (s SubscriptionSqlRepository) DeleteByClientIdTopicShareName(clientId string, topic string, shareName string) error {
	return s.Db.Where("share_name = ? and topic = ? and client_id = ?", shareName, topic, clientId).Delete(&Subscription{}).Error
}

func (s SubscriptionSqlRepository) GetAll() []model.Subscription {
	subs := []model.Subscription{}
	if err := s.Db.Find(&subs).Error; err != nil {
		slog.Error("could not query for subscriptions", "err", err)
	}
	return subs
}
