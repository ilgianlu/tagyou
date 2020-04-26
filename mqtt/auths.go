package mqtt

type Auths interface {
	createAuth(a Auth) error
	remAuth(clientId string, username string) error
	findAuth(clientId string) (Auth, bool)
	checkAuth(clientId string, username string, password string) bool
}
