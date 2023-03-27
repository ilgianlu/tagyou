package model

/*
*
Users enabled to connect to http apis
*/
type User struct {
	ID        uint
	Username  string
	Password  []byte
	CreatedAt int64
}
