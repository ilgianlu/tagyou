package model

/*
*
Users enabled to connect to http apis
*/
type User struct {
	ID        int64
	Username  string
	Password  []byte
	CreatedAt int64
}
