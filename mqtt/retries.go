package mqtt

type Retries interface {
	addRetry(retry Retry) error
	remRetry(clientId string, packetIdentifier int) error
	findRetry(clientId string, packetIdentifier int) (Retry, bool)
}
