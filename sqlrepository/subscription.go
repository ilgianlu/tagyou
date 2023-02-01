package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Subscription struct {
	ID                uint   `gorm:"primary_key"`
	ClientId          string `gorm:"uniqueIndex:sub_pars_idx"`
	Topic             string `gorm:"uniqueIndex:sub_pars_idx"`
	RetainHandling    uint8
	RetainAsPublished uint8
	NoLocal           uint8
	Qos               uint8
	ProtocolVersion   uint8
	Enabled           bool
	CreatedAt         int64
	SessionID         uint
	Shared            bool   `gorm:"default:false"`
	ShareName         string `gorm:"uniqueIndex:sub_pars_idx"`
}

type SubscriptionSqlRepository struct {
	Db *gorm.DB
}

func (s SubscriptionSqlRepository) CreateOne(sub model.Subscription) error {
	err := s.Db.Create(&sub).Error
	return err
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
		log.Error().Err(err).Msg("could not query for subscriptions")
	}
	return subs
}

func (s SubscriptionSqlRepository) FindOrderedSubscriptions(topics []string, shared bool, orderField string) []model.Subscription {
	subs := []model.Subscription{}
	if err := s.Db.Where("topic IN (?)", topics).Where("shared = ?", shared).Order(orderField).Find(&subs).Error; err != nil {
		log.Error().Err(err).Msg("could not query for subscriptions")
	}
	return subs
}

func (s SubscriptionSqlRepository) IsOnline(sub model.Subscription) bool {
	session := model.Session{}
	if err := s.Db.Where("id = ?", sub.SessionID).First(&session).Error; err != nil {
		return false
	} else {
		return session.Connected
	}
}

func (s SubscriptionSqlRepository) DeleteByClientIdTopicShareName(clientId string, topic string, shareName string) error {
	return s.Db.Where("share_name = ? and topic = ? and client_id = ?", shareName, topic, clientId).Delete(&Subscription{}).Error
}
