package model

/**
Acl

[{"pattern": "/topic/1"}, {"pattern": "/topic/2"}]

*/

type Client struct {
	ID           uint
	ClientId     string
	Username     string
	Password     []byte
	SubscribeAcl string
	PublishAcl   string
	CreatedAt    int64
}
