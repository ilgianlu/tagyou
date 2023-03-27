package model

type User struct {
	ID        uint
	Username  string
	Password  []byte
	CreatedAt int64
}
