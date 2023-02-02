package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/topic"
	"gorm.io/gorm"
)

type Retain struct {
	Topic              string `gorm:"primaryKey"`
	ApplicationMessage []byte
	CreatedAt          int64
}

type RetainSqlRepository struct {
	Db *gorm.DB
}

func (r RetainSqlRepository) FindRetains(subscribedTopic string) []model.Retain {
	var allRetains []model.Retain
	r.Db.Find(&allRetains)

	retains := []model.Retain{}
	for _, ret := range allRetains {
		if topic.Match(ret.Topic, subscribedTopic) {
			retains = append(retains, ret)
		}
	}

	return retains
}

func (r RetainSqlRepository) Create(retain model.Retain) error {
	return r.Db.Create(&retain).Error
}

func (r RetainSqlRepository) Delete(retain model.Retain) error {
	return r.Db.Where("topic = ?", retain.Topic).Delete(&retain).Error
}
