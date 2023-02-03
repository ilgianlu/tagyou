package repository

import "github.com/ilgianlu/tagyou/model"

type SubscriptionRepository interface {
	Create(model.Subscription) error
	DeleteByClientIdTopicShareName(clientId string, topic string, shareName string) error
	FindToUnsubscribe(shareName string, topic string, clientId string) (model.Subscription, error)
	FindSubscriptions(topics []string, shared bool) []model.Subscription
	GetAll() []model.Subscription
}
