package repository

import "github.com/ilgianlu/tagyou/model"

type SubscriptionRepository interface {
	CreateOne(model.Subscription) error
	DeleteOne(model.Subscription) error
	FindToUnsubscribe(shareName string, topic string, clientId string) (model.Subscription, error)
	FindSubscriptions(topics []string, shared bool) []model.Subscription
	FindOrderedSubscriptions(topics []string, shared bool, orderField string) []model.Subscription
	IsOnline(sub model.Subscription) bool
}
