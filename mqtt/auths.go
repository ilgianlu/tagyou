package mqtt

type Auths interface {
	createAuth(a Auth) error
	remAuth(clientId string, username string) error
	findAuth(clientId string) (Auth, bool)
}
