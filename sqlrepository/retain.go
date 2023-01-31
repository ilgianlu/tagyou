package sqlrepository

import (
	"strings"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
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
	trimmedTopic := trimWildcard(subscribedTopic)
	var retains []model.Retain
	r.Db.Where("topic LIKE ?", strings.Join([]string{trimmedTopic, "%"}, "")).Find(&retains)
	return retains
}

func trimWildcard(topic string) string {
	lci := len(topic) - 1
	lc := topic[lci]
	if string(lc) == conf.WILDCARD_MULTI_LEVEL {
		topic = topic[:lci]
	}
	return topic
}
