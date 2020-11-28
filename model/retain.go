package model

type Retain struct {
	Topic              string `gorm:"primaryKey"`
	ApplicationMessage []byte
	CreatedAt          int64
}
