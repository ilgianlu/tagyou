package model

type User struct {
	Username             string
	Password             []byte `json:"-"`
	CreatedAt            int64
	InputPassword        string
	InputPasswordConfirm string
}
